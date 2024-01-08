package clients

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	goprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	utils "github.startlite.cn/itapp/startlite/internal/pkg/infra/util/utils"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

type GrpcClient struct {
	Config GrpcConfig

	ConnMu sync.Mutex
	Conn   *grpc.ClientConn
}

type GrpcConfig struct {
	ServerAddress string
	Cert          tls.Certificate
	Insecure      bool
}

func NewGrpcInsecureConfig(serverAddress string) (*GrpcConfig, error) {
	if serverAddress == "" {
		return nil, errorx.Errorf("NewGrpcInsecureConfig: missing server address")
	}
	return &GrpcConfig{
		ServerAddress: serverAddress,
		Insecure:      true,
	}, nil
}

func NewGrpcSecureConfig(serverAddress, clientCert, clientKey string) (*GrpcConfig, error) {
	if serverAddress == "" {
		return nil, errorx.Errorf("NewGrpcSecureConfig: missing server address")
	}
	if clientCert == "" {
		return nil, errorx.Errorf("NewGrpcSecureConfig: missing grpc clientCert")
	}
	if clientKey == "" {
		return nil, errorx.Errorf("NewGrpcSecureConfig: missing grpc clientKey")
	}

	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, errorx.Errorf("NewGrpcSecureConfig: could not create cert from client cert and key pair. %s", err.Error())
	}

	return &GrpcConfig{
		ServerAddress: serverAddress,
		Cert:          cert,
		Insecure:      false,
	}, nil
}

// SetGrpcConnection dials and sets connection of GrpcClient member
// The proper GrpcClient settings need to be set first
func (c *GrpcClient) SetGrpcConnection() error {
	if c.Config.ServerAddress == "" {
		return errorx.Errorf("SetGrpcConnection: missing GrpcConfig")
	}

	var opts = []grpc.DialOption{
		GetKeepaliveDialOption(),
	}

	if c.Config.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		auth := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{c.Config.Cert}})
		opts = append(opts, grpc.WithTransportCredentials(auth))
	}

	opts = append(opts,
		grpc.WithUnaryInterceptor(goprom.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(goprom.StreamClientInterceptor),
	)

	conn, err := grpc.Dial(c.Config.ServerAddress, opts...)
	if err != nil {
		c.Conn = nil
		return err
	}

	c.Conn = conn
	return nil
}

func GetKeepaliveDialOption() grpc.DialOption {
	// the below settings were originally chosen to align with common-scala: https://github.xxx.com/energy/common-scala/tree/main/api-grpc#api-grpc
	// they were lowered as an action item from: TESRE-5046, PWG-3335 to match https://github.xxx.com/energy/influx-tsapi/blob/main/src/main/resources/reference.conf#L192-L200
	return grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             2 * time.Second,
		PermitWithoutStream: true,
	})
}

// GetOutgoingPopulatedGrpcContext returns context with Authorization and X-Request-ID
// headers populated to be used as the context in rpc calls.
func GetOutgoingPopulatedGrpcContext(ctx context.Context, token string) context.Context {
	md := metadata.Pairs(
		"Authorization", token,
		"X-Request-ID", utils.GetRequestIdFromContext(ctx),
	)
	return metadata.NewOutgoingContext(ctx, md)
}
