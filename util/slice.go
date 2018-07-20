// string
package util

import (
	"sort"
)

func NoDuplicateString(arr []string) (ret []string) {
	sort.Strings(arr)
	for i, v := range arr {
		if i > 0 && arr[i-1] == v {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}
