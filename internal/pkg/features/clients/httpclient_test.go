package clients

import (
	"bytes"
	"context"
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.startlite.cn/itapp/startlite/internal/pkg/infra/util/utils"
)

type testParams struct {
	A string
	B bool
	C int
}

func TestNewInternalClientConfig(t *testing.T) {
	tests := []struct {
		desc          string
		host          string
		cert          string
		key           string
		expectedError string
	}{{
		desc: "success",
		host: "valid-host",
		cert: "valid-cert",
		key:  "valid-key",
	}, {
		desc:          "no host",
		host:          "",
		cert:          "valid-cert",
		key:           "valid-key",
		expectedError: "Host not provided for Internal HTTP Client",
	}, {
		desc:          "no cert",
		host:          "valid-host",
		cert:          "",
		key:           "valid-key",
		expectedError: "Cert not provided for Internal HTTP Client",
	}, {
		desc:          "no key",
		host:          "valid-host",
		cert:          "valid-cert",
		key:           "",
		expectedError: "Key not provided for Internal HTTP Client",
	}}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config, err := NewInternalClientConfig(test.host, test.cert, test.key)
			if test.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected no error got %s", err.Error())
				}
				if test.host != config.Host {
					t.Fatalf("Expected config.Host %s, got %s", test.host, config.Host)
				}
				if test.cert != config.Cert {
					t.Fatalf("Expected config.Cert %s, got %s", test.cert, config.Cert)
				}
				if test.key != config.Key {
					t.Fatalf("Expected config.Key %s, got %s", test.key, config.Key)
				}
			} else {
				if config != nil {
					t.Fatalf("Expected nil config, got %v", config)
				}
				if err == nil {
					t.Fatalf("Expected error %s, got nil", test.expectedError)
				}
				if err.Error() != test.expectedError {
					t.Fatalf("Expected error %s, got %s", test.expectedError, err.Error())
				}
			}
		})
	}
}

func TestNewExternalClientConfig(t *testing.T) {
	tests := []struct {
		desc          string
		host          string
		expectedError string
	}{{
		desc: "success",
		host: "valid-host",
	}, {
		desc:          "no host",
		host:          "",
		expectedError: "Host not provided for External HTTP Client",
	}}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config, err := NewExternalClientConfig(test.host)
			if test.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected no error got %s", err.Error())
				}
				if test.host != config.Host {
					t.Fatalf("Expected config.Host %s, got %s", test.host, config.Host)
				}
			} else {
				if config != nil {
					t.Fatalf("Expected nil config, got %v", config)
				}
				if err == nil {
					t.Fatalf("Expected error %s, got nil", test.expectedError)
				}
				if err.Error() != test.expectedError {
					t.Fatalf("Expected error %s, got %s", test.expectedError, err.Error())
				}
			}

		})
	}
}

func TestSetupHttpClient(t *testing.T) {
	config, _ := NewExternalClientConfig("host")
	timeout := time.Duration(5)
	maxConns := 2
	config.RestTimeout = &timeout
	config.KeepAliveTimeout = &timeout
	config.IdleConnTimeout = &timeout
	config.MaxConnsPerHost = &maxConns

	client, err := setupHTTPClient(*config)
	if err != nil {
		t.Fatalf("Expected nil client got %s", err.Error())
	}
	if client == nil {
		t.Fatal("Expected an http.Client got nil")
	}
}

func TestGenerateTlsConfig(t *testing.T) {
	conf := generateTlsConfig(nil)
	if conf != nil {
		t.Fatalf("Expected nil but got a tls.Config %v", conf)
	}

	cert := &tls.Certificate{}
	conf = generateTlsConfig(cert)
	if conf == nil {
		t.Fatalf("Expected tls.Config ut got nil")
	}
}

func TestNewRequestWithPromContext(t *testing.T) {
	params := map[string]interface{}{"user_id": "email-address"}
	ctx := utils.NewServiceContextWithValue(context.Background(), utils.RequestIDCtxID, "test-request-id")
	req, err := NewRequestWithPromContext(ctx, http.MethodPost, "host/", "", nil, params, bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatalf("Expected nil got error %s", err.Error())
	}
	if req == nil {
		t.Fatal("Expected a request got nil")
	}
	if req.Header.Get("X-Request-ID") != "test-request-id" {
		t.Fatal("Expected X-Request-ID: test-request-id")
	}
}

func TestNewRequestWithProm(t *testing.T) {
	params := map[string]interface{}{"user_id": "email-address"}
	req, err := NewRequestWithProm(http.MethodPost, "host/", "", nil, params, bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatalf("Expected nil got error %s", err.Error())
	}
	if req == nil {
		t.Fatal("Expected a request got nil")
	}
}

func TestQueryFromInterface(t *testing.T) {
	queryParams, err := QueryFromInterface(testParams{
		A: "abc",
		B: true,
		C: 123,
	})
	if err != nil {
		t.Fatal("failed to create query params from struct")
	}

	if queryParams["A"] != "abc" || queryParams["B"] != "true" || queryParams["C"] != "123" {
		t.Fatalf(
			"failed to correctly serialize query struct values: %v, %v, %v",
			queryParams["A"], queryParams["B"], queryParams["C"],
		)
	}
}
