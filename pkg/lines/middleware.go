package lines

import (
	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
)

const (
	reqContextKey = "__REQ_CONTEXT_KEY__"
)

type (
	Middleware func(appCtx appx.AppContext, reqCtx appx.ReqContext) error
)

func MiddlewareWrapper(appCtx appx.AppContext, mid Middleware) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqCtx := getOrCreateReqContext(appCtx, ctx)
		mid(appCtx, reqCtx)
		// if err := mid(appCtx,reqCtx); err != nil {
		// 	logx.Error("error when calling middleware", "error", err)
		// }
	}
}

func getOrCreateReqContext(appCtx appx.AppContext, ctx *gin.Context) appx.ReqContext {
	val, exists := ctx.Get(reqContextKey)
	if exists {
		if reqCtx, ok := val.(appx.ReqContext); ok {
			return reqCtx
		}
	}
	rc := NewReqContext(appCtx, ctx)
	ctx.Set(reqContextKey, rc)
	return rc
}
