package lines

import (
	"fmt"
	"net/http"
	"os"

	_ "go.uber.org/automaxprocs"

	"github.com/gin-gonic/gin"

	"github.com/gin-gonic/contrib/static"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/lines/healthcheck"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/lines/middlewarex"
	"github.startlite.cn/itapp/startlite/pkg/lines/runlevelx"
	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
)

// InitApp will initialize app context with IoC container inject.Injector
func InitApp() *appContext {
	appCtx := NewAppContext("local")

	logx.Init(appCtx.IsLocal())

	//router := gin.New()
	healthcheck.Init(appCtx, appCtx.router)
	appCtx.router.Use(static.Serve("/", static.LocalFile("../../ui/dist", true)))
	appCtx.router.Use(MiddlewareWrapper(appCtx, middlewarex.Log))
	appCtx.router.Use(MiddlewareWrapper(appCtx, middlewarex.Recovery))
	SetupCORS(appCtx, appCtx.router)
	return appCtx
}

func InitTestApp(appName string) (appx.AppContext, *gin.Engine) {
	appCtx := NewAppContext("local")
	//featurex.Provide(appCtx)
	logx.Init(true)
	router := gin.New()
	appCtx.Provide(router)

	serverConfig := &featurex.Server{
		ServerConfig: typesx.ServerConfig{
			Name:     appName,
			Port:     0,
			RunLevel: "testing",
			Host:     "",
		},
		RunLevel: runlevelx.RunLevelLocal,
	}
	appCtx.Provide(serverConfig)
	err := appCtx.Apply(appCtx)
	if err != nil {
		appCtx.Fatal("can't init test app context: %s", err)
	}
	return appCtx, router
}

func StartServer(appCtx *appContext) {
	server := appCtx.Server
	host := fmt.Sprintf(":%d", server.Port)
	if server.IsLocal() {
		host = "127.0.0.1" + host
	}

	srv := &http.Server{
		Addr:    host,
		Handler: appCtx.GetRouter(),
	}

	logx.Info(appCtx.WhoAmI()+" is listening on", "host", host)
	err := srv.ListenAndServe()
	if err != nil {
		appCtx.Fatal(err.Error())
	}
}
