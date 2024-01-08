package features

import (
	goconfluence "github.com/virtomize/confluence-go-api"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/servicex/types"
)

type Confluence struct {
	types.ConfluenceConfig
	*goconfluence.API
}

func NewConfluence(appCtx appx.AppContext, cl *featurex.ConfigLoader) (*Confluence, error) {
	res := &Confluence{}
	cl.Load(res)

	api, err := goconfluence.NewAPI(res.Host, res.Username, res.Password)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	res.API = api

	return res, nil
}

func MustNewConfluence(appCtx appx.AppContext, cl *featurex.ConfigLoader) *Confluence {
	res, err := NewConfluence(appCtx, cl)
	if err != nil {
		appCtx.Fatal("MustNewConfluence err : %v", err)
	}

	return res
}
