package types

import (
	"github.com/riferrei/srclient"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
)

type Consumer interface {
	Handler(key, value []byte, ctx appx.ReqContext) error
	InitSchema(topicName string, client *srclient.SchemaRegistryClient)
}
