package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type JWTConfig struct {
	Secret     string
	ExpireTime int
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found, using default values: %v", err)
		setDefaults()
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
		},
		Database: DatabaseConfig{
			Driver: viper.GetString("database.driver"),
			DSN:    viper.GetString("database.dsn"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("jwt.secret"),
			ExpireTime: viper.GetInt("jwt.expire_time"),
		},
	}
}

func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.dsn", "association.db")
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expire_time", 24)
}
