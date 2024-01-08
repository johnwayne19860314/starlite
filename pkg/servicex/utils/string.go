package utils

import (
	"strings"

	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/textx"
)

func InArrayString(str string, array []string) bool {
	for _, elt := range array {
		if elt == str {
			return true
		}
	}
	return false
}

func InArrayAnyString(src []string, array []string) bool {
	for _, v := range src {
		if InArrayString(v, array) {
			return true
		}
	}
	return false
}

func Blank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func StringIsNilOrEmpty(s *string) bool {
	if s == nil {
		return true
	}
	return textx.Blank(*s)
}
