package tchconf

import (
	"github.startlite.cn/itapp/startlite/internal/pkg/types"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
)

type TCHConf types.TCHConfig

func NewTCHConf(cl *featurex.ConfigLoader) *TCHConf {
	tchConf := &TCHConf{}
	cl.Load(tchConf)
	return tchConf
}
