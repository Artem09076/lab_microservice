package config

import (
	"os"
)

type Config struct {
	KeyCloakURL    string
	KeyCloakRealm  string
	KeyCloakClient string
	KeyCloakSecret string
	DatabaseURL    string
	Port           string
}

func Load() (*Config, error) {
	cfg := Config{
		KeyCloakURL:    os.Getenv("KEYCLOAK_URL"),
		KeyCloakRealm:  os.Getenv("KEYCLOAK_REALM"),
		KeyCloakClient: os.Getenv("KEYCLOAK_CLIENT"),
		KeyCloakSecret: os.Getenv("KEYCLOAK_SECRET"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		Port:           "10000",
	}
	return &cfg, nil
}
