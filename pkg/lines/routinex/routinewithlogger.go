package routinex

import (
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

func RunSafeWithLogger(fn func(), logger logx.ReqLogger) (panicked error) {
	defer func() {
		if p := recover(); p != nil {
			logger.Error("recover from panic", "info", p)
			panicked = errorx.EnsureStack(p, 1)
		}
	}()

	fn()
	return nil
}
