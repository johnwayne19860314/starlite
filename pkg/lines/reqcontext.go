package lines

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/inject"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/lines/xidx"
)

type (
	reqContext struct {
		context.Context
		xid string
		inject.Injector
		GinContext *gin.Context
		logx.ReqLogger
	}
)

func (r reqContext) GetLogger() logx.ReqLogger {
	return r.ReqLogger
}

func (r reqContext) GetXid() string {
	return r.xid
}

func (r *reqContext) SetOrGenerateXid(xid string) string {
	if xid == "" {
		xid = xidx.GenXid()
	}
	r.xid = xid
	return xid
}

func (r reqContext) Gin() *gin.Context {
	return r.GinContext
}

func NewReqContext(appCtx appx.AppContext, c *gin.Context) appx.ReqContext {
	rc := &reqContext{
		Context:    context.Background(),
		Injector:   inject.New(),
		GinContext: c,
	}
	// xid
	xid := ""
	if c.Request != nil {
		xid = xidx.ExtractXidFromRequest(c.Request)
	}
	rc.SetOrGenerateXid(xid)
	rc.ReqLogger = logx.MustNewReqLogger(rc.GetXid())
	rc.ProvideAs(rc.ReqLogger, (*logx.ReqLogger)(nil))
	// link to parent injector
	rc.SetParent(appCtx.GetInjector())
	rc.ProvideAs(rc, (*appx.ReqContext)(nil))
	return rc
}

func NewMockReqContext(appCtx appx.AppContext, path string) appx.ReqContext {
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginCtx.Request, _ = http.NewRequestWithContext(context.Background(), http.MethodGet, path, strings.NewReader(""))
	ginCtx.Request.Header = make(http.Header)
	return NewReqContext(appCtx, ginCtx)
}
