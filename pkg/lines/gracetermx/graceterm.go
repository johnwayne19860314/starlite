package gracetermx

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/lines/routinex"
	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
)

type (
	GraceTerm interface {
		BeforeExit(fn func())
		Done() <-chan typesx.PlaceHolder
		StoppingSignal() <-chan typesx.PlaceHolder
	}

	graceTerm struct {
		*sync.RWMutex
		cleaningUp   chan typesx.PlaceHolder
		cleanups     []func()
		exitSignal   chan typesx.PlaceHolder
		notifySignal chan os.Signal
	}
)

func New(signals ...os.Signal) GraceTerm {
	gt := &graceTerm{
		RWMutex:      &sync.RWMutex{},
		cleaningUp:   make(chan typesx.PlaceHolder),
		notifySignal: make(chan os.Signal, 1),
		exitSignal:   make(chan typesx.PlaceHolder),
	}
	signal.Notify(gt.notifySignal, signals...)

	routinex.GoSafe(func() {
		sig := <-gt.notifySignal
		logx.Info("got term signal, cleaning up", "signal", sig)
		close(gt.cleaningUp)
		gt.cleanUp()
	})
	return gt
}

func (gt graceTerm) StoppingSignal() <-chan typesx.PlaceHolder {
	return gt.cleaningUp
}

func (gt graceTerm) cleanUp() {
	gt.Lock()
	defer gt.Unlock()

	var wg sync.WaitGroup
	for _, fn := range gt.cleanups {
		wg.Add(1)
		routinex.GoSafe(func() {
			fn()
			wg.Done()
		})
	}
	wg.Wait()
	logx.Info("all cleanup done, send done signal")
	err := logx.Logger().Sync()
	if err != nil {
		fmt.Println(err.Error())
	}
	close(gt.exitSignal)
}

func (gt graceTerm) Done() <-chan typesx.PlaceHolder {
	return gt.exitSignal
}

func (gt *graceTerm) BeforeExit(fn func()) {
	gt.Lock()
	defer gt.Unlock()
	gt.cleanups = append(gt.cleanups, fn)
}
