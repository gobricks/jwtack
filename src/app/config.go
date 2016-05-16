package app

import (
	"os"
)

const (
	defaultPort = "36801"
)

type AppConfig struct {
	Port string
}

func LoadConfig() AppConfig {
	return AppConfig{envString("JWTACK_PORT", defaultPort)}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}