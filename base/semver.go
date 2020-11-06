package base

import (
	"strings"
)

// IsValidSemver reports whether v is a valid semantic version string.
func IsValidSemver(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return false
	}
	if strings.Index(v, "v") == 0 || strings.Index(v, "V") == 0 {
		v = v[1:]
	}
	for _, vc := range strings.Split(v, ".") {
		if !parseInt(vc) {
			return false
		}
	}
	return true
}

func parseInt(v string) bool {
	if v == "" {
		return false
	}
	// fmt.Println(v, " len: ", len(v))
	if v[0] < '0' || '9' < v[0] {
		return false
	}
	i := 0
	for i < len(v) && '0' <= v[i] && v[i] <= '9' {
		i++
	}
	// fmt.Println(v, " i: ", i)
	if v[0] == '0' && i != 1 {
		return false
	}
	if len(v) != i && len(v) > 0 {
		return false
	}
	return true
}
