package tokenx

import (
	"fmt"
	"strings"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/textx"
)

const (
	BJMSessionToken     = "BJMSessionToken"
	WorkflowToken       = "workflowToken"
	SSOToken            = "ssoToken"
	HEADERBearer        = "Bearer"
	HEADERAuthorization = "Authorization"
)

type Tokener interface {
	GetToken() string
}

type tokener struct {
	token string
}

func (i *tokener) GetToken() string {
	return i.token
}

func NewTokener(token string) Tokener {
	return &tokener{
		token: token,
	}
}

func GetToken(reqCtx appx.ReqContext) string {
	var i Tokener
	err := reqCtx.Find(&i)
	if err == nil && !textx.Blank(i.GetToken()) {
		return i.GetToken()
	}
	token := ""
	if reqCtx.Gin() != nil && reqCtx.Gin().Request != nil {
		token, err = reqCtx.Gin().Cookie(WorkflowToken)
		if err == nil && !textx.Blank(token) {
			return token
		}
		token, err = reqCtx.Gin().Cookie(SSOToken)
		if err == nil && !textx.Blank(token) {
			return token
		}
		token = reqCtx.Gin().GetHeader(HEADERAuthorization)
		if strings.HasPrefix(token, HEADERBearer) {
			return strings.Fields(token)[1]
		}
	}

	return token
}

func GetBearerHeader(reqCtx appx.ReqContext) string {
	token := GetToken(reqCtx)
	if textx.Blank(token) {
		return ""
	}
	return fmt.Sprintf("%s %s", HEADERBearer, token)
}
