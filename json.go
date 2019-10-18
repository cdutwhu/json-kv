// ********** ALL are Based On JQ Formatted JSON ********** //

package json2

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"regexp"
)

// IsJSON :
func (jstr jStr) IsJSON() bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(jstr), &js) == nil
}

// scan :                         L   posarr     pos L
func (jstr jStr) scan() (int, map[int][]int, map[int]int, error) {
	Lm, offset := -1, 0
	if s := string(jstr); jstr.IsJSON() {
		mLvlFParr := make(map[int][]int)
		for i := 0; i <= LvlMax; i++ {
			mLvlFParr[i] = []int{}
		}
		mFPosLvl := make(map[int]int)

		// L0 : object
		if s[0] == '{' {
		NEXT:
			for i := 0; i < len(s); i++ {
				// modify levels for array-object
				if S(s[i:]).HPAny(sTAOS...) {
					offset++
				}
				if S(s[i:]).HPAny(sTAOE...) {
					offset--
				}

				for j := 3; j <= 39; j += 2 {
					T, L := TL(j)
					e := i + j

					if e < len(s) && s[i:e] == T && s[e] == '"' { // xIn(s[e], StartTrait) {
						// remove fakes (remove string array)
						for k := e + 1; k < len(s)-1; k++ {
							if s[k] == '"' {
								if s[k+1] != ':' {
									continue NEXT
								}
								break
							}
						}

						L -= offset
						mLvlFParr[L] = append(mLvlFParr[L], e) // store '"' position
						mFPosLvl[e] = L
						continue NEXT
					}
				}
			}
		}

		// remove empty levels
		for i := LvlMax; i >= 0; i-- {
			if v := mLvlFParr[i]; len(v) == 0 {
				delete(mLvlFParr, i)
				continue
			}
			Lm = i
			break
		}

		return Lm, mLvlFParr, mFPosLvl, nil
	}
	return Lm, nil, nil, errors.New("Not a valid JSON string")
}

// fields :
func (jstr jStr) fields(mLvlFPos map[int][]int) []map[int]string {
	s, keys := string(jstr), MapKeys(mLvlFPos).([]int)
	nLVL := keys[len(keys)-1]
	mFPosFNameList := []map[int]string{map[int]string{}} // L0 is empty
	for L := 1; L <= nLVL; L++ {                         // from L1 to Ln
		mFPosFName := make(map[int]string)
		for _, p := range mLvlFPos[L] {
			pe := p + 1
			for i := p + 1; s[i] != DQ; i++ {
				pe = i
			}
			mFPosFName[p] = s[p+1 : pe+1]
		}
		mFPosFNameList = append(mFPosFNameList, mFPosFName)
	}
	return mFPosFNameList
}

// pl2 -> pl1. pl1, pl2 are sorted.
func merge2fields(mFPosFName1, mFPosFName2 map[int]string) map[int]string {
	pl2Parent, pl2Path, iPos := make(map[int]string), make(map[int]string), 0
	pl1, pl2 := MapKeys(mFPosFName1).([]int), MapKeys(mFPosFName2).([]int)
	for _, p2 := range pl2 {
		for i := iPos; i < len(pl1)-1; i++ {
			if p2 > pl1[i] && p2 < pl1[i+1] {
				iPos = i
				pl2Parent[p2] = mFPosFName1[pl1[i]]
				break
			}
		}
		if p2 > pl1[len(pl1)-1] {
			pl2Parent[p2] = mFPosFName1[pl1[len(pl1)-1]]
		}
		pl2Path[p2] = pl2Parent[p2] + pLinker + mFPosFName2[p2]
	}
	return MapsJoin(mFPosFName1, pl2Path).(map[int]string)
}

// rely on "fields outcome"
func fPaths(mFPosFNameList ...map[int]string) map[int]string {
	mFPosFPath := make(map[int]string)
	nL := len(mFPosFNameList)
	posLists := make([][]int, nL)
	for i, mFPosFName := range mFPosFNameList {
		if len(mFPosFName) == 0 {
			continue
		}
		posList := MapKeys(mFPosFName).([]int)
		posLists[i] = posList
	}
	mFPosFNameMerge := mFPosFNameList[1]
	for i := 1; i < nL-1; i++ {
		mFPosFNameMerge = merge2fields(mFPosFNameMerge, mFPosFNameList[i+1])
		mFPosFPath = mFPosFNameMerge
	}
	return mFPosFPath
}

