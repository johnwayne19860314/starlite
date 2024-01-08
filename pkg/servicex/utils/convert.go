package utils

import (
	"encoding/json"

	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

func ConvertByJson(source interface{}, target interface{}) error {
	marshal, err := json.Marshal(source)
	if err != nil {
		return errorx.Wrap(err, "failed to marshal object")
	}

	err = json.Unmarshal(marshal, target)
	if err != nil {
		return errorx.Wrap(err, "failed to unmarshal object")
	}

	return nil
}

func BoolP(src bool) *bool {
	return &src
}
