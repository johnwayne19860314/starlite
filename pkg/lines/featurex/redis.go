package featurex

import (
	"fmt"
	"time"

	"github.com/mediocregopher/radix/v3"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
)

type Redis typesx.RedisConfig

func (cfg Redis) MustResolve(appCtx appx.AppContext) *radix.Pool {
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Port == 0 {
		cfg.Port = 6379
	}
	if cfg.PoolSize == 0 {
		cfg.PoolSize = 10
	}

	dialOpts := []radix.DialOpt{
		radix.DialReadTimeout(10 * time.Second),
		radix.DialWriteTimeout(10 * time.Second),
	}
	if cfg.Secret != "" {
		dialOpts = append(dialOpts, radix.DialAuthPass(cfg.Secret))
	}
	if cfg.DB != 0 {
		dialOpts = append(dialOpts, radix.DialSelectDB(cfg.DB))
	}
	cf := radix.DefaultConnFunc
	if len(dialOpts) > 0 {
		cf = func(network, addr string) (radix.Conn, error) {
			return radix.Dial(network, addr, dialOpts...)
		}
	}

	clt, err := radix.NewPool("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), cfg.PoolSize, radix.PoolConnFunc(cf))
	if err != nil {
		appCtx.Fatal(err.Error())
	}

	return clt
}

func NewRedis(appCtx appx.AppContext, cl *ConfigLoader) *radix.Pool {
	redisCfg := &Redis{}
	cl.Load(redisCfg)

	return redisCfg.MustResolve(appCtx)
}
