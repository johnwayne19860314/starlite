package textx

import "strings"

// Blank returns true if trimmed string is empty
func Blank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// AnyBlank test if there is any empty string passed in as parameters
func AnyBlank(ss ...string) bool {
	for _, str := range ss {
		if Blank(str) {
			return true
		}
	}
	return false
}

// IsSame compares two string by trimming white spaces and ignore case
func IsSame(x string, y string) bool {
	return strings.EqualFold(strings.TrimSpace(x), strings.TrimSpace(y))
}

// ContainString returns true if e is found in s using a case in-sensitive comparison after trimming extra spaces
func ContainString(s []string, e string) bool {
	for _, a := range s {
		if IsSame(a, e) {
			return true
		}
	}
	return false
}
