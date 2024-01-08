package features

import (
	"time"

	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

type ActionTrier struct {
	Err error

	Fn        func() error
	ReqLogger logx.ReqLogger

	TryDuration  []time.Duration
	CurrentIndex int
}

func (at *ActionTrier) Do() {
	err := at.Fn()
	at.Err = err
	if err == nil {
		return
	}

	at.ReqLogger.Warn("ActionTrier do err", "err", err, "try_index", -1)

	for index, v := range at.TryDuration {
		at.CurrentIndex = index
		time.Sleep(v)

		err = at.Fn()
		at.Err = err
		if err == nil {
			return
		}
		at.ReqLogger.Warn("ActionTrier do err", "err", err, "try_index", index)
	}

}
