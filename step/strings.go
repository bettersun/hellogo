package step

import "strings"

// IsAllEmpty 切片的所有元素为空或空格
func IsAllEmpty(slice []string) (isEmpty bool) {
	if len(slice) == 0 {
		return true
	}

	isEmpty = true
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			isEmpty = false
			break
		}
	}

	return isEmpty
}

// IsInSuffix 判断目标字符串的末尾是否含有切片中指定的字符串
func IsInSuffix(s string, suffixes []string) bool {
	isIn := false

	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	for _, v := range suffixes {
		v = strings.TrimSpace(v)
		if strings.HasSuffix(s, v) {
			isIn = true
			break
		}
	}

	return isIn
}

// IsInSlice 判断目标字符串是否是在切片中
func IsInSlice(slice []string, s string) bool {
	if len(slice) == 0 {
		return false
	}

	isIn := false
	for _, f := range slice {
		if f == s {
			isIn = true
			break
		}
	}

	return isIn
}

// RemoveSpaceEmpty 删除字符串切片中的空字符串
func RemoveSpaceEmpty(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	for i, v := range slice {
		if strings.TrimSpace(v) == "" {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveSpaceEmpty(slice)
		}
	}

	return slice
}

// RemoveRepeated 去除字符串切片中的重复字符串
func RemoveRepeated(s []string) []string {
	m := make(map[string]int8, 0)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = 0
		}
	}

	var result []string
	for k, _ := range m {
		result = append(result, k)
	}

	return result
}
