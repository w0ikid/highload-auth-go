package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv    string
	LogLevel  string
	HTTP      HTTPConfig
	Postgres  PostgresConfig
	Dragonfly DragonflyConfig
	JWT       JWTConfig
}

type HTTPConfig struct {
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type DragonflyConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.DBName,
		p.SSLMode,
	)
}

func Load(prefix ...string) Config {
	_ = godotenv.Load()

	var servicePrefix string
	if len(prefix) > 0 {
		servicePrefix = prefix[0]
	}

	return Config{
		AppEnv:   getEnv("APP_ENV", "dev", servicePrefix),
		LogLevel: getEnv("LOG_LEVEL", "info", servicePrefix),
		HTTP: HTTPConfig{
			Port: getEnv("APP_PORT", "8080", servicePrefix),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost", servicePrefix),
			Port:     getEnv("POSTGRES_PORT", "5432", servicePrefix),
			User:     getEnv("POSTGRES_USER", "postgres", servicePrefix),
			Password: getEnv("POSTGRES_PASSWORD", "postgres", servicePrefix),
			DBName:   getEnv("POSTGRES_DB_NAME", "auth_db", servicePrefix),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable", servicePrefix),
		},
		Dragonfly: DragonflyConfig{
			Host:     getEnv("DRAGONFLY_HOST", "localhost", servicePrefix),
			Port:     getEnv("DRAGONFLY_PORT", "6379", servicePrefix),
			Password: getEnv("DRAGONFLY_PASSWORD", "", servicePrefix),
			DB:       getEnvInt("DRAGONFLY_DB", "0", servicePrefix),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "super-secret-key-replace-me", servicePrefix),
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 30 * 24 * time.Hour,
		},
	}
}

// SMART getEnv: first check service specific variable (e.g. AUTH_POSTGRES_HOST),
// if not found, check general variable (e.g. POSTGRES_HOST),
// if not found, return fallback
func getEnv(key, fallback string, prefix ...string) string {
	if len(prefix) > 0 && prefix[0] != "" {
		fullKey := fmt.Sprintf("%s_%s", prefix[0], key)
		if value := os.Getenv(fullKey); value != "" {
			return value
		}
	}

	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvSlice(key, fallback string, prefix ...string) []string {
	value := getEnv(key, fallback, prefix...)
	return strings.Split(value, ",")
}

func getEnvInt(key, fallback string, prefix ...string) int {
	value := getEnv(key, fallback, prefix...)
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}

func getEnvBool(key, fallback string, prefix ...string) bool {
	value := getEnv(key, fallback, prefix...)
	return value == "true" || value == "1"
}