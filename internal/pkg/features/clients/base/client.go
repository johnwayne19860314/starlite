package base

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	goprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.startlite.cn/itapp/startlite/internal/pkg/features/clients"
	"github.startlite.cn/itapp/startlite/internal/pkg/types"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultRestTimeout           = 5 * time.Second
	defaultMaxConnectionsPerHost = 100
	defaultIdleConnTimeout       = 30 * time.Second
	defaultKeepAliveTimeout      = 60 * time.Second
	defaultRequestQueueTimeout   = 7 * time.Second
	defaultRequestLimit          = 1
)

var instance *ClientProp

type ClientProp struct {
	certType        string
	commandClientCa string
}

type BaseClient interface {
	get() (string, string)
}

func (c *ClientProp) get() (string, string) {
	return c.certType, c.commandClientCa
}
func MustNewBaseClient(appCtx appx.AppContext, cl *featurex.ConfigLoader) BaseClient {

	clientConfig := &types.ClientConfig{}
	cl.Load(clientConfig)

	instance = &ClientProp{
		certType:        clientConfig.CertType,
		commandClientCa: clientConfig.ClientCA,
	}

	return instance

}

func BuildConnection(host, port, cert, key, svc string, insecure bool) (*grpc.ClientConn, error) {

	if host == "" || port == "" {
		logx.Error("Either host or port are not specified ", "host", host, "port", port)
		return nil, errorx.Errorf("Either host {%s} or port {%s} are not specified ", host, port)
	}

	var address = fmt.Sprintf("%v:%v", host, port)

	var opts = []grpc.DialOption{
		clients.GetKeepaliveDialOption(),
		grpc.WithUnaryInterceptor(goprom.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(goprom.StreamClientInterceptor),
	}

	if !insecure {
		logx.Info("========== use tls certification========", "service", svc)
		if cert == "" || key == "" {
			logx.Error("Either cert or key are not specified ", cert, key)
			return nil, errorx.Errorf("Either cert {%s} or key {%s} are not specified ", cert, key)
		}
		var ca tls.Certificate
		ca, err := GetCert(cert, key)
		if err != nil {
			logx.Error("Unable to get cert from files")
			return nil, err
		}

		var tlsConfig = generateTlsConfig(ca)
		// if svc == constant.COMMAND_SERVICE {
		// 	// Load the CA certificate file
		// 	var caCert []byte
		// 	var err error
		// 	if instance.certType == constant.FILE_TYPE {
		// 		caCert, err = os.ReadFile(instance.commandClientCa)
		// 		if err != nil {
		// 			logx.Error("Unable to get Client CA")
		// 			return nil, err
		// 		}
		// 	} else {
		// 		caCert = []byte(instance.commandClientCa)
		// 	}

		// 	// Create a certificate pool and add the CA certificate to it
		// 	caCertPool := x509.NewCertPool()
		// 	caCertPool.AppendCertsFromPEM(caCert)

		// 	tlsConfig.ClientCAs = caCertPool
		// 	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		// 	tlsConfig.RootCAs = caCertPool
		// }

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		logx.Info("========== not use tls certification========", "service", svc)
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logx.Error("Failed to connect: command service " + address)
		return nil, err
	}
	return conn, nil
}

func GetCert(certInput, key string) (cert tls.Certificate, err error) {
	if certInput == "" || key == "" {
		return
	}

	// if instance.certType == constant.FILE_TYPE {
	// 	cert, err = tls.LoadX509KeyPair(certInput, key)
	// } else {
	// 	cert, err = tls.X509KeyPair([]byte(certInput), []byte(key))
	// }

	if err != nil {
		err = errorx.Errorf("can't loadX509KeyPair, %v", err)
		return
	}
	return
}

func generateTlsConfig(cert tls.Certificate) (tlsConfig *tls.Config) {
	tlsConfig = &tls.Config{}
	if len(cert.Certificate) == 0 {
		return
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	return
}

func GetToken(server, certPath, keyPath string) string {
	url := server + "tokens"
	method := "POST"

	// Load client cert
	cert, err := GetCert(certPath, keyPath)
	if err != nil {
		logx.Error(err.Error())
		return ""
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}
	// Setup request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		logx.Error(err.Error())
		return ""
	}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		logx.Error(err.Error())
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logx.Error(err.Error())
		return ""
	}
	t := new(tokenResponse)
	err = json.Unmarshal(body, t)
	if err != nil {
		logx.Error(err.Error())
		return err.Error()
	}
	return t.Data.Token
}

type tokenResponse struct {
	Meta struct {
		RequestId    string `json:"request_id"`
		ClaimsLimits struct {
		} `json:"claims_limits"`
		Pagination interface{} `json:"pagination"`
	} `json:"meta"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Links interface{} `json:"links"`
}

func GenerateReqID() string {
	bytes := make([]byte, 16) // 16 bytes gives us 32 hex characters
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

type TokenServer struct {
	Server string
	Cert   string
	Key    string
}

func SetupHttpClient(ClientCert, ClientKey string) (newClient *metrics.HTTPClient, err error) {
	var tlsConfig = &tls.Config{}
	ServerCaCert := ""
	if ServerCaCert != "" {
		clientCACert, err := os.ReadFile(ServerCaCert)
		if err != nil {
			err = errorx.Errorf("Can't read cert file %v error is %v", ServerCaCert, err)
			return nil, err
		}
		clientCertPool := x509.NewCertPool()
		clientCertPool.AppendCertsFromPEM(clientCACert)
		tlsConfig = &tls.Config{
			RootCAs: clientCertPool,
		}
	}

	if ClientCert != "" && ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(ClientCert, ClientKey)
		if err != nil {
			return nil, errorx.New("can't LoadX509KeyPair " + err.Error())
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		//tlsConfig.BuildNameToCertificate()
	}
	MaxConnectionsPerHost := defaultMaxConnectionsPerHost
	RestTimeout := defaultRestTimeout
	KeepAliveTimeout := defaultKeepAliveTimeout
	IdleConnTimeout := defaultIdleConnTimeout
	// logx.Info("Asset Max Connections Per Host: %v\n", MaxConnectionsPerHost)
	// logx.Info("Asset Timeout: %v\n", RestTimeout)
	client := &http.Client{
		Transport: clients.InstrumentPromRoundTripper(&http.Transport{
			MaxConnsPerHost:     MaxConnectionsPerHost,
			MaxIdleConnsPerHost: MaxConnectionsPerHost,
			Dial: (&net.Dialer{
				KeepAlive: KeepAliveTimeout,
				Timeout:   RestTimeout,
			}).Dial,
			TLSClientConfig:     tlsConfig,
			TLSHandshakeTimeout: RestTimeout,
			IdleConnTimeout:     IdleConnTimeout,
		}),
		Timeout: RestTimeout,
	}

	newClient = &metrics.HTTPClient{Client: client}

	return
}
