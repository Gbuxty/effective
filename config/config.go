package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	pathConfigFile = ".env"
	dotenv         = "dotenv"
)

type Config struct {
	HTTP     *HTTPServer
	Postgres *PostgresConfig
	APIUrl   *APIUrl
}

type HTTPServer struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type APIUrl struct {
	AgifyUrl       string
	GenderizeUrl   string
	NationalizeUrl string
}

func Load() (*Config, error) {
	viper.SetConfigFile(pathConfigFile)
	viper.SetConfigType(dotenv)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		HTTP: &HTTPServer{
			Port:            viper.GetInt("HTTP_SERVER_PORT"),
			ReadTimeout:     viper.GetDuration("HTTP_READ_TIMEOUT"),
			WriteTimeout:    viper.GetDuration("HTTP_WRITE_TIMEOUT"),
			ShutdownTimeout: viper.GetDuration("HTTP_SHUTDOWN_TIMEOUT"),
		},
		Postgres: &PostgresConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetInt("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("POSTGRES_DB"),
			SSLMode:  viper.GetString("POSTGRES_SSL_MODE"),
		},
		APIUrl: &APIUrl{
			AgifyUrl:       viper.GetString("AGIFY_URL"),
			GenderizeUrl:   viper.GetString("GENDERIZE_URL"),
			NationalizeUrl: viper.GetString("NATIONALIZE_URL"),
		},
	}
	return cfg, nil
}

func (p PostgresConfig) ToDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DBName, p.SSLMode)
}
