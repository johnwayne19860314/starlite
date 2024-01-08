package features

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/xidx"
)

func NewMockReqContext(appCtx appx.AppContext, xid string) appx.ReqContext {
	return lines.NewReqContext(appCtx, NewMockGinContext(xid))
}

func NewMockGinContext(xid string) *gin.Context {
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginCtx.Request, _ = http.NewRequestWithContext(context.Background(), http.MethodGet, "/", strings.NewReader(""))
	ginCtx.Request.Header = make(http.Header)
	ginCtx.Request.Header.Set(xidx.HeaderXid, xid)

	return ginCtx
}
