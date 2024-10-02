package config

import (
	"os"
	"strings"

	"github.com/jose-lico/log-processing-microservices/common/envs"
)

type RESTConfig struct {
	Env string

	Host string
	Port string

	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func NewRESTConfig() *RESTConfig {
	env := os.Getenv("ENV")

	return &RESTConfig{
		Env: env,

		Host:             os.Getenv("HOST"),
		Port:             os.Getenv("PORT"),
		AllowedOrigins:   strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowedMethods:   strings.Split(os.Getenv("ALLOWED_METHODS"), ","),
		AllowedHeaders:   strings.Split(os.Getenv("ALLOWED_HEADERS"), ","),
		ExposedHeaders:   strings.Split(os.Getenv("EXPOSED_HEADERS"), ","),
		AllowCredentials: envs.GetEnvAsBool("ALLOW_CREDENTIALS"),
		MaxAge:           envs.GetEnvAsInt("MAX_AGE"),
	}
}
