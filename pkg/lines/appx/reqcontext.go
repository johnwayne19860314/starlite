package appx

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.startlite.cn/itapp/startlite/pkg/lines/inject"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

type ReqContext interface {
	context.Context
	inject.Injector
	Gin() *gin.Context
	GetXid() string
	GetLogger() logx.ReqLogger
	logx.ReqLogger
}
