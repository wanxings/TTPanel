package util

import "strings"

// CompareVersion 比较版本号 v1 > v2 返回 1，v1 < v2 返回 -1，v1 == v2 返回 0
func CompareVersion(version1 string, version2 string) int {
	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")

	for i := 0; i < len(v1) || i < len(v2); i++ {
		var n1, n2 int
		if i < len(v1) {
			n1 = parseInt(v1[i])
		}
		if i < len(v2) {
			n2 = parseInt(v2[i])
		}
		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}
	return 0
}

func parseInt(s string) int {
	var result int
	for _, c := range s {
		result = result*10 + int(c-'0')
	}
	return result
}
