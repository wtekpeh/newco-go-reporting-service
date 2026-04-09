package config

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() *Config {
	return &Config{
		AppPort:    "8081",
		DBHost:     "127.0.0.1",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "orbi@2020",
		DBName:     "newco_db",
	}
}
