package govconf

import (
	"github.startlite.cn/itapp/startlite/internal/pkg/types"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
)

type GOVConf types.GOVConfig

func NewGOVConf(cl *featurex.ConfigLoader) *GOVConf {
	govConf := &GOVConf{}
	cl.Load(govConf)
	return govConf
}
