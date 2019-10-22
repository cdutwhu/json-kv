package jkv

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestLoad2JSON(t *testing.T) {
	defer tmTrack(time.Now())

	if jsonbytes, e := ioutil.ReadFile("./data/testjqrst2.json"); e == nil {
		jkv2 := NewJKV(string(jsonbytes))
		fPln("--- Init 2 ---")
		// fPln(jkv2.Unfold(-1, nil))
		// for ipath, val := range jkv2.mIPathValue {
		// 	fPln(ipath, val)
		// }

		fPln(" --------------------------------------------------------- ")

		if jsonbytes, e := ioutil.ReadFile("./data/testjqrst.json"); e == nil {
			jkv1 := NewJKV(string(jsonbytes))
			fPln("--- Init 1 ---")

			fPln(jkv1.Unfold(0, jkv2.mIPathValue))

			// for ipath, val := range jkv1.mIPathValue {
			// 	fPln(ipath, val)
			// }
		}
	}
}
