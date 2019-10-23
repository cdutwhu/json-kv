package preprocess

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestJQ(t *testing.T) {
	// fmt.Println(prepareJQ("../", "./", "./utils/"))
	// fmt.Println(os.Getwd())

	fmt.Println(FmtJSONStr("{\"abc\": 123}", "../", "./", "./utils/"))

	if data, err := ioutil.ReadFile("../data/sample.json"); err == nil {
		// fmt.Println(string(data))
		fmt.Println(FmtJSONStr(string(data), "../", "./", "./utils/"))
	} else {
		fmt.Println(err.Error())
	}

	fmt.Println(FmtJSONFile("../../data/sample.json", "../", "./", "./utils/"))
}
