package jkv

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	defer tmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("./data/testjqrst.json"); e == nil {
		jkv := NewJKV(string(jsonbytes))
		LVL, mLvlFParr, mFPosLvl, _ := jkv.scan()
		fPln("levels:", LVL)
		for k, v := range mLvlFParr {
			fPln(k, v)
		}
		for k, v := range mFPosLvl {
			fPln(k, v)
		}
	}
}

func TestFieldByPos(t *testing.T) {
	defer tmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("./data/testjqrst.json"); e == nil {
		jkv := NewJKV(string(jsonbytes))
		LVL, mLvlFParr, _, _ := jkv.scan()
		// for k, v := range mLvlFParr {
		// 	fPln(k, v)
		// }
		mFPosFNameList := jkv.fields(mLvlFParr)
		for i := 1; i <= LVL; i++ {
			fPln("---------------->", i)
			mFPosFName := mFPosFNameList[i]
			for k, v := range mFPosFName {
				_, t := jkv.fValueType(k)
				fPf("%-8d%-20s%-10s\n", k, v, t.Str())
				// if t.IsPrimitive() {
				// 	fPf("%-8d%-20s%-10s\n", k, v, t.Str())
				// } else {
				// 	fPf("%-8d%-20s\n", k, v)
				// }
			}
		}
	}
}

func TestFType(t *testing.T) {
	defer tmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("./data/testjqrst.json"); e == nil {
		jkv := NewJKV(string(jsonbytes))
		value, typ := jkv.fValueType(1617)
		fPln(typ.Str())

		if typ == ARR|OBJ {
			objs := fValuesOnObjs(value)
			fPln(objs[1])
		}

	}
}

func TestInit(t *testing.T) {
	defer tmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("./data/infmt.json"); e == nil {
		jkv := NewJKV(string(jsonbytes))
		jkv.Init()
	}
	fPln("break")
}

func TestUnfold(t *testing.T) {
	defer tmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("./data/testjqrst.json"); e == nil {
		jkv := NewJKV(string(jsonbytes))
		jkv.Init()
		fPln("--- Init ---")
		fPln(jkv.Unfold())

		fPln(jkv.mOIDLvl["ab6f49bdc782bb205ce4682a141d8ad9854f2a34"])
		fPln(jkv.mOIDLvl["c10256e43dba81f035c95631437cf45ba41c0170"])
	}
}

func TestQuery(t *testing.T) {
	defer tmTrack(time.Now())
	param := "NAPTestItemLocalId"
	value := "x00101935"
	if jsonbytes, e := ioutil.ReadFile("./data/testjqrst.json"); e == nil {
		// jstr := jStr(string(jsonbytes))
		jkv := NewJKV(string(jsonbytes))
		jkv.Init()

		fPln("--- Init ---")

		path := func(string) string {
			return "NAPCodeFrame~~TestletList~~Testlet~~TestItemList~~TestItem~~TestItemContent~~NAPTestItemLocalId"
		}(param)

		//path1 := "NAPCodeFrame~~TestletList~~Testlet~~TestItemList~~TestItem~~TestItemContent~~NAPTestItemLocalId"
		//value1 := "\"x00101923-00-AIA\""
		// path2 := "NAPCodeFrame~~TestletList~~Testlet~~NAPTestletRefId"
		// value2 := "\"2b7c9606-09b9-43c2-a935-6a2db78bf2c9\""

		if mLvlOIDs, maxL := jkv.QueryPV(path, value); mLvlOIDs != nil && len(mLvlOIDs) > 0 {

			for _, oid := range mLvlOIDs[maxL] {
				fPln(oid, jkv.mOIDObj[oid])
			}

			// for _, lvl := range MapKeys(mLvlOIDs).([]int) {
			// 	for _, oid := range mLvlOIDs[lvl] {
			// 		fPf("[%s] %s\n", oid, mOIDObj[oid])
			// 		if mOIDType[oid].IsObjArr() {
			// 			fPf("ex: array object\n")
			// 			for _, oid := range AOIDStrToOIDs(mOIDObj[oid]) {
			// 				fPf("[%s] %s\n", oid, mOIDObj[oid])
			// 			}
			// 		}
			// 	}
			// 	fPln(" ----------------------------------------------------------------- ")
			// }

			// fPln(mOIDLvl["fe7262a928bbe05f8a42bab98ebec56a8e1e9379"])
			// fPln(mOIDLvl["887450b46a52ccad78f6a74f34c2699c649b17cd"])
		}
	}
}
