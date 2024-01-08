package typesx

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

	WorkflowConfig struct {
		URL string `yaml:"url" env:"WORKFLOW_URL"`
	}

	RedisConfig struct {
		Host     string `yaml:"host" env:"REDIS_HOST"`
		Port     int    `yaml:"port" env:"REDIS_PORT"`
		PoolSize int    `yaml:"poolsize" env:"REDIS_POOL_SIZE"`
		Secret   string `yaml:"secret" env:"REDIS_SECRET"`
		DB       int    `yaml:"db" env:"REDIS_DB"`
	}

	OAuth2Config struct {
		SSOIssuer    string `yaml:"sso_issuer" env:"OAUTH2_ISSUER"`
		ClientID     string `yaml:"client_id" env:"OAUTH2_CLIENT_ID"`
		ClientSecret string `yaml:"client_secret" env:"OAUTH2_CLIENT_SECRET"`
		PublicKey    string `yaml:"public_key" env:"OAUTH2_PUBLIC_KEY"`
	}

	NewRedisConfig struct {
		Addr             string `yaml:"addr" env:"REDIS_ADDR"` // cluster contains at least one ","
		Password         string `yaml:"password" env:"REDIS_PASSWORD"`
		DB               int    `yaml:"db" env:"REDIS_DB"`
		PoolSize         int    `yaml:"pool_size" env:"REDIS_POOL_SIZE"`
		MasterName       string `yaml:"master_name" env:"REDIS_MASTER_NAME"`
		SentinelPassword string `yaml:"sentinel_password" env:"REDIS_SENTINEL_PASSWORD"`
		KeyPrefix        string `yaml:"key_prefix" env:"REDIS_KEY_PREFIX"`
	}

	LdapConfig struct {
		URL      string `yaml:"url" env:"LDAP_URL"`
		Username string `yaml:"username" env:"LDAP_USERNAME"`
		Password string `yaml:"password" env:"LDAP_PASSWORD"`
	}

	TwpConsumerConfig struct {
		Name           string `yaml:"name" env:"TWP_CONSUMER_NAME"`
		Topics         string `yaml:"topics" env:"TWP_CONSUMER_TOPICS"`
		ConsumerGroup  string `yaml:"consumer_group" env:"TWP_CONSUMER_GROUP"`
		Partition      int    `yaml:"partition" env:"TWP_CONSUMER_PARTITION"`
		BoostrapServer string `yaml:"bootstrap_server" env:"TWP_CONSUMER_BOOTSTRAP_SERVER"`
		TLSEnable      bool   `yaml:"tls_enable"  env:"TWP_CONSUMER_TLS_ENABLE"`
		TLSCA          string `yaml:"tls_ca" env:"TWP_CONSUMER_CA"`
		TLSCert        string `yaml:"tls_cert" env:"TWP_CONSUMER_CERT"`
		TLSKey         string `yaml:"tls_key" env:"TWP_CONSUMER_KEY"`
		Secret         string `yaml:"secret" env:"TWP_CONSUMER_SECRET"`
		Apps           string `yaml:"apps" env:"TWP_CONSUMER_APPS"`
	}
)
