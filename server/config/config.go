package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Server struct {
	ListenAddress       string `koanf:"listen_address"`
	JWTKey              string `koanf:"jwt_key"`
	JWTPublicKey        string `koanf:"jwt_public_key"`
	AllowedOrigins      string `koanf:"allowed_origins"`
	CookieDomain        string `koanf:"cookie_domain"`
	JWTKeyPath          string `koanf:"jwt_private_key_path"`
	JWTPublicKeyPath    string `koanf:"jwt_public_key_path"`
	GoogleEmailPassword string `koanf:"google_email_password"`
	GoogleService       string `koanf:"google_service"`
	GmailHost           string `koanf:"gmail_host"`
	GmailPort           string `koanf:"gmail_port"`
	EmailPort           string `koanf:"email_port"`
	FromEmail           string `koanf:"from_email"`
	DatabaseUrl         string `koanf:"database_url"`
	DBHost              string `koanf:"db_host"`
	DBUser              string `koanf:"db_user"`
	DBName              string `koanf:"db_name"`
	DBPassword          string `koanf:"db_password"`
	DBPort              string `koanf:"db_port"`
}

type Log struct {
	Level slog.Level `koanf:"level"`
}

type Cloudinary struct {
	CloudName          string `koanf:"cloud_name"`
	ApiKey             string `koanf:"api_key"`
	ApiSecret          string `koanf:"api_secret"`
	BaseDeliveryUrl    string `koanf:"base_delivery_url"`
	SecretDeliveryUrl  string `koanf:"secret_delivery_url"`
	ApiBaseUrl         string `koanf:"api_base_url"`
	ApiProvisioningUrl string `koanf:"api_provisioning_url"`
	Url                string `koanf:"url"`
}

type Config struct {
	Server     Server     `koanf:"server"`
	Log        Log        `koanf:"log"`
	Cloudinary Cloudinary `koanf:"cloudinary"`
}

var k = koanf.New(".")

/*
Loads configurations either or both from yaml file or .env file
.env file overrides configurations in yaml file
*/
func Load(configPath *string) (*Config, error) {
	var errYAML error
	var errENV error
	configurations := &Config{}

	if configPath != nil {
		path := file.Provider(*configPath)
		// load configuration from  yaml file
		errYAML = k.Load(path, yaml.Parser())
	}

	// or load/override configurations with env file
	errENV = utils.CoalesceError(godotenv.Load("../.env"), godotenv.Load(".env"))

	// continue if reading either works and throw error otherwise
	if errYAML != nil && errENV != nil {
		fmt.Println("here", errENV)
		return configurations, utils.CoalesceError(errENV, errYAML)
	}

	// env variables should be prexied with CS_ and __ to indicate nest
	err := k.Load(env.Provider("CONFIG_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CONFIG_")), "__", ".", -1)
	}), nil)

	if err != nil {
		return configurations, err
	}

	err = k.UnmarshalWithConf("", configurations, koanf.UnmarshalConf{Tag: "koanf"})

	if err != nil {
		return configurations, err
	}

	//Load JWT keys from file path
	if err = loadJWTKeys(configurations); err != nil {
		return configurations, err
	}

	return configurations, nil
}

func loadJWTKeys(config *Config) error {
	if config.Server.JWTKeyPath == "" || config.Server.JWTPublicKeyPath == "" {
		return fmt.Errorf("jwt keys path is required")
	}

	privateKey, err := os.ReadFile(config.Server.JWTKeyPath)
	if err != nil {
		return fmt.Errorf("Failed to read private key: %v", err)
	}

	publicKey, err := os.ReadFile(config.Server.JWTPublicKeyPath)
	if err != nil {
		return fmt.Errorf("Failed to read public key: %v", err)
	}

	config.Server.JWTKey = string(privateKey)
	config.Server.JWTPublicKey = string(publicKey)
	return nil
}
