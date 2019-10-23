package preprocess

import (
	"fmt"
	"testing"
)

func TestJQ(t *testing.T) {
	// data := string(must(ioutil.ReadFile("../xjy/files/content.json")).([]byte))
	// fPln(FmtJSONStr(data, "./util/"))
	// fPln(" **************** ")
	// fPln(FmtJSONFile("../../xjy/files/content.json", "./util/"))
	// fPln()
	// jqrst := FmtJSONFile("./xapi.json", "./util/")
	// ioutil.WriteFile("./util/jqrst.json", []byte(jqrst), 0666)

	// *****************

	// fmt.Println(prepareJQ("../", "./", "./utils/"))
	// fmt.Println(os.Getwd())

	// fmt.Println(FmtJSONStr("{\"abc\": 123}", "../", "./", "./utils/"))

	// if data, err := ioutil.ReadFile("../data/sample.json"); err == nil {
	// 	// fmt.Println(string(data))
	// 	fmt.Println(FmtJSONStr(string(data), "../", "./", "./utils/"))
	// } else {
	// 	fmt.Println(err.Error())
	// }

	fmt.Println(FmtJSONFile("../../data/sample.json", "../", "./", "./utils/"))
}
