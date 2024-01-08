package utils

import "strings"

var xxxUserWhiteList = func() []string {
	return []string{
		"legal@qq.com",
		"sucim@qq.com",
		//"sa-gfsh-bjm-agent@xxxmotors.com",
	}
}

func IsxxxUser(email string) bool {
	if InArrayString(email, xxxUserWhiteList()) {
		return true
	}
	return strings.HasSuffix(strings.TrimSpace(email), "@xxx.com")
}
