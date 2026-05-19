package config

import "os"

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	KeycloakServerURL string
	KeycloakRealm     string
	KeycloakClientID  string

	OllamaBaseURL string
	OllamaModel   string
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func Load() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", "8081"),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "newco_db"),

		KeycloakServerURL: getEnv("KEYCLOAK_SERVER_URL", ""),
		KeycloakRealm:     getEnv("KEYCLOAK_REALM", ""),
		KeycloakClientID:  getEnv("KEYCLOAK_CLIENT_ID", ""),

		OllamaBaseURL: getEnv("OLLAMA_BASE_URL", "http://127.0.0.1:11434"),
		OllamaModel:   getEnv("OLLAMA_MODEL", "qwen3:4b"),
	}
}
