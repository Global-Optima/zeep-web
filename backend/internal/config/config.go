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
	Env           string `mapstructure:"ENV"`
	IsDevelopment bool
	IsTest        bool

	Database DatabaseConfig `mapstructure:",squash"`
	Server   ServerConfig   `mapstructure:",squash"`
	JWT      JWTConfig      `mapstructure:",squash"`
	S3       S3Config       `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
	Kafka    KafkaConfig    `mapstructure:",squash"`
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
