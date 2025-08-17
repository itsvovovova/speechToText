package config

import (
	"os"
)

const ()

var CurrentConfig = NewConfig()

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	Logger   *LoggerConfig
	RabbitMQ *RabbitMQConfig
	Redis    *RedisConfig
}

type ServerConfig struct {
	Port     string
	Host     string
	HostPort string
}

type DatabaseConfig struct {
	Username     string
	Name         string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	SSLMode      string
}

type LoggerConfig struct {
	Level  string
	Format string
}

type RedisConfig struct {
	Host string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	UserPort string
}

func NewConfig() *Config {
	var databaseConfig = DatabaseConfig{
		Username:     os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		Name:         os.Getenv("DB_NAME"),
		DatabaseName: os.Getenv("DB_NAME"),
		SSLMode:      os.Getenv("DB_SSL_MODE"),
	}

	var loggerConfig = LoggerConfig{
		Level:  os.Getenv("LOG_LEVEL"),
		Format: os.Getenv("LOG_FORMAT"),
	}

	var serverConfig = ServerConfig{
		Port:     os.Getenv("SERVER_PORT"),
		Host:     os.Getenv("SERVER_HOST"),
		HostPort: os.Getenv("HOST_PORT"),
	}

	var redisConfig = RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
	}

	var rabbitMQConfig = RabbitMQConfig{
		Host:     os.Getenv("RABBIT_HOST"),
		Port:     os.Getenv("RABBIT_PORT"),
		Username: os.Getenv("RABBIT_USER"),
		Password: os.Getenv("RABBIT_PASSWORD"),
		UserPort: os.Getenv("RABBIT_USER_PORT"),
	}

	var Config = &Config{
		Server:   &serverConfig,
		Database: &databaseConfig,
		Logger:   &loggerConfig,
		RabbitMQ: &rabbitMQConfig,
		Redis:    &redisConfig,
	}
	return Config
}
