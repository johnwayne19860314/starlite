package features

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"time"

	_ "github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"github.startlite.cn/itapp/startlite/pkg/lines"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/lines/routinex"
)

func MustNewConsumer(appCtx appx.AppContext, cl *featurex.ConfigLoader, handler func(key, value []byte, ctx appx.ReqContext) error) error {

	config := Kafka{}
	cl.Load(&config)
	consumerConfig := Consumer{}
	cl.Load(&consumerConfig)

	routinex.GoSafe(func() {
		for {
			routinex.RunSafe(func() {
				_ = consumerLoop(appCtx, &consumerConfig, &config, handler)
			})
		}
	})

	return nil
}

func consumerLoop(ctx appx.AppContext, consumer *Consumer, broker *Kafka, handler func(key, value []byte, reqCtx appx.ReqContext) error) error {

	reqCtx1 := lines.NewMockReqContext(ctx, "/consumerChargeSession")
	logger := &kafkaInfoLogger{Logger: reqCtx1.GetLogger()}
	errorLogger := &kafkaErrorLogger{Logger: reqCtx1.GetLogger()}
	bootstrapservers := broker.Addrs
	kafkaConfig := kafka.ReaderConfig{
		Brokers: bootstrapservers,
		GroupID: consumer.GroupID,
		//Topic: sub.Topics[0],
		GroupTopics:    consumer.Topics,
		Logger:         logger,
		ErrorLogger:    errorLogger,
		MinBytes:       10e3, // 1K
		MaxBytes:       10e5, // 100k
		MaxAttempts:    10000,
		CommitInterval: time.Second * 5,
	}
	if broker.TLSEnable {
		certAndKey, err := tls.X509KeyPair([]byte(broker.ClientCert), []byte(broker.ClientKey))
		if err != nil {
			return errors.New("load consumer cert and key failed")
		}
		caBytes := []byte(broker.CaCert)
		if err != nil {
			return errors.New("load consumer ca failed")
		}
		clientCertPool := x509.NewCertPool()
		ok := clientCertPool.AppendCertsFromPEM(caBytes)
		if !ok {
			return errorx.New("append root ca failed")
		}
		tlsConfig := tls.Config{RootCAs: clientCertPool, Certificates: []tls.Certificate{certAndKey}}
		dialer := kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			TLS:       &tlsConfig,
		}
		kafkaConfig.Dialer = &dialer
	}

	r := kafka.NewReader(kafkaConfig)
	for {
		reqCtx := lines.NewMockReqContext(ctx, "/consumerChargeSession")
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		_ = handler(m.Key, m.Value, reqCtx)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	return nil
}
