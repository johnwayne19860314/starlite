package featurex

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/inject"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/lines/routinex"
	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
	"github.startlite.cn/itapp/startlite/pkg/lines/xidx"
)

// MustStartTwpConsumer start twp consumer
//
//	func errorHandle(err error) {
//		fmt.Println(err)
//	}
//
//	func messageHandle(event *featurex.WorkflowEvent) {
//		fmt.Println("handle", event.Data.ProcessInstanceId)
//	}
//
//	func aroundMessageHandle(reqCtx appx.ReqContext, chain featurex.HandleChain, event *featurex.WorkflowEvent) error {
//		fmt.Println("before handle", event.Data.ProcessDefinitionId)
//		if err := chain.Next(reqCtx); err != nil {
//			return err
//		}
//		fmt.Println("after handle", event.Data.ProcessDefinitionId)
//		return nil
//	}
//
//	func skipMessageHandle(chain featurex.HandleChain, event *featurex.WorkflowEvent) {
//		if time.Now().Unix()%2 == 0 {
//			fmt.Println(event.Data.ProcessInstanceId, "skipped")
//			chain.Complete()
//		}
//	}
//
//	func allMessageHandle(event *featurex.WorkflowEvent) {
//		fmt.Println("receive", event.Data.ProcessInstanceId)
//	}
//
//	func main() {
//		featurex.OnTwpError(errorHandle)
//		featurex.OnTwpMessage(allMessageHandle, skipMessageHandle, aroundMessageHandle, messageHandle)
//		//featurex.OnTwpMessagePrepend(messagePrependHandle)
//		appCtx := lines.NewAppContext()
//		cl := featurex.NewConfigLoader()
//		appCtx.Provide(cl)
//		appCtx.Invoke(featurex.MustStartTwpConsumer)
//		select {}
//	}
func MustStartTwpConsumer(appCtx appx.AppContext, cl *ConfigLoader) {
	twp := NewTwpConsumer(cl)
	appCtx.Provide(twp)
	if err := twp.Start(appCtx); err != nil {
		panic(err)
	}
}

type TwpConsumer struct {
	typesx.TwpConsumerConfig

	block          cipher.Block
	handles        []interface{}
	prependHandles []interface{}
	errorHandles   []interface{}
	appSet         map[string]bool
	reader         *kafka.Reader
}

var defaultTwpErrorHandles []interface{}
var defaultTwpMessageHandles []interface{}
var defaultTwpMessagePrependHandles []interface{}

func OnTwpError(handles ...interface{}) {
	for _, handle := range handles {
		if reflect.TypeOf(handle).Kind() != reflect.Func {
			panic("OnTwpError handle must be a function")
		}
		defaultTwpErrorHandles = append(defaultTwpErrorHandles, handle)
	}

}

func OnTwpMessage(handles ...interface{}) {
	for _, handle := range handles {
		if reflect.TypeOf(handle).Kind() != reflect.Func {
			panic("OnTwpMessage handle must be a function")
		}
		defaultTwpMessageHandles = append(defaultTwpMessageHandles, handle)
	}
}

func OnTwpMessagePrepend(handles ...interface{}) {
	for _, handle := range handles {
		if reflect.TypeOf(handle).Kind() != reflect.Func {
			panic("OnTwpMessagePrepend handle must be a function")
		}
		defaultTwpMessagePrependHandles = append(defaultTwpMessagePrependHandles, handle)
	}
}
func NewTwpConsumer(cl *ConfigLoader) *TwpConsumer {
	var consumer TwpConsumer
	cl.Load(&consumer)
	consumer.init()
	return &consumer
}

func NewTwpConsumerFromConfig(cfg typesx.TwpConsumerConfig) *TwpConsumer {
	consumer := &TwpConsumer{
		TwpConsumerConfig: cfg,
	}
	consumer.init()
	return consumer
}

func (consumer *TwpConsumer) init() {
	if consumer.Secret != "" {
		key := []byte(consumer.Secret)
		sumKey := sha1.Sum(key)
		block, err := aes.NewCipher(sumKey[:16])
		if err != nil {
			panic(err)
		}
		consumer.block = block
	}

	consumer.appSet = make(map[string]bool)
	if consumer.Apps != "" {
		apps := strings.Split(consumer.Apps, ",")
		for _, app := range apps {
			consumer.appSet[strings.TrimSpace(app)] = true
		}
	} else {
		consumer.appSet["*"] = true
	}
}

