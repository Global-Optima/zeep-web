package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	SecretKey string
}

type ServerConfig struct {
	Port int
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigFile("config.yml")
	viper.AddConfigPath(".")
	// TODO: Add dependent on mode load (e.g. Development, Testing)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
