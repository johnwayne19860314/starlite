package features

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/servicex/types"
)

type Kafka types.KafkaBrokerConfig
type Consumer types.KafkaConsumerConfig
type Producer types.KafkaProducerConfig

type KafkaProducer struct {
	producer kafka.Writer
	topic    string
}

/*
*
in current producer, no topic assigned when the producer created, topic will be assigned when sending the message.
*/
func MustInitKafkaProducer(appCtx appx.AppContext, cl *featurex.ConfigLoader) *KafkaProducer {
	config := Consumer{}
	cl.Load(&config)
	producerConfig := Producer{}
	cl.Load(&producerConfig)

	p := KafkaProducer{}
	p.producer = kafka.Writer{Addr: kafka.TCP(config.KafkaBrokerConfig.Addrs...)}
	p.topic = producerConfig.Topic
	return &p
}

func (p *KafkaProducer) ProduceSingle(req appx.ReqContext, value interface{}) error {
	message := kafka.Message{Key: []byte(req.GetXid())}

	vb, err := json.Marshal(value)
	if err != nil {
		return errorx.WithStack(err)
	}

	message.Value = vb
	message.Topic = p.topic
	err = p.producer.WriteMessages(context.Background(), message)
	if err != nil {
		req.Error("failed to write messages", "error", err, "topic", value)
	}

	retryCount := 0
	for err != nil && retryCount < 3 {
		err = p.producer.WriteMessages(context.Background(), message)
		if err != nil {
			req.Error("failed to write messages in retry", "error", err, "topic", value, "retryCount", retryCount)
		}
		retryCount++
	}
	return err
}
