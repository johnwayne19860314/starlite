package stringutil

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	EMPTY_STRING = ""
)

// HasText check if string contain at least one none space character
func HasText(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

// AppendError append error msg after
func AppendError(s string, err error) string {
	eMsg := s
	if err != nil {
		eMsg = fmt.Sprintf("%s; error: %v", eMsg, err)
	}
	return eMsg
}

func JsonString(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// JsonStringOnly used to return json string, if have error, return empty string
func JsonStringOnly(obj interface{}) string {
	jsonStr, _ := JsonString(obj)
	return jsonStr
}
