package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

func LoadTestConfig(optionalPath ...string) (*Config, error) {
	// Reset Viper to avoid interference from previous config loads.
	viper.Reset()

	// Set the environment to test.
	viper.Set("ENV", EnvTest)

	// Determine the test config file path.
	var testEnvPath string
	if len(optionalPath) > 0 && optionalPath[0] != "" {
		testEnvPath = optionalPath[0]
	} else {
		// Use the caller directory to build a default path for test.env.
		_, callerFile, _, ok := runtime.Caller(1)
		if !ok {
			return nil, fmt.Errorf("failed to get caller information")
		}
		// You may adjust this if your test.env file is in a different location.
		testEnvPath = filepath.Join(filepath.Dir(callerFile), "test.env")
	}

	log.Printf("Loading test config from: %s", testEnvPath)

	// Check if the file exists.
	if _, err := os.Stat(testEnvPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %s", testEnvPath)
	} else if err != nil {
		return nil, fmt.Errorf("error checking configuration file: %w", err)
	}

	// Tell Viper to use the test config file.
	viper.SetConfigFile(testEnvPath)
	viper.AutomaticEnv()
	bindAllEnvVars() // bind env variables as before

	// Read in the config file.
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Failed to read config file: %v. Falling back to env variables.", err)
		// Not fatal—we continue with env variables.
	}

	// Unmarshal into our Config struct.
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Set flags based on ENV.
	config.IsDevelopment = config.Env == EnvDevelopment
	config.IsTest = config.Env == EnvTest

	// Cache the config.
	cfg = &config
	return cfg, nil
}

func GetConfig() *Config {
	if cfg == nil {
		log.Fatal("Config not initialized. Call LoadConfig() first.")
	}
	return cfg
}
