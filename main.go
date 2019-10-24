package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"./jkv"
	pp "./preprocess"
)

func main() {
	exe := filepath.Base(os.Args[0])
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [-o='output'] <inputdata.json> <mask.json>\n", exe)
		return
	}

	inputfp, maskfp := os.Args[1], os.Args[2]
	if len(os.Args) == 4 {
		inputfp, maskfp = os.Args[2], os.Args[3]
	}

	outputPtr := flag.String("o", "result.json", "a string")
	flag.Parse()
	output := *outputPtr
	if !strings.HasSuffix(output, ".json") {
		output = output + ".json"
	}

	jsonData := strings.ReplaceAll(pp.FmtJSONFile(inputfp), "\r\n", "\n")
	jsonMask := strings.ReplaceAll(pp.FmtJSONFile(maskfp), "\r\n", "\n")
	jkvD, jkvM := jkv.NewJKV(jsonData), jkv.NewJKV(jsonMask)
	masked, _ := jkvD.Unfold(0, jkvM.MapIPathValue)
	ioutil.WriteFile(output, []byte(masked), 0666)
}
