package types

const (
	//EnvDev defines dev environment
	EnvDev = "local"
	// EnvEng defines eng environment
	EnvEng = "eng"
	// EnvProd defines prod environment
	EnvProd = "prod"
)

type Config struct {
	CLAPI  OutboundApi `mapstructure:"CHARGINGLOCATIONSAPI"`
	CSM    OutboundApi `mapstructure:"CHARGINGSTATEMACHINE"`
	CSS    OutboundApi `mapstructure:"CHARGINGSESSIONSSERVICE"`
	Spark  OutboundApi `mapstructure:"SPARK"`
	JWTKey string      `mapstructure:"JWTPRIVATEKEY"`
}

type OutboundApi struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Cert        string `mapstructure:"cert"`
	Key         string `mapstructure:"key"`
	ClientID    string `mapstructure:"client_id"`
	ClientToken string `mapstructure:"client_token"`
}

type (
	ServerConfig struct {
		Name     string `yaml:"name" env:"SERVER_NAME"`
		Port     int    `yaml:"port" env:"SERVER_PORT"`
		RunLevel string `yaml:"run_level" env:"SERVER_RUN_LEVEL"`
		Host     string `yaml:"host" env:"SERVER_HOST"`
	}

	DatabaseConfig struct {
		Driver string `yaml:"driver" env:"DATABASE_DRIVER"`
		DSN    string `yaml:"dsn" env:"DATABASE_DSN"`
	}

	AssetApi struct {
		AssetServer string `yaml:"assetserver" env:"ASSET_HOST"`
		Cert        string `yaml:"cert" env:"IEM_CERT"`
		Key         string `yaml:"key" env:"IEM_KEY"`
	}
	BiApi struct {
		Host string `yaml:"biserver" env:"BI_HOST"`
		Cert string `yaml:"cert" env:"IEM_CERT"`
		Key  string `yaml:"key" env:"IEM_KEY"`
	}
	ChargingLocationsApi struct {
		Host string `yaml:"host" env:"CLAPI_HOST"`
		Port string `yaml:"port" env:"CLAPI_PORT"`
		Cert string `yaml:"cert" env:"IEM_CERT"`
		Key  string `yaml:"key" env:"IEM_KEY"`
	}

	TokenApi struct {
		TokenServer string `yaml:"host" env:"TOKEN_HOST"`
		// Cert string `yaml:"cert" env:"ASSET_CERT"`
		// Key  string `yaml:"key" env:"ASSET_KEY"`
	}
	CommandService struct {
		Server   string `yaml:"server" env:"COMMAND_HOST"`
		Port     string `yaml:"port" env:"COMMAND_PORT"`
		Cert     string `yaml:"cert" env:"IEM_CERT"`
		Key      string `yaml:"key" env:"IEM_KEY"`
		Insecure bool   `yaml:"insecure" env:"COMMAND_INSECURE"`
	}
	ClientConfig struct {
		CertType string `yaml:"type" env:"CERT_TYPE"`
		ClientCA string `yaml:"client_ca" env:"COMMAND_CLIENT_CA"`
	}
	AzureADConfig struct {
		ClientID     string `yaml:"client_id" env:"AZURE_CLIENT_ID"`
		TenantID     string `yaml:"tenant_id" env:"AZURE_TENANT_ID"`
		ClientSecret string `yaml:"client_secret" env:"AZURE_CLIENT_SECRET"`
		RedirectURL  string `yaml:"redirect_url" env:"AZURE_REDIRECT_URL"`
	}

	ChargingStateMachine struct {
		Host string `yaml:"host" env:"CSM_HOST"`
		Port string `yaml:"port" env:"CSM_PORT"`
		Cert string `yaml:"cert" env:"IEM_CERT"`
		Key  string `yaml:"key" env:"IEM_KEY"`
	}

	ChargingSessionsService struct {
		Host   string `yaml:"host" env:"CSS_HOST"`
		Port   string `yaml:"port" env:"CSS_PORT"`
		Cert   string `yaml:"cert" env:"IEM_CERT"`
		Key    string `yaml:"key" env:"IEM_KEY"`
		JWTKey string `yaml:"jwt_key" env:"CSS_JWT_KEY"`
	}

	Spark struct {
		Host        string `yaml:"host" env:"SPARK_HOST"`
		ClientID    string `yaml:"client_id" env:"SPARK_CLIENT_ID"`
		ClientToken string `yaml:"client_token" env:"SPARK_CLIENT_TOKEN"`
	}

	JWTConfig struct {
		HeaderAuthorization string `yaml:"header_authorization" env:"TCH_JWT_AUTHORIZATION"`
		HeaderBearer        string `yaml:"header_bearer" env:"TCH_JWT_BEARER"`
		TokenExpirationSec  int    `yaml:"token_expiration_sec" env:"TCH_JWT_EXPIRATION_SEC"`
		TokenSecret         string `yaml:"token_secret" env:"TCH_JWT_TOKEN_SECRET"`
		Leeway              int    `yaml:"leeway" env:"TCH_JWT_LEEWAY"`
	}

	EncrypxConfig struct {
		TimestampLeeway int    `yaml:"timestamp_leeway" env:"TCH_ENCRYPX_TIMESTAMP_LEEWAY"`
		TimeLayout      string `yaml:"time_layout" env:"TCH_ENCRYPX_TIME_LAYOUT"`
	}

	TCHConfig struct {
		xxxOperatorId     string `yaml:"xxx_operator_id" env:"TCH_OPERATOR_xxx_OPERATOR_ID"`
		xxxOperatorSecret string `yaml:"xxx_operator_secret" env:"TCH_OPERATOR_xxx_OPERATOR_SECRET"`
		xxxOutSeq         string `yaml:"xxx_out_seq" env:"TCH_OPERATOR_xxx_OUT_SEQ"`
	}

	GOVConfig struct {
		Endpoint      string   `yaml:"gov_endpoint" env:"GOV_ENDPOINT"`
		TCECWhiteList []string `yaml:"gov_tcec_white_list" env:"GOV_TCEC_WHITE_LIST"`
	}

	KafkaClientConfig struct {
		Name                 string `yaml:"name" env:"KAFKA_CLIENT_TCH_NAME"`
		ConsumerTopics       string `yaml:"consumer_topics" env:"KAFKA_CLIENT_TCH_CONSUMER_TOPICS"`
		ConsumerGroup        string `yaml:"consumer_group" env:"KAFKA_CLIENT_TCH_GROUP"`
		ProducerTopics       string `yaml:"producer_topics" env:"KAFKA_CLIENT_TCH_PRODUCER_TOPICS"`
		PairedConsumerTopics string `yaml:"paired_consumer_topics" env:"KAFKA_CLIENT_TCH_PAIRED_CONSUMER_TOPICS"`
		PairedConsumerGroup  string `yaml:"paired_consumer_group" env:"KAFKA_CLIENT_TCH_PAIRED_CONSUMER_GROUP"`
		BoostrapServer       string `yaml:"bootstrap_server" env:"KAFKA_CLIENT_TCH_BOOTSTRAP_SERVER"`
		TLSCA                string `yaml:"tls_ca" env:"KAFKA_CLIENT_TCH_CA"`
		TLSCert              string `yaml:"tls_cert" env:"KAFKA_CLIENT_TCH_CERT"`
		TLSKey               string `yaml:"tls_key" env:"KAFKA_CLIENT_TCH_KEY"`
	}
)
