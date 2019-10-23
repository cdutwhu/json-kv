package jkv

import (
	"testing"
	"time"

	pp "./preprocess"
)

func TestLoad2JSON(t *testing.T) {
	defer tmTrack(time.Now())

	jsonPolicy := pp.FmtJSONFile("./data/test2.json")
	jsonPolicy = sReplaceAll(jsonPolicy, "\r\n", "\n")
	jsonData := pp.FmtJSONFile("./data/test1.json")
	jsonData = sReplaceAll(jsonData, "\r\n", "\n")
	jkvP := NewJKV(jsonPolicy)
	jkvD := NewJKV(jsonData)
	fPln(jkvD.Unfold(0, jkvP.mIPathValue))
}
