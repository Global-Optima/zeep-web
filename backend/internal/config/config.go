package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvTest        = "test"
)

type Config struct {
	Env          string `mapstructure:"ENV"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       int    `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
	ServerPort   int    `mapstructure:"SERVER_PORT"`
	ClientUrl    string `mapstructure:"CLIENT_URL"`
	S3AccessKey  string `mapstructure:"PSKZ_ACCESS_KEY"`
	S3SecretKey  string `mapstructure:"PSKZ_SECRET_KEY"`
	S3Endpoint   string `mapstructure:"PSKZ_ENDPOINT"`
	S3BucketName string `mapstructure:"PSKZ_BUCKETNAME"`

	IsDevelopment bool
	IsTest        bool
}

var cfg *Config

func LoadConfig(file string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.IsDevelopment = cfg.Env == EnvDevelopment
	cfg.IsTest = cfg.Env == EnvTest

	return cfg, nil
}

func GetConfig() *Config {
	if cfg == nil {
		log.Fatal("Config not loaded")
	}
	return cfg
}
