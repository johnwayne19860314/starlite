package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

func JwtGetHeaderMap(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errorx.New("token format err")
	}
	jwtHeader, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(jwtHeader, &m)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return m, nil
}

func JwtGetHeaderKeyId(token string) (string, error) {
	m, err := JwtGetHeaderMap(token)
	if err != nil {
		return "", err
	}

	v, ok := m["kid"]
	if ok {
		return v.(string), nil
	}

	v, ok = m["x5t"]
	if ok {
		return v.(string), nil
	}

	return "", errorx.New("not found kid in jwt header")
}
