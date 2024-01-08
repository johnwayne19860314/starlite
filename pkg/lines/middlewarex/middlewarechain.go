package middlewarex

import (
	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
)

type MiddlewareChain []gin.HandlerFunc

func AppendMiddlewares(appCtx appx.AppContext, middlewares ...gin.HandlerFunc) MiddlewareChain {
	var mc MiddlewareChain
	err := appCtx.Find(&mc)
	if err == nil {
		mc = append(mc, middlewares...)
	} else {
		mc = middlewares
	}

	appCtx.Provide(&mc)
	return mc
}