// ********************************************************** //

// fValuesOnObjs :
func fValuesOnObjs(strobjs string) (objs []string) {
	L, mLPStart, mLPEnd := 0, make(map[int][]int), make(map[int][]int)
	for p := 0; p < len(strobjs); p++ {
		c := strobjs[p]
		if c == '{' {
			L++
			mLPStart[L] = append(mLPStart[L], p)
		}
		if c == '}' {
			mLPEnd[L] = append(mLPEnd[L], p)
			L--
		}
	}
	pstarts, pends := mLPStart[1], mLPEnd[1]
	for i := 0; i < len(pstarts); i++ {
		s, e := pstarts[i], pends[i]
		objs = append(objs, strobjs[s:e+1])
	}
	return objs
}

// fValueType :
func (jstr jStr) fValueType(p int) (v string, t JTYPE) {
	getV := func(str string, s int) string {
		for i := s + 1; i < len(str); i++ {
			if S(str[i:]).HPAny(Trait1EndV, Trait2EndV) {
				return str[s:i]
			}
		}
		panic("Shouldn't be here @ getV")
	}
	getOV := func(str string, s int) string {
		nLCB, nRCB := 0, 0
		for i := s; i < len(str); i++ {
			switch str[i] {
			case '{':
				nLCB++
			case '}':
				nRCB++
			}
			if nLCB == nRCB && S(str[i:]).HPAny("},\n", "}\n") {
				return str[s : i+1]
			}
		}
		panic("Shouldn't be here @ getOV")
	}
	getAV := func(str string, s int) string {
		nLBB, nRBB := 0, 0
		for i := s; i < len(str); i++ {
			switch str[i] {
			case '[':
				nLBB++
			case ']':
				nRBB++
			}
			if nLBB == nRBB && S(str[i:]).HPAny("],\n", "]\n") {
				return str[s : i+1]
			}
		}
		panic("Shouldn't be here @ getAV")
	}

	s := string(jstr)
	v1c, pv := byte(0), 0
	for i := p; i < len(s); i++ {
		if S(s[i:]).HP(TraitFV) {
			pv = i + len(TraitFV)
			v1c = s[pv]
			break
		}
	}
	switch v1c {
	case DQ:
		t, v = STR, getV(s, pv)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		t, v = NUM, getV(s, pv)
	case 't', 'f':
		t, v = BOOL, getV(s, pv)
	case 'n':
		t, v = NULL, getV(s, pv)
	case '{':
		t, v = OBJ, getOV(s, pv)
	case '[':
		t, v = ARR, getAV(s, pv)
		{
			for i := pv + 1; i < len(s); i++ {
				c := s[i]
				if c == LF || c == SP {
					continue
				}
				switch c {
				case DQ:
					t |= STR
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
					t |= NUM
				case 't', 'f':
					t |= BOOL
				case 'n':
					t |= NULL
				case '{':
					t |= OBJ
				default:
					panic("Invalid JSON array element type")
				}
				break
			}
		}
	default:
		panic(fSf("[%d] @ Invalid JSON element type", p))
	}
	return
}

// pathType :
func (jstr jStr) pathType(fpath string, psSort []int, mFPosFPath map[int]string) JTYPE {
	for _, p := range psSort {
		if fpath == mFPosFPath[p] {
			_, t := jstr.fValueType(p)
			return t
		}
	}
	panic("Shouldn't be here @ posByPath")
}

