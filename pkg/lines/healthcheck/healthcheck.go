package healthcheck

import (
	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
)

func Init(appCtx appx.AppContext, router *gin.Engine) {
	SetupPing(appCtx, router)
	SetupHealthz(appCtx, router)
	SetupMetrics(appCtx, router)
}
