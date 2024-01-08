package healthcheck

import (
	"runtime"

	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/constantx"
)

func SetupHealthz(appCtx appx.AppContext, router *gin.Engine) {
	// aliveness API
	router.GET(constantx.HealthcheckPrefix+appCtx.WhoAmI()+"/healthz", func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		c.JSON(200, gin.H{
			// application info
			"app":      appCtx.WhoAmI(),
			"runLevel": appCtx.GetRunLevel(),

			// golang info
			"goroutines": runtime.NumGoroutine(),
			"cgoCalls":   runtime.NumCgoCall(),
			"memStats":   m,
		})
	})

}
