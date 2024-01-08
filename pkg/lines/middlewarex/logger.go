package middlewarex

import (
	"time"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

func Log(appCtx appx.AppContext, reqCtx appx.ReqContext) (err error) {
	reqTime := time.Now()

	reqCtx.Gin().Next()

	respTime := time.Now()
	httpCode := reqCtx.Gin().Writer.Status()
	msg := "request done successfully"
	if httpCode >= 400 {
		msg = "request done failed"
		err = errorx.Errorf("request faild with code %v ", httpCode)
	}
	reqCtx.Info(msg,
		"HttpCode", httpCode,
		"Method", reqCtx.Gin().Request.Method, "Path", reqCtx.Gin().Request.URL.Path, "Query", reqCtx.Gin().Request.URL.RawQuery,
		"ReqTime", reqTime.Format("2006-01-02 15:04:05.000"), "RespTime", respTime.Format("2006-01-02 15:04:05.000"),
		"RequestDuration", float64(respTime.Sub(reqTime).Milliseconds())/1000)
	return
}
