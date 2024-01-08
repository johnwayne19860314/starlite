package features

import (
	"fmt"

	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

type kafkaInfoLogger struct {
	Logger logx.ReqLogger
}

type kafkaErrorLogger struct {
	Logger logx.ReqLogger
}

func (l *kafkaInfoLogger) Printf(format string, i ...interface{}) {
	l.Logger.Info(fmt.Sprint(format, i))
}

func (l *kafkaErrorLogger) Printf(format string, i ...interface{}) {
	l.Logger.Error(fmt.Sprint(format, i))
}
