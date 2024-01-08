package lines

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/bizerr"
	"github.startlite.cn/itapp/startlite/pkg/lines/constantx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/middlewarex"
	"github.startlite.cn/itapp/startlite/pkg/lines/xidx"
)

const requestErr = "request end up with error"

type (
	EntryPoint func(appx.ReqContext) (interface{}, error)

	RawResponse struct {
		StatusCode    int
		ContentLength int64
		ContentType   string
		Reader        io.Reader
		ExtraHeaders  map[string]string
	}
)

func NewRawResponse(contentType string, reader io.Reader) *RawResponse {
	return &RawResponse{
		StatusCode:    -1,
		ContentLength: -1,
		ContentType:   contentType,
		Reader:        reader,
		ExtraHeaders:  make(map[string]string),
	}
}

func SetupCORS(appCtx appx.AppContext, router *gin.Engine) {
	cfg := cors.Config{
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "token", "authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
		MaxAge:           12 * time.Hour,
	}
	if appCtx.IsLocal() {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = []string{"*.xxx.cn", "*.xxx.com", "*.xxxmotors.com"}
	}

	// cors - put it before router handle
	router.Use(cors.New(cfg))
}

func UseMiddlewareChain(appCtx appx.AppContext, router *gin.Engine) {
	var mc middlewarex.MiddlewareChain
	err := appCtx.Find(&mc)
	if err == nil {
		router.Use(mc...)
	}
}

func SetupHttpServer(appCtx appx.AppContext, routes []Route) {
	//routes = ExtendRoutes(appCtx, routes)

	for _, route := range routes {
		if route.Group == "" {
			route.Group = constantx.APIPrefix + appCtx.WhoAmI()
		}
		if route.CustomHandlerFunc != nil {
			appCtx.GetRouter().Handle(route.Method, route.MakeRelativePath(), route.CustomHandlerFunc)
		} else {
			appCtx.GetRouter().Handle(route.Method, route.MakeRelativePath(), HttpHandler(appCtx, route.Handler))
		}
	}
}

// HttpHandler help setup ReqContext for controller
func HttpHandler(appCtx appx.AppContext, point EntryPoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rc := getOrCreateReqContext(appCtx, ctx)
		ret, err := point(rc)
		ctx.Writer.Header().Set(xidx.HeaderXid, rc.GetXid())
		if err != nil {
			var httpCode int
			if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
				httpCode = http.StatusNotFound
				ret = map[string]interface{}{
					"message":  err.Error(),
					"xid":      rc.GetXid(),
					"response": ret,
				}
			} else if bizErr, ok := err.(bizerr.BizLogicError); ok {
				httpCode = bizErr.HttpStatusCode()
				ret = map[string]interface{}{
					"code":     bizErr.ErrorNo(),
					"message":  bizErr.Message(),
					"xid":      rc.GetXid(),
					"response": ret,
				}
			} else {
				httpCode = http.StatusInternalServerError
				ret = map[string]interface{}{
					"message":  err.Error(),
					"xid":      rc.GetXid(),
					"response": ret,
				}
			}

			if _, ok := err.(bizerr.BizLogicError); ok {
				rc.Error(fmt.Sprintf("%s: \n%+v\n", requestErr, err))
			} else {
				if _, ok := err.(errorx.StackTraceWrapper); ok {
					rc.Error(fmt.Sprintf("%+v\n", err))
				} else {
					rc.Error(requestErr, "error", err)
				}
			}

			ctx.AbortWithStatusJSON(httpCode, ret)
			return
		}

		if r, ok := ret.(*RawResponse); ok {
			statusCode := http.StatusOK
			if r.StatusCode != -1 {
				statusCode = r.StatusCode
			}
			ctx.DataFromReader(statusCode, r.ContentLength, r.ContentType, r.Reader, r.ExtraHeaders)
			return
		}
		ret = map[string]interface{}{
			"message":  "ok",
			"xid":      rc.GetXid(),
			"data": ret,
		}
		ctx.JSON(http.StatusOK, ret)
	}
}
