package clients

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.startlite.cn/itapp/startlite/internal/pkg/infra/util/utils"
)

func TestNewGrpcInsecureConfig(t *testing.T) {
	tests := []struct {
		desc          string
		serverAddress string
		expectedError string
	}{{
		desc:          "success",
		serverAddress: "test-server",
	}, {
		desc:          "No server address",
		serverAddress: "",
		expectedError: "NewGrpcInsecureConfig: missing server address",
	}}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			res, err := NewGrpcInsecureConfig(test.serverAddress)
			if test.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected nil got %s", err.Error())
				}
				if res.ServerAddress != test.serverAddress {
					t.Fatalf("Expected GrpcConfig with server address %s got %s", test.serverAddress, res.ServerAddress)
				}
				if res.Insecure != true {
					t.Fatalf("Expected GrpcConfig with Insecure: true, got false")
				}
			}
			if test.expectedError != "" && err.Error() != test.expectedError {
				t.Fatalf("Expected error %s got %s", test.expectedError, err.Error())
			}
		})
	}
}

func TestNewGrpcSecureConfig(t *testing.T) {
	tests := []struct {
		desc          string
		serverAddress string
		cert          string
		key           string
		expectedError string
	}{{
		desc:          "No server address",
		expectedError: "NewGrpcSecureConfig: missing server address",
	}, {
		desc:          "No client cert",
		serverAddress: "test-server",
		expectedError: "NewGrpcSecureConfig: missing grpc clientCert",
	}, {
		desc:          "No client key",
		serverAddress: "test-server",
		cert:          "test-cert",
		expectedError: "NewGrpcSecureConfig: missing grpc clientKey",
	}, {
		desc:          "bad cert-key pair",
		serverAddress: "test-server",
		cert:          "test-cert",
		key:           "test-key",
		expectedError: "NewGrpcSecureConfig: could not create cert from client cert and key pair. open test-cert: no such file or directory",
	}}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			res, err := NewGrpcSecureConfig(test.serverAddress, test.cert, test.key)
			if test.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected nil got %s", err.Error())
				}
				if res.ServerAddress != test.serverAddress {
					t.Fatalf("Expected GrpcConfig with server address %s got %s", test.serverAddress, res.ServerAddress)
				}
				if res.Insecure != false {
					t.Fatalf("Expected GrpcConfig with Insecure: false, got true")
				}
			}
			if test.expectedError != "" && err.Error() != test.expectedError {
				t.Fatalf("Expected error %s got %s", test.expectedError, err.Error())
			}
		})
	}
}

func TestGetOutgoingPopulatedGrpcContext(t *testing.T) {
	ctxWithRequestID := utils.NewServiceContextWithValue(context.Background(), utils.RequestIDCtxID, "test-request-id")
	tests := []struct {
		name        string
		incomingCtx context.Context
		token       string
		expectedCtx context.Context
	}{
		{
			name:        "empties",
			incomingCtx: context.Background(),
			token:       "",
			expectedCtx: metadata.NewOutgoingContext(context.Background(), metadata.Pairs("Authorization", "", "X-Request-ID", "")),
		},
		{
			name:        "populated",
			incomingCtx: ctxWithRequestID,
			token:       "test-token",
			expectedCtx: metadata.NewOutgoingContext(ctxWithRequestID, metadata.Pairs("Authorization", "test-token", "X-Request-ID", "test-request-id")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := GetOutgoingPopulatedGrpcContext(test.incomingCtx, test.token)
			assert.Equal(t, test.expectedCtx, ctx)
		})
	}
}
