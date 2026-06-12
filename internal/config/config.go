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

	OpenAIBaseURL string
	OpenAIAPIKey  string
	OpenAIModel   string

	TavilyBaseURL string
	TavilyAPIKey  string
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

		OpenAIBaseURL: getEnv(
			"OPENAI_BASE_URL",
			"https://api.openai.com/v1/chat/completions",
		),

		OpenAIAPIKey: getEnv(
			"OPENAI_API_KEY",
			"",
		),

		OpenAIModel: getEnv(
			"OPENAI_MODEL",
			"gpt-4.1-mini",
		),

		TavilyBaseURL: getEnv(
			"TAVILY_BASE_URL",
			"https://api.tavily.com/search",
		),

		TavilyAPIKey: getEnv(
			"TAVILY_API_KEY",
			"",
		),
	}
}
