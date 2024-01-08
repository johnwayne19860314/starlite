package lines

import (
	"fmt"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/lines/gracetermx"
	"github.startlite.cn/itapp/startlite/pkg/lines/inject"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/lines/routinex"
)

type appContext struct {
	ready bool // indicates if app is ready for use
	*featurex.Server
	gracetermx.GraceTerm
	router *gin.Engine
	inject.Injector
	//context.Context
}

func (appCtx *appContext) GetInjector() inject.Injector {
	return appCtx.Injector
}

// Quit exits current application
func (appCtx *appContext) Quit(exitCode int) {
	logx.Info(fmt.Sprintf("%s quit with code %d.", appCtx.WhoAmI(), exitCode))
	// @todo should use context.context to do some clean up work
	os.Exit(exitCode)
}

// Fatal logs error message then quits the application
func (appCtx *appContext) Fatal(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	logx.Error(s)
	panic(s)
}

func (appCtx *appContext) GetReady() {
	logx.Info(appCtx.WhoAmI()+" is ready", "runlevel", appCtx.GetRunLevel(), "module", appCtx.WhoAmI())
	appCtx.ready = true
}

func (appCtx *appContext) IsReady() bool {
	return appCtx.ready
}

func (appCtx *appContext) GetRouter() *gin.Engine {
	return appCtx.router
}

func NewAppContext(env string) *appContext {
	appCtx := &appContext{
		GraceTerm: gracetermx.New(syscall.SIGINT, syscall.SIGTERM),
		Injector:  inject.New(),
		router:    gin.New(),
		Server:    featurex.MustProvideServer(env),
	}
	routinex.GoSafe(func() {
		<-appCtx.Done()
		appCtx.Quit(0)
	})

	return appCtx
}
