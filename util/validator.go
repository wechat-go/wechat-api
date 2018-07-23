// validator
package util

import (
	gv "github.com/asaskevich/govalidator"
)

func IsPhone(str string) bool {
	pl := []string{"134", "135", "136", "137", "138", "139", "147", "150", "151", "152", "157", "158", "159", "178", "182", "183", "184", "187", "188", "198", //移动
		"130", "131", "132", "145", "155", "156", "175", "176", "185", "186", "166", //联通
		"133", "153", "173", "177", "180", "181", "189", "199"} //电信

	if len(str) != 11 {
		return false
	}
	if !gv.IsNumeric(str) {
		return false
	}
	nstr := SubStr(str, 0, 3)
	for _, v := range pl {
		if v == nstr {
			return true
		}
	}
	return false
}
