package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/tokenx"
)

func GetTokenFromRequest(ctx *gin.Context, tokenName string) string {
	token := ctx.GetHeader(tokenName)
	if strings.HasPrefix(token, tokenx.HEADERBearer) {
		return strings.Fields(token)[1]
	}

	token = ctx.GetHeader(tokenx.HEADERAuthorization)
	if strings.HasPrefix(token, tokenx.HEADERBearer) {
		return strings.Fields(token)[1]
	}

	token, _ = ctx.Cookie(tokenName)
	return token
}

func GetRolesRequired(reqCtx appx.ReqContext) []string {
	var routes []lines.Route
	if err := reqCtx.Find(&routes); err != nil {
		panic("Auth failed to work, can not find routes.")
	}

	fullPath := reqCtx.Gin().FullPath()

	for _, each := range routes {
		if each.Method == reqCtx.Gin().Request.Method && each.Group+each.Path == fullPath {
			return each.RolesRequired
		}
	}

	return nil
}