// Init : prepare <>
func (jstr jStr) Init() error {
	if _, mLvlFParr, _, err := jstr.scan(); err == nil {
		mFPosFNameList := jstr.fields(mLvlFParr)
		// for iL, mPN := range mFPosFNameList {
		// 	fPln("------Level------:", iL)
		// 	for p, name := range mPN {
		// 		v, t := jstr.fValueType(p)
		// 		if !t.IsLeafValue() {
		// 			oid := uuid.New().String()
		// 			v = oid
		// 		}
		// 		fPln(t.Str(), name, v)
		// 	}
		// }

		fpaths := fPaths(mFPosFNameList...)
		for _, p := range MapKeys(fpaths).([]int) {
			v, t := jstr.fValueType(p)

			oid := ""
			if !t.IsLeafValue() {
				if !jStr(v).IsJSON() {
					panic("fetching value error")
				}
				oid = fSf("%x", sha1.Sum([]byte(v))) // e.g. 7c0873bf5b18e999ba3efcf482e80c9869b5bdf7
				mOIDObj[oid] = v
				v = oid
				if t.IsObj() || t.IsObjArr() {
					mOIDType[oid] = t
				}
				// switch {
				// case t.IsObj():
				// 	mOIDType[oid] = OBJ
				// case t.IsObjArr():
				// 	mOIDType[oid] = ARR | OBJ
				// }
			}

			fp := fpaths[p]
			fip := fSf("%s@%d", fp, mPathMIdx[fp])
			mPathMIdx[fp]++
			mIPathValue[fip] = v
			mIPathPos[fip] = p
			// fPf("DEBUG: %-5d%-5d[%-7s]  [%-60s]  %s\n", i, p, t.Str(), fip, v)

			if !t.IsLeafValue() {
				mIPathOID[fip] = oid
			}
		}

		//
		for ipath := range mIPathOID {
			n := sCount(ipath, pLinker) + 1
			lsLvlIPaths[n] = append(lsLvlIPaths[n], ipath)
			// fPf("%s [%d] %s\n", oid, n, ipath)
		}

		for i := 1; i < len(lsLvlIPaths); i++ {
			if Ls, LsPrev := lsLvlIPaths[i], lsLvlIPaths[i-1]; len(Ls) > 0 && len(LsPrev) > 0 {
				for _, ipathP := range LsPrev {
					pathP := S(ipathP).RmTailFromLast("@").V()
					chk := pathP + pLinker
					for _, ipath := range Ls {
						if S(ipath).HP(chk) {
							oidP, oid := mIPathOID[ipathP], mIPathOID[ipath]
							objP, obj := mOIDObj[oidP], mOIDObj[oid]
							mOIDObj[oidP] = sReplaceAll(objP, obj, oid)
							mOIDLvl[oidP], mOIDLvl[oid] = i-1, i
						}
					}
				}
			}
		}

		// [obj-arr whole value string] -> [aoid arr string]
		for oid := range mOIDObj {
			if strOIDs := ExpAOID(oid); strOIDs != "" {
				mOIDObj[oid] = strOIDs
				lvl := mOIDLvl[oid]
				for _, aoid := range AOIDStrToOIDs(strOIDs) {
					mOIDLvl[aoid] = lvl
				}
			}
		}

		return nil
	}

	return errors.New("scan error")
}

// ExpAOID : only can be used after mOIDType assigned
func ExpAOID(aoid string) string {
	if typ, ok := mOIDType[aoid]; ok && typ.IsObjArr() {
		strobjs := mOIDObj[aoid]
		objs := fValuesOnObjs(strobjs)
		for _, obj := range objs {
			oid := fSf("%x", sha1.Sum([]byte(obj)))
			mOIDObj[oid] = obj
			strobjs = sReplace(strobjs, obj, oid, 1)
		}
		return strobjs
	}
	return ""
}

// AOIDStrToOIDs :
func AOIDStrToOIDs(aoidstr string) (oids []string) {
	nComma := sCount(aoidstr, ",")
	r, _ := regexp.Compile("[a-f0-9]{40}")
	oids = r.FindAllString(aoidstr, -1)
	if aoidstr[0] != '[' || aoidstr[len(aoidstr)-1] != ']' || (oids != nil && len(oids) != nComma+1) {
		panic("error format @ AOIDStr")
	}
	return
}

