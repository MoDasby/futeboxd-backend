package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres Postgres
	Football Football
	Server   Server
	Resend   Resend
	Frontend Frontend
	AWS      AWS
	Cookies  Cookies
	Axiom    Axiom
}

type Axiom struct {
	ApiKey string
}

type Cookies struct {
	Name     string
	HttpOnly bool
	Secure   bool
	SameSite string
	Path     string
	Domain string
}

type Resend struct {
	ApiKey string
}

type Frontend struct {
	URL string
}

type Server struct {
	Port int
}

type Postgres struct {
	User     string
	DBName   string
	Password string
	Host     string
	Port     int
}

type Football struct {
	Url string
}

type AWS struct {
	Secret             string
	Key                string
	S3Endpoint         string
	S3DownloadEndpoint string
	Region             string
}

func getConfigFilename() string {
	if IsProduction() {
		return "./config"
	}

	return "./config/config-dev"
}

func IsProduction() bool {
	return os.Getenv("ENV") == "prod"
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName(getConfigFilename())
	v.AddConfigPath(".")
	v.SetConfigType("yml")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	config, err := parseConfig(v)

	if IsProduction() {
		secretPassword, err := os.ReadFile(os.Getenv("DB_PASSWORD_FILE"))
		if err != nil {
			return nil, err
		}

		config.Postgres.Password = strings.TrimSpace(string(secretPassword))
	}

	return config, err
}

func parseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		slog.Error(fmt.Sprintf("unable to decode into struct, %v", err))
		return nil, err
	}

	return &c, nil
}
