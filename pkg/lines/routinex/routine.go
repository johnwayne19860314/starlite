package routinex

import (
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
)

func GoRestartUnlessStopped(appStopSignal <-chan typesx.PlaceHolder, fn func()) chan typesx.PlaceHolder {
	stopChan := make(chan typesx.PlaceHolder)
	var runner func()
	runner = func() {
		defer func() {
			if p := recover(); p != nil {
				logx.Error("recover from panic", "info", p)
			}
			select {
			case <-stopChan:
				return
			case <-appStopSignal:
				return
			default:
			}
			go func() {
				runner()
			}()
		}()

		fn()
	}
	go runner()
	return stopChan
}

// GoSafe will always catch panic and log it out, routine will not cause app crash
func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer Recover()

	fn()
}

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logx.Error("recover from panic", "info", p)
	}
}
