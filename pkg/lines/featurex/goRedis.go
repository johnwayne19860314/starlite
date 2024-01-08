package featurex

import (
	"strings"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/rds"
	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
)

type GoRedis typesx.NewRedisConfig

func NewRedisClient(appCtx appx.AppContext, cl *ConfigLoader) *rds.RedisClient {
	redisCfg := GoRedis{}
	cl.Load(&redisCfg)

	conf := rds.RedisConfig{
		Addrs:            strings.Split(redisCfg.Addr, ","),
		Password:         redisCfg.Password,
		DB:               redisCfg.DB,
		PoolSize:         redisCfg.PoolSize,
		MasterName:       redisCfg.MasterName,
		SentinelPassword: redisCfg.SentinelPassword,
		KeyPrefix:        redisCfg.KeyPrefix,
	}

	client, err := rds.NewRedisProvider(&conf)
	if err != nil {
		appCtx.Fatal(err.Error())
	}
	return client
}
