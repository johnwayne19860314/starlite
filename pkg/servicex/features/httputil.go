package features

import "github.startlite.cn/itapp/startlite/pkg/lines/errorx"

type ClientResponse interface {
	StatusCode() int
	GetBody() []byte
}

func CheckClientResponse(resp ClientResponse, err error) error {
	if err != nil {
		return errorx.WithStack(err)
	}
	if resp.StatusCode() >= 400 {
		return errorx.Errorf("resp err: StatusCode: %d, Body: %s", resp.StatusCode(), string(resp.GetBody()))
	}

	return nil
}
