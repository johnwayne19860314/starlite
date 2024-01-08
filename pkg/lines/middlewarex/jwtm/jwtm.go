package jwtm

import (
	"fmt"
	"net/http"

	"github.startlite.cn/itapp/startlite/pkg/lines"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/claimx"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/tokenx"
)

const (
	OAuth2AccessToken = "workflowToken"

	jwtErrClaims = "can't get claims"
)

// JWT will inject claims from token in cookie or header
func JWT(reqCtx appx.ReqContext) {
	var (
		oa            *featurex.OAuth2
		routes        []lines.Route
		requiredRoles []string
	)
	c := reqCtx.Gin()
	if err := reqCtx.Find(&routes); err != nil {
		panic("JWT failed to work, can not find routes.")
	}
	fullPath := c.FullPath()
	// If the RolesRequired is nil, just ignore the JWT checking.
	// If the RolesRequired is empty list, make sure the JWT is valid.
	// If the RolesRequired has any roles, make sure at least one role is included in the JWT's claim.
	for _, each := range routes {
		if each.Method == c.Request.Method && each.Group+each.Path == fullPath {
			if each.RolesRequired == nil {
				c.Next()
				return
			} else {
				requiredRoles = each.RolesRequired
				break
			}
		}
	}
	if requiredRoles == nil {
		c.Next()
		return
	}

	if err := reqCtx.Find(&oa); err != nil {
		panic("JWT failed to work, can not find OAuth.")
	}
	token := tokenx.GetToken(reqCtx)
	if token == "" { // at this point, the API requires roles.
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
			"message": "API need token",
			"xid":     reqCtx.GetXid(),
		})
		return
	}

	claims, err := oa.GetClaims(token)
	if err != nil {
		stackErr := errorx.WithStack(err)
		reqCtx.Error(jwtErrClaims, "error", stackErr)
		reqCtx.Gin().SetCookie(OAuth2AccessToken, "", 0, "/", "", true, true)
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
			"message": fmt.Sprintf("%+v", stackErr),
			"xid":     reqCtx.GetXid(),
		})
		return
	}
	reqCtx.Provide(&claims)
	reqCtx.With("user", claimx.GetClaim(&claims, "email"))

	if len(requiredRoles) == 0 { // only require valid token
		c.Next()
		return
	}
	realmAccess, exist := claims["realm_access"]
	if !exist {
		panic("no realm_access in token")
	}

	rolesInterface, exist := realmAccess.(map[string]interface{})["roles"]
	if !exist {
		panic("no roles in token")
	}
	rolesInterfaceArray, ok := rolesInterface.([]interface{})
	if !ok {
		panic("invalid roles in token")
	}
	// If any of the roles that is defined in the routes is included in the token, let it pass.
	for _, requiredRole := range requiredRoles {
		for _, has := range rolesInterfaceArray {
			if requiredRole == has { // string and interface{string} seems can be compared together in go magically
				c.Next()
				return
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
		"message": fmt.Sprintf("API need role %v", requiredRoles),
		"xid":     reqCtx.GetXid(),
	})
}
