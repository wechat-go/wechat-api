// string
package util

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func GetObj(body io.ReadCloser) map[string]interface{} {
	var obj map[string]interface{}
	con, _ := ioutil.ReadAll(body)
	json.Unmarshal(con, &obj)
	return obj
}

func OToInt(a interface{}) int {
	if a == nil {
		return 0
	} else {
		return a.(int)
	}
}

func OToFloat64(a interface{}) float64 {
	if a == nil {
		return 0
	} else {
		return a.(float64)
	}
}

func OToString(a interface{}) string {
	if a == nil {
		return ""
	} else {
		return a.(string)
	}
}
