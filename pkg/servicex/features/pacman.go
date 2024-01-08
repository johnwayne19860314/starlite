package features

import (
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/servicex/types"
)

type Pacman types.PacmanConfig

func NewPacmanConfig(cl *featurex.ConfigLoader) *Pacman {
	pacman := &Pacman{}
	cl.Load(pacman)
	return pacman
}