func (consumer *TwpConsumer) Start(appCtx appx.AppContext) error {
	reqCtx := newKafkaReqContext(appCtx)
	kafkaConfig := kafka.ReaderConfig{
		Brokers: []string{consumer.BoostrapServer},

		Logger: kafka.LoggerFunc(func(s string, i ...interface{}) {
			reqCtx.Info(fmt.Sprintf(s, i...))
		}),
		ErrorLogger: kafka.LoggerFunc(func(s string, i ...interface{}) {
			reqCtx.Error(fmt.Sprintf(s, i...))
		}),
		MinBytes:       10e3, // 1K
		MaxBytes:       10e5, // 100k
		MaxAttempts:    10000,
		CommitInterval: time.Second * 5,
	}
	if consumer.ConsumerGroup != "" {
		kafkaConfig.GroupID = consumer.ConsumerGroup
		kafkaConfig.GroupTopics = []string{consumer.Topics}
	} else {
		kafkaConfig.Topic = consumer.Topics
		kafkaConfig.Partition = consumer.Partition
	}

	if consumer.TLSEnable {
		certAndKey, err := tls.X509KeyPair([]byte(consumer.TLSCert), []byte(consumer.TLSKey))
		if err != nil {
			return errorx.New("load consumer cert and key failed")
		}
		caBytes := []byte(consumer.TLSCA)
		if err != nil {
			return errorx.New("load consumer ca failed")
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

	reader := kafka.NewReader(kafkaConfig)
	consumer.reader = reader
	appCtx.Provide(consumer)

	var errorHandles []interface{}
	var messageHandles []interface{}

	errorHandles = append(errorHandles, defaultTwpErrorHandles...)
	errorHandles = append(errorHandles, consumer.errorHandles...)
	messageHandles = append(messageHandles, defaultTwpMessagePrependHandles...)
	messageHandles = append(messageHandles, consumer.prependHandles...)
	messageHandles = append(messageHandles, valueHandle, appFilterHandle, decryptHandle)
	messageHandles = append(messageHandles, defaultTwpMessageHandles...)
	messageHandles = append(messageHandles, consumer.handles...)

	routinex.GoSafe(func() {
		defer func() {
			if err := consumer.Close(); err != nil {
				reqCtx.Error("Kafka Close", "error", err)
			}
		}()

		for {
			reqCtx := newKafkaReqContext(appCtx)
			message, err := reader.ReadMessage(context.Background())
			if err == io.EOF {
				break
			}

			if err != nil {
				reqCtx.Error("Kafka ReadMessage", "error", err)
				continue
			}

			handler := messageHandler{
				index:   0,
				handles: messageHandles,
			}

			if err := handler.handle(reqCtx, &message); err != nil {
				if len(errorHandles) > 0 {
					reqCtx.ProvideAs(errorx.WithStack(err), (*error)(nil))
					errorHandler := messageHandler{
						index:   0,
						handles: errorHandles,
					}
					err = errorHandler.handle(reqCtx, &message)
				}
				if err != nil {
					reqCtx.Error("Kafka handle Message", "error", err)
				}
				continue
			}
		}
	})

	return nil
}

func (consumer *TwpConsumer) Close() error {
	if consumer.reader != nil {
		if err := consumer.reader.Close(); err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}

func (consumer *TwpConsumer) OnMessage(handles ...interface{}) {
	consumer.handles = nil
	for _, handle := range handles {
		if reflect.TypeOf(handle).Kind() != reflect.Func {
			panic("OnMessage handle must be a function")
		}
		consumer.handles = append(consumer.handles, handle)
	}
}

func (consumer *TwpConsumer) OnMessagePrepend(handles ...interface{}) {
	consumer.prependHandles = nil
	for _, handle := range handles {
		if reflect.TypeOf(handle).Kind() != reflect.Func {
			panic("OnMessage handle must be a function")
		}
		consumer.prependHandles = append(consumer.prependHandles, handle)
	}
}

func (consumer *TwpConsumer) OnError(handles ...interface{}) {
	for _, handle := range handles {
		if reflect.TypeOf(handle).Kind() != reflect.Func {
			panic("OnError handle must be a function")
		}
		consumer.errorHandles = append(consumer.errorHandles, handle)
	}
}

func valueHandle(reqCtx appx.ReqContext, msg *kafka.Message) error {
	event := WorkflowEvent{}
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return errorx.WithStack(err)
	}
	reqCtx.Provide(&event)
	return nil
}

func appFilterHandle(event *WorkflowEvent, handleChain HandleChain, consumer *TwpConsumer) error {
	process := event.Data.ProcessDefinitionKey
	if process == "" {
		process = strings.Split(event.Data.ProcessDefinitionId, ":")[0]
	}
	if consumer.appSet["*"] {
		return nil
	}

	if !consumer.appSet[process] {
		handleChain.Complete()
	}
	return nil
}

func decryptHandle(event *WorkflowEvent, consumer *TwpConsumer) error {
	var err error
	if event.Data != nil && consumer.block != nil {
		if event.Data.Variables, err = decryptVariables(consumer.block, event.Data.Variables); err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}

func decryptVariables(cipher cipher.Block, variables map[string]interface{}) (map[string]interface{}, error) {
	newVariables := make(map[string]interface{})
	for k, v := range variables {
		if str, ok := v.(string); ok {
			newStr, err := doDecrypt(cipher, []byte(str))
			if err != nil {
				return nil, errorx.WithStack(err)
			}
			var newObj interface{}
			if err = json.Unmarshal(newStr, &newObj); err != nil {
				return nil, errorx.WithStack(err)
			}
			newVariables[k] = newObj
		} else {
			newVariables[k] = v
		}
	}
	return newVariables, nil
}

func doDecrypt(block cipher.Block, source []byte) ([]byte, error) {
	target := make([]byte, base64.StdEncoding.DecodedLen(len(source)))
	n, err := base64.StdEncoding.Decode(target, source)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	target = target[:n]
	decrypt, err := aesDecryptECB(block, target)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return decrypt, nil
}

func aesDecryptECB(block cipher.Block, encrypted []byte) (decrypted []byte, err error) {
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, block.BlockSize(); bs < len(encrypted); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim], nil
}

type kafkaReqContext struct {
	context.Context
	inject.Injector
	xid string
	logx.ReqLogger
}

func (k *kafkaReqContext) Gin() *gin.Context {
	return nil
}

func (k *kafkaReqContext) GetXid() string {
	return k.xid
}

func (k *kafkaReqContext) GetLogger() logx.ReqLogger {
	return k.ReqLogger
}

// newKafkaReqContext avoid import cycle
func newKafkaReqContext(appContext appx.AppContext) appx.ReqContext {
	k := &kafkaReqContext{
		Context:  context.Background(),
		Injector: inject.New(),
	}
	k.xid = xidx.GenXid()
	k.ReqLogger = logx.MustNewReqLogger(k.xid)
	k.ProvideAs(k.ReqLogger, (*logx.ReqLogger)(nil))
	k.SetParent(appContext.GetInjector())
	k.ProvideAs(k, (*appx.ReqContext)(nil))
	return k
}

var _ appx.ReqContext = &kafkaReqContext{}

type WorkflowEvent struct {
	Producer  string        `json:"producer"`
	Event     string        `json:"event"`
	Xid       string        `json:"xid"`
	Timestamp string        `json:"timestamp"`
	Data      *WorkflowData `json:"data"`
}

type WorkflowData struct {
	// TaskData
	ProcessDefinitionId string                 `json:"processDefinitionId"`
	ProcessInstanceId   string                 `json:"processInstanceId"`
	FormKey             string                 `json:"formKey"`
	Priority            int                    `json:"priority"`
	TaskDefinitionKey   string                 `json:"taskDefinitionKey"`
	ExecutionId         string                 `json:"executionId"`
	CreateTime          string                 `json:"createTime"` // "Tue Oct 11 12:16:25 CST 2022"
	ClaimTime           string                 `json:"claimTime"`
	EndTime             string                 `json:"endTime"` // ??
	DueDate             string                 `json:"dueDate"`
	DelegationState     string                 `json:"delegationState"`
	CandidateGroups     []string               `json:"candidateGroups" copier:"-"`
	CandidateUsers      []string               `json:"candidateUsers" copier:"-"`
	Participants        []string               `json:"participants" copier:"-"`
	BusinessKey         string                 `json:"businessKey"`
	TaskId              string                 `json:"taskId"`
	Owner               string                 `json:"owner"`
	Outcome             string                 `json:"outcome"`
	Assignee            string                 `json:"assignee"`
	Variables           map[string]interface{} `json:"variables" copier:"-"`

	// ProcessData
	//ProcessInstanceId    string                 `json:"processInstanceId"`
	//Variables            map[string]interface{} `json:"variables"`
	StartUserId string `json:"startUserId"`
	//BusinessKey          string                 `json:"businessKey"`
	StartTime            string `json:"startTime"`
	ProcessDefinitionKey string `json:"processDefinitionKey"`
	//Participants         []string               `json:"participants"`
	//Outcome              string                 `json:"outcome"`
}

type messageHandler struct {
	index   int
	handles []interface{}
}

type HandleChain interface {
	Complete()
	Next(reqContext appx.ReqContext) error
}

var _ HandleChain = &messageHandler{}

func (messageHandler *messageHandler) Complete() {
	messageHandler.index = len(messageHandler.handles)
}

func (messageHandler *messageHandler) handle(reqCtx appx.ReqContext, msg *kafka.Message) error {
	reqCtx.ProvideAs(messageHandler, (*HandleChain)(nil))
	reqCtx.Provide(msg)
	if err := messageHandler.Next(reqCtx); err != nil {
		return errorx.WithStack(err)
	}
	return nil
}

func (messageHandler *messageHandler) Next(reqCtx appx.ReqContext) error {
	for messageHandler.index < len(messageHandler.handles) {
		handle := messageHandler.handles[messageHandler.index]
		messageHandler.index++
		ret, err := reqCtx.Invoke(handle)
		if err != nil {
			return errorx.WithStack(err)
		}

		if len(ret) > 0 {
			err, ok := ret[len(ret)-1].Interface().(error)
			if ok && err != nil {
				return errorx.WithStack(err)
			}
		}
	}
	return nil

}
