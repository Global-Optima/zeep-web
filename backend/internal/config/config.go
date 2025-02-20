package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

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

	Database  DatabaseConfig  `mapstructure:",squash"`
	Server    ServerConfig    `mapstructure:",squash"`
	JWT       JWTConfig       `mapstructure:",squash"`
	S3        S3Config        `mapstructure:",squash"`
	Redis     RedisConfig     `mapstructure:",squash"`
	Kafka     KafkaConfig     `mapstructure:",squash"`
	Filtering FilteringConfig `mapstructure:",squash"`
}

var cfg *Config

func LoadConfig(file string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.AutomaticEnv()

	viper.SetDefault("JWT_CUSTOMER_ACCESS_TOKEN_TTL", DEFAULT_CUSTOMER_ACCESS_TOKEN_TTL)
	viper.SetDefault("JWT_CUSTOMER_REFRESH_TOKEN_TTL", DEFAULT_CUSTOMER_REFRESH_TOKEN_TTL)
	viper.SetDefault("JWT_EMPLOYEE_ACCESS_TOKEN_TTL", DEFAULT_EMPLOYEE_ACCESS_TOKEN_TTL)
	viper.SetDefault("JWT_EMPLOYEE_REFRESH_TOKEN_TTL", DEFAULT_EMPLOYEE_REFRESH_TOKEN_TTL)

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

func LoadTestConfig(optionalPath ...string) (*Config, error) {
	viper.Reset()
	viper.Set("ENV", EnvTest)

	var testEnvPath string
	if len(optionalPath) > 0 && optionalPath[0] != "" {
		testEnvPath = optionalPath[0]
	} else {
		_, callerFile, _, ok := runtime.Caller(1)
		if !ok {
			return nil, fmt.Errorf("failed to get caller information")
		}
		dir := filepath.Dir(callerFile)
		testEnvPath = filepath.Join(dir, "test.env")
	}

	log.Printf("Loading test config from: %s", testEnvPath)

	if _, err := os.Stat(testEnvPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %s", testEnvPath)
	} else if err != nil {
		return nil, fmt.Errorf("error checking configuration file: %w", err)
	}

	return LoadConfig(testEnvPath)
}

func GetConfig() *Config {
	if cfg == nil {
		log.Fatal("Config not loaded")
	}
	return cfg
}
