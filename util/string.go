// string
package util

import (
	"strings"
)

func StrToSlic(str, sep string) []string {
	s := strings.Split(str, sep)
	index := 0
	endIndex := len(s) - 1

	var result = make([]string, 0)
	for k, v := range s {
		if v == "" {
			result = append(result, s[index:k]...)
			index = k + 1
		} else if k == endIndex {
			result = append(result, s[index:endIndex+1]...)
		}
	}
	return result
}

func HandleStr(all string, target string, endstr string) string {
	if strings.Index(all, target) == -1 {
		return "-1"
	} else {
		start := strings.Index(all, target) + len(target)
		newStr := all[start:]
		end := strings.Index(newStr, endstr)
		return newStr[:end]
	}
}

func SubStr(s string, start, end int) string {
	rs := []byte(s)
	rl := len(rs)
	if start < 0 {
		start = 0
	}
	if start > end {
		start, end = end, start
	}
	if end > rl {
		end = rl
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
