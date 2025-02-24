package config

import (
	"log"
	"os"
	"strings"
	"sync"

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

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		viper.SetConfigFile(".env") // Load from .env file if available
		viper.AutomaticEnv()        // Read from system environment variables
		bindAllEnvVars()            // Automatically bind all env variables

		// Try loading .env file but do NOT fail if it doesn't exist
		if _, err := os.Stat(".env"); err == nil {
			err := viper.ReadInConfig()
			if err != nil {
				log.Printf("Warning: Failed to read config file: %v. Falling back to env variables.", err)
			}
		}

		// Unmarshal environment variables into Config struct
		var config Config
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		// Determine environment
		config.IsDevelopment = config.Env == EnvDevelopment
		config.IsTest = config.Env == EnvTest

		cfg = &config
	})

	return cfg
}

// Automatically bind all environment variables dynamically
func bindAllEnvVars() {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		_ = viper.BindEnv(key)
	}
}

// GetConfig returns the loaded config instance
func GetConfig() *Config {
	if cfg == nil {
		log.Fatal("Config not initialized. Call LoadConfig() first.")
	}
	return cfg
}
