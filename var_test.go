package jkv

import "testing"

func TestSHA(t *testing.T) {
	fPln("MD5", MD5Str("a"))       // 0cc175b9c0f1b6a831c399e269772661
	fPln("SHA1", SHA1Str("a"))     // 86f7e437faa5a7fce15d1ddcb9eaeaea377667b8
	fPln("SHA256", SHA256Str("a")) // ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb
}
