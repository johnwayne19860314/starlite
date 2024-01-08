package middlewarex

import (
	"fmt"
	"net/http"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/routinex"
)

const (
	ServerPanicMessage = "server panicked, please search logs with xid"
)

func Recovery(appCtx appx.AppContext, reqCtx appx.ReqContext) (err error){
	panickErr := routinex.RunSafeWithLogger(reqCtx.Gin().Next, reqCtx.GetLogger())
	if panickErr != nil {
		if appCtx.IsPrd() {
			reqCtx.Gin().JSON(http.StatusInternalServerError, map[string]string{
				"message": ServerPanicMessage,
				"xid":     reqCtx.GetXid(),
			})
			return
		}
		reqCtx.Gin().JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf("%+v", panickErr),
			"xid":     reqCtx.GetXid(),
		})
	}
	err = panickErr
	return
}
