package jkv

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestLoad2JSON(t *testing.T) {
	defer tmTrack(time.Now())

	if jsonbytes, e := ioutil.ReadFile("./data/test2.json"); e == nil {
		jkv2 := NewJKV(string(jsonbytes))
		jkv2.Init()
		fPln("--- Init 2 ---")
		fPln(jkv2.Unfold())
		for ipath, val := range jkv2.mIPathValue {
			fPln(ipath, val)
		}

		if jsonbytes, e := ioutil.ReadFile("./data/test1.json"); e == nil {
			jkv1 := NewJKV(string(jsonbytes))
			jkv1.Init()
			fPln("--- Init 1 ---")

			
			for ipath, val := range jkv1.mIPathValue {
				fPln(ipath, val)
			}
		}
	}
}