// QueryPV : value ("*.*") for no path checking
func QueryPV(path string, value interface{}) (mLvlOIDs map[int][]string, maxLvl int) {
	mLvlOIDs = make(map[int][]string)
	nGen, valstr := sCount(path, pLinker), ""
	switch value.(type) {
	case string:
		valstr = fSf("\"%v\"", value)
	default:
		valstr = fSf("%v", value)
	}

	ignore := false
	if valstr == "\"*.*\"" {
		ignore = true
	}

	for i := 0; i < mPathMIdx[path]; i++ {
		ipath := fSf("%s@%d", path, i)
		if v, ok := mIPathValue[ipath]; ok && (v == valstr || ignore) {
			pos, PIPath := mIPathPos[ipath], ""
			for upgen := 1; upgen <= nGen; upgen++ {
				ppath := S(ipath).RmTailFromLastN(pLinker, upgen).V()
				for j := 0; j < mPathMIdx[ppath]; j++ {
					pipath := fSf("%s@%d", ppath, j)
					ppos := mIPathPos[pipath]
					if ppos > pos {
						break
					}
					PIPath = pipath
				}
				if pid, ok := mIPathValue[PIPath]; ok {
					if _, ok := mOIDObj[pid]; ok {
						iLvl := nGen - upgen + 1
						if !XIn(pid, mLvlOIDs[iLvl]) {
							mLvlOIDs[iLvl] = append(mLvlOIDs[iLvl], pid)
							if iLvl > maxLvl {
								maxLvl = iLvl
							}
						}
					}
					// break // if  search only the first one, break here !
				}
			}
		}
	}
	return mLvlOIDs, maxLvl
}

// Unfold :
func Unfold() string {
	if len(lsLvlIPaths[1]) == 0 {
		return "empty json"
	}
	if len(lsLvlIPaths[1]) != 0 && len(lsLvlIPaths[2]) == 0 {
		return "return whole json"
	}

	lvl2path := S(lsLvlIPaths[2][0]).RmTailFromLast("@").V()
	if mLvlOIDs, _ := QueryPV(lvl2path, "*.*"); mLvlOIDs != nil && len(mLvlOIDs) > 0 {
		for _, lvl := range MapKeys(mLvlOIDs).([]int) {
			for _, oid := range mLvlOIDs[lvl] {
				fPf("[%s] %s\n", oid, mOIDObj[oid])
				if mOIDType[oid].IsObjArr() {
					fPf(" *** ex ***: array object\n")
					for _, oid := range AOIDStrToOIDs(mOIDObj[oid]) {
						fPf("[%s] %s\n", oid, mOIDObj[oid])
					}
				}
			}
			fPln(" ----------------------------------------------------------------- ")
		}
	} else {
		fPln(mLvlOIDs)
	}

	return ""
}

// Query : unfinished ...
// func Query(paths []string, values []interface{}) map[int][]string {
// 	lPaths, lVals := len(paths), len(values)
// 	if lPaths != lVals {
// 		panic("paths' count & values' count are not same!")
// 	}

// 	mLvlOIDs, pathShared, maxLvl := make(map[int][]string), "", 0
// 	for i := 0; i < lPaths; i++ {
// 		path, value := paths[i], values[i]
// 		mlvloids, maxl := QueryPV(path, value)
// 		if len(mlvloids) == 0 {
// 			return nil
// 		}
// 		if i == 0 {
// 			mLvlOIDs, pathShared, maxLvl = mlvloids, path, maxl
// 			continue
// 		}

// 		pathShared = func(s1, s2 string) string {
// 			minl := int(math.Min(float64(len(s1)), float64(len(s2))))
// 			for i := 0; i < minl; i++ {
// 				if s1[i] != s2[i] {
// 					return s1[:i]
// 				}
// 			}
// 			return ""
// 		}(pathShared, path)

// 		if maxl > maxLvl {
// 			maxLvl = maxl
// 		}

// 		lvl := sCount(pathShared, pLinker)
// 		IDs1, IDs2 := mLvlOIDs[lvl], mlvloids[lvl]
// 	NEXT:
// 		for j, id1 := range IDs1 {
// 			for _, id2 := range IDs2 {
// 				if id1 == id2 {
// 					continue NEXT
// 				}
// 			}
// 			// remove id1 from IDs1
// 			IDs1[j] = IDs1[len(IDs1)-1]
// 			mLvlOIDs[lvl] = mLvlOIDs[lvl][:len(mLvlOIDs[lvl])-1]
// 		}
// 		if len(mLvlOIDs[lvl]) == 0 {
// 			return nil
// 		}

// 		// refresh mLvlIDs
// 		// if i > 0 {
// 		// 	for l := 1; l <= maxLvl; l++ {
// 		// 		IDs1, IDs2 = mLvlIDs[l], mlvlids[l]

// 		// 	}
// 		// }
// 	}
// 	return mLvlOIDs
// }
