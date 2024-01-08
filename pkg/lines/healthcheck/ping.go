package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/constantx"
	"github.startlite.cn/itapp/startlite/pkg/lines/timex"
	"github.startlite.cn/itapp/startlite/pkg/lines/versionx"
)

func SetupPing(appCtx appx.AppContext, router *gin.Engine) {
	startup := timex.CnNowString()
	router.GET(constantx.HealthcheckPrefix+appCtx.WhoAmI()+"/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":          appCtx.WhoAmI(),
			"runLevel":     appCtx.GetRunLevel(),
			"host":         appCtx.Host(),
			"startup":      startup,
			"version":      versionx.Version,
			"linesVersion": constantx.LinesVersion,
			"gitHash":      versionx.GitHash,
			"gitBranch":    versionx.GitBranch,
			"buildDate":    versionx.BuildDate,
		})
	})
}
