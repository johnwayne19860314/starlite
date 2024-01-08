package utils

import (
	"bufio"
	"bytes"
	ttl "text/template"

	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

func ExecuteTemplate(content string, context interface{}) ([]byte, error) {
	parse, err := ttl.New("").Parse(content)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = parse.ExecuteTemplate(w, parse.ParseName, context)
	if err != nil {
		return nil, errorx.Wrap(err, "error when generating "+parse.ParseName)
	}
	err = w.Flush()
	if err != nil {
		return nil, errorx.Wrap(err, "error when generating "+parse.ParseName)
	}
	return buf.Bytes(), nil
}
