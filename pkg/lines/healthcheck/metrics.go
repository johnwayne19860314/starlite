package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/constantx"
)

func SetupMetrics(appCtx appx.AppContext, router *gin.Engine) {
	router.GET(constantx.HealthcheckPrefix+appCtx.WhoAmI()+"/metrics", gin.WrapH(promhttp.Handler()))
}
