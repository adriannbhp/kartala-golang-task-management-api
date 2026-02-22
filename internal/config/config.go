package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Secret   SecretConfig
	Auth     AuthConfig
}

type AppConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	DBNameTest string
	SSLMode    string
}

type SecretConfig struct {
	JwtSecret string
	ApiKey    string
}

type AuthConfig struct {
	AccessTokenDuration  string
	RefreshTokenDuration string
}

func LoadConfig() (Config, error) {
	// Skip .env loading in tests to ensure isolation, unless explicitly requested
	if os.Getenv("GO_ENV") != "test" || os.Getenv("LOAD_ENV_IN_TEST") == "true" {
		// Try to load .env from current directory or parent directories
		err := godotenv.Load()
		if err != nil {
			// Try various relative paths to find the .env file
			paths := []string{"../../.env", "../../../.env", "../../../../.env"}
			loaded := false
			for _, path := range paths {
				if err := godotenv.Load(path); err == nil {
					log.Println("Loaded .env from", path)
					loaded = true
					break
				}
			}
			if !loaded {
				log.Println("No .env file found in parent directories, using environment variables")
			}
		} else {
			log.Println("Loaded .env from current directory")
		}
	}

	cfg := Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "users_db"),
			DBNameTest: getEnv("DB_NAME_TEST", "kartala_test"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Secret: SecretConfig{
			JwtSecret: getEnv("JWT_SECRET", ""),
			ApiKey:    getEnv("API_KEY", ""),
		},
		Auth: AuthConfig{
			AccessTokenDuration:  getEnv("ACCESS_TOKEN_DURATION", "15m"),
			RefreshTokenDuration: getEnv("REFRESH_TOKEN_DURATION", "168h"),
		},
	}

	if cfg.Secret.JwtSecret == "" {
		if os.Getenv("GO_ENV") == "test" && os.Getenv("LOAD_ENV_IN_TEST") != "true" {
			cfg.Secret.JwtSecret = "test_jwt_secret_only_for_testing_purposes_123"
		} else {
			return cfg, errors.New("JWT_SECRET is required in .env file")
		}
	}

	if cfg.Secret.ApiKey == "" {
		if os.Getenv("GO_ENV") == "test" && os.Getenv("LOAD_ENV_IN_TEST") != "true" {
			cfg.Secret.ApiKey = "test_api_key_123"
		} else {
			return cfg, errors.New("API_KEY is required in .env file")
		}
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
