package appx

import (
	"github.com/gin-gonic/gin"
	"github.startlite.cn/itapp/startlite/pkg/lines/gracetermx"
	"github.startlite.cn/itapp/startlite/pkg/lines/inject"
)

type (
	AppContext interface {
		// Server Feature
		WhoAmI() string
		GetRunLevel() string
		Host() string
		APIEndpoint() string
		IsTesting() bool
		IsLocal() bool
		IsPrd() bool

		GetReady()
		IsReady() bool
		GetRouter() *gin.Engine
		Quit(exitCode int)
		Fatal(format string, args ...interface{})

		gracetermx.GraceTerm

		GetInjector() inject.Injector
		inject.Injector
		
	}
)
