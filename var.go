// ********** ALL Based On JQ Formatted JSON ********** //

package jkv

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"time"

	u "github.com/cdutwhu/go-util"
	w "github.com/cdutwhu/go-wrappers"
)

type (
	b     = byte
	S     = w.Str
	I32   = w.I32
	JTYPE int
	// jStr  string
)

var (
	StartTrait = []byte{
		b('"'), // [array : string] OR [object : field]
		// b('{'), // [array : object]
		// b('n'),         // [array : null]
		// b('t'), b('f'), // [array : bool]
		// b('1'), b('2'), b('3'), b('4'), b('5'), b('6'), b('7'), b('8'), b('9'), b('-'), b('0'), // [array : number]
	}

	LF, SP, DQ = byte('\n'), byte(' '), byte('"')
)

var (
	fPf         = fmt.Printf
	fPln        = fmt.Println
	fSf         = fmt.Sprintf
	sSpl        = strings.Split
	sJoin       = strings.Join
	sCount      = strings.Count
	sReplace    = strings.Replace
	sReplaceAll = strings.ReplaceAll
	sIndex      = strings.Index
	IF          = u.IF
	MapKeys     = u.MapKeys
	MapKVs      = u.MapKVs
	MapsJoin    = u.MapsJoin
	MapsMerge   = u.MapsMerge
	MapPrint    = u.MapPrint
	SliceCover  = u.SliceCover
	MatchAssign = u.MatchAssign
	XIn         = u.XIn
	BLANK       = w.BLANK
)

const (
	TraitScan = "\n                                                                " // 64 spaces

	AOS0  = "[\n  {\n    "                                                 // 2, 4
	AOE0  = "\n  }\n]"                                                     // 2, 0
	AOS1  = "[\n    {\n      "                                             // 4, 6
	AOE1  = "\n    }\n  ]"                                                 // 4, 2
	AOS2  = "[\n      {\n        "                                         // 6, 8
	AOE2  = "\n      }\n    ]"                                             // 6, 4
	AOS3  = "[\n        {\n          "                                     // 8, 10
	AOE3  = "\n        }\n      ]"                                         // 8, 6
	AOS4  = "[\n          {\n            "                                 // 10, 12
	AOE4  = "\n          }\n        ]"                                     // 10, 8
	AOS5  = "[\n            {\n              "                             // 12, 14
	AOE5  = "\n            }\n          ]"                                 // 12, 10
	AOS6  = "[\n              {\n                "                         // 14, 16
	AOE6  = "\n              }\n            ]"                             // 14, 12
	AOS7  = "[\n                {\n                  "                     // 16, 18
	AOE7  = "\n                }\n              ]"                         // 16, 14
	AOS8  = "[\n                  {\n                    "                 // 18, 20
	AOE8  = "\n                  }\n                ]"                     // 18, 16
	AOS9  = "[\n                    {\n                      "             // 20, 22
	AOE9  = "\n                    }\n                  ]"                 // 20, 18
	AOS10 = "[\n                      {\n                        "         // 22, 24
	AOE10 = "\n                      }\n                    ]"             // 22, 20
	AOS11 = "[\n                        {\n                          "     // 24, 26
	AOE11 = "\n                        }\n                      ]"         // 24, 22
	AOS12 = "[\n                          {\n                            " // 26, 28
	AOE12 = "\n                          }\n                        ]"     // 26, 24

	TraitFV    = "\": "
	Trait1EndV = ",\n"
	Trait2EndV = "\n"
	Trait3EndV = "],\n"
	Trait4EndV = "]\n"

	pathLinker = "~~"
	// nowritestr = "\"------\""
	// noreadstr  = "\"******\""

	LvlMax = 20 // init 20 max level in advances
)

// readonly var
var (
	sTAOS      = []string{AOS0, AOS1, AOS2, AOS3, AOS4, AOS5, AOS6, AOS7, AOS8, AOS9, AOS10, AOS11, AOS12}
	sTAOE      = []string{AOE0, AOE1, AOE2, AOE3, AOE4, AOE5, AOE6, AOE7, AOE8, AOE9, AOE10, AOE11, AOE12}
	pLinker    = pathLinker
	rMD5, _    = regexp.Compile("[a-f0-9]{32}")
	rSHA1, _   = regexp.Compile("[a-f0-9]{40}")
	rSHA256, _ = regexp.Compile("[a-f0-9]{64}")
)

// var (
// 	// lsLvlIPaths is 2D slice for each Level's each ipath
// 	lsLvlIPaths = [][]string{
// 		{}, {}, {}, {}, {},
// 		{}, {}, {}, {}, {},
// 		{}, {}, {}, {}, {},
// 		{}, {}, {}, {}, {},
// 		{}, {}, {}, {}, {},
// 	}
// 	mPathMIdx   = make(map[string]int)    //
// 	mIPathPos   = make(map[string]int)    //
// 	mIPathValue = make(map[string]string) //
// 	mIPathOID   = make(map[string]string) //
// 	mOIDIPath   = make(map[string]string) //
// 	mOIDObj     = make(map[string]string) //
// 	mOIDLvl     = make(map[string]int)    // from 1 ...
// 	mOIDType    = make(map[string]JTYPE)  // oid's type is OBJ or ARR|OBJ
// )

// JKV :
type JKV struct {
	json        string
	lsLvlIPaths [][]string        // 2D slice for each Level's each ipath
	mPathMIdx   map[string]int    //
	mIPathPos   map[string]int    //
	mIPathValue map[string]string //
	mIPathOID   map[string]string //
	mOIDIPath   map[string]string //
	mOIDObj     map[string]string //
	mOIDLvl     map[string]int    // from 1 ...
	mOIDType    map[string]JTYPE  // oid's type is OBJ or ARR|OBJ
}

// ********************************************************** //

// SHA1Str :
func SHA1Str(s string) string {
	return fSf("%x", sha1.Sum([]byte(s)))
}

// SHA256Str :
func SHA256Str(s string) string {
	return fSf("%x", sha256.Sum256([]byte(s)))
}

// MD5Str :
func MD5Str(s string) string {
	return fSf("%x", md5.Sum([]byte(s)))
}

// T : JSON line Search Feature.
func T(lvl int) string {
	return TraitScan[0 : 2*lvl+1]
}

// PT :
func PT(T string) string {
	return T[0 : len(T)-2]
}

// NT :
func NT(T string) string {
	return T[0 : len(T)+2]
}

// TL : get T & L by nchar
func TL(nChar int) (string, int) {
	lvl := (nChar - 1) / 2
	return T(lvl), lvl
}

func tmTrack(start time.Time) {
	elapsed := time.Since(start)
	fPf("took %s\n", elapsed)
}
