package types

type M2MConfig struct {
	URL          string `yaml:"url" env:"M2M_URL"`
	ClientId     string `yaml:"client_id" env:"M2M_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" env:"M2M_CLIENT_SECRET"`
	GrantType    string `yaml:"grant_type" env:"M2M_GRANT_TYPE"`
}

type KeycloakM2MConfig struct {
	Url      string `yaml:"url" env:"KEYCLOAK_M2M_URL"`
	Username string `yaml:"username" env:"KEYCLOAK_M2M_USERNAME"`
	Secret   string `yaml:"secret" env:"KEYCLOAK_M2M_SECRET"`
	Realm    string `yaml:"realm" env:"KEYCLOAK_M2M_REALM"`
}

type KeycloakAdminConfig struct {
	URL          string `yaml:"url" env:"KEYCLOAK_ADMIN_URL"`
	AuthUrl      string `yaml:"auth_url" env:"KEYCLOAK_ADMIN_AUTH_URL"`
	AuthUsername string `yaml:"auth_username" env:"KEYCLOAK_ADMIN_AUTH_USERNAME"`
	AuthPassword string `yaml:"auth_password" env:"KEYCLOAK_ADMIN_AUTH_PASSWORD"`
	AuthRealm    string `yaml:"auth_realm" env:"KEYCLOAK_ADMIN_AUTH_REALM"`
}

type SmsConfig struct {
	URL       string `yaml:"url" env:"SMS_URL"`
	SpCode    string `yaml:"sp_code" env:"SMS_SP_CODE"`
	LoginName string `yaml:"login_name" env:"SMS_LOGIN_NAME"`
	Password  string `yaml:"password" env:"SMS_PASSWORD"`
}

type S3Config struct {
	Region          string `yaml:"region" env:"S3_REGION"`
	Bucket          string `yaml:"bucket" env:"S3_BUCKET"`
	AccessKeyID     string `yaml:"access_key_id" env:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `yaml:"secret_access_key" env:"S3_SECRET_ACCESS_KEY"`
	Endpoint        string `yaml:"endpoint" env:"S3_ENDPOINT"`
	Style           string `yaml:"style" env:"S3_STYLE"`
}

type SftpConfig struct {
	Host       string `yaml:"host" env:"SFTP_HOST"`
	Port       string `yaml:"port" env:"SFTP_PORT"`
	User       string `yaml:"user" env:"SFTP_USER"`
	Password   string `yaml:"password" env:"SFTP_PASSWORD"`
	PrivateKey string `yaml:"private_key" env:"SFTP_PRIVATE_KEY"`
}

type ConfluenceConfig struct {
	Host     string `yaml:"host" env:"CONFLUENCE_HOST"`
	Username string `yaml:"username" env:"CONFLUENCE_USERNAME"`
	Password string `yaml:"password" env:"CONFLUENCE_PASSWORD"`
}

type ADFSConfig struct {
	Endpoint string `yaml:"endpoint" env:"ADFS_ENDPOINT"`
	Key      string `yaml:"key" env:"ADFS_KEY"`
}

type WorkflowModelerConfig struct {
	Endpoint string `yaml:"endpoint" env:"WORKFLOWMODELER_ENDPOINT"`
}

type PacmanConfig struct {
	PacmanToken string `yaml:"pacman_token" env:"PACMAN_TOKEN"`
}

type KafkaPubSubConfig struct {
	Brokers            string `yaml:"brokers" env:"PACMAN_BROKERS"`
	DLQTopic           string `yaml:"dlq_topic" env:"PACMAN_DLQ_TOPIC"`
	GroupID            string `yaml:"group_id" env:"PACMAN_GROUP_ID"`
	GroupPrefix        string `yaml:"group_prefix" env:"PACMAN_GROUP_PREFIX"`
	SessionTimeout     int    `yaml:"session_timeout" env:"PACMAN_SESSION_TIMEOUT"`
	TestPublishEnabled bool   `yaml:"test_publish_enabled" env:"PACMAN_TEST_PUBLISH_ENABLED"`

	// TLS
	TLSPathToCAFile   string `yaml:"kafka_ca_file" env:"PACMAN_CA_FILE"`
	TLSPathToCertFile string `yaml:"kafka_cert_file" env:"PACMAN_CERT_FILE"`
	TLSPathToKeyFile  string `yaml:"kafka_key_file" env:"PACMAN_KEY_FILE"`
}

type KafkaBrokerConfig struct {
	Addrs      []string `yaml:"addrs" env:"KAFKA_ADDR"`
	TLSEnable  bool     `yaml:"tls_enable" env:"KAFKA_TLS_ENABLE"`
	ClientCert string   `yaml:"client_cert" env:"KAFKA_CLIENT_CERT"`
	ClientKey  string   `yaml:"client_key" env:"KAFKA_CLIENT_KEY"`
	CaCert     string   `yaml:"ca_cert" env:"KAFKA_CA_CERT"`
}

type KafkaProducerConfig struct {
	KafkaBrokerConfig
	Topic string `yaml:"topic" env:"KAFKA_PRODUCER_TOPIC"`
}

type KafkaConsumerConfig struct {
	KafkaBrokerConfig
	Topics  []string `yaml:"topics" env:"KAFKA_CONSUMER_TOPICS"`
	GroupID string   `yaml:"group_id" env:"KAFKA_CONSUMER_GROUP_ID"`
}
type AmsConfig struct {
	Endpoint string `yaml:"endpoint" env:"AMS_ENDPOINT"`
	AppId    string `yaml:"app_Id" env:"AMS_APP_ID"`
	FuncId   string `yaml:"func_Id" env:"AMS_FUNC_ID"`
}

type SsoConfig struct {
	URL          string `yaml:"url" env:"SSO_URL"`
	ClientId     string `yaml:"client_id" env:"SSO_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" env:"SSO_CLIENT_SECRET"`
	GrantType    string `yaml:"grant_type" env:"SSO_GRANT_TYPE"`
	Scope        string `yaml:"scope" env:"SSO_SCOPE"`
	Resource     string `yaml:"resource" env:"SSO_RESOURCE"`
}

type SmtpConfig struct {
	Endpoint    string `yaml:"endpoint" env:"SMTP_ENDPOINT"`
	Port        string `yaml:"port" env:"SMTP_PORT"`
	User        string `yaml:"user" env:"SMTP_USER"`
	Password    string `yaml:"password" env:"SMTP_PASSWORD"`
	AddressFrom string `yaml:"address_from" env:"SMTP_ADDRESS_FROM"`
	AddressName string `yaml:"address_name" env:"SMTP_ADDRESS_Name"`
}

type WarpNotification struct {
	URL          string `yaml:"url" env:"WARP_NOTIFICATION_URL"`
	MessagePath  string `yaml:"message_path" ENV:"WARP_NOTIFICATION_MESSAGE_PATH"`
	TemplatePath string `yaml:"template_path" ENV:"WARP_NOTIFICATION_TEMPLATE_PATH"`
	Resource     string `yaml:"resource" env:"WARP_NOTIFICATION_RESOURCE"`
}
