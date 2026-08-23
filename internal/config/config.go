package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env      string
	HTTP     HTTPConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type HTTPConfig struct {
	Addr string
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret string
}

type LogConfig struct {
	Level slog.Level
}

func Load() (Config, error) {
	cfg := Config{
		Env: envString("APP_ENV", "development"),
		HTTP: HTTPConfig{
			Addr: envString("HTTP_ADDR", ":8080"),
		},
		// Database: DatabaseConfig{
		// 	URL:             strings.TrimSpace(os.Getenv("DATABASE_URL")),
		// 	MaxOpenConns:    envInt("DB_MAX_OPEN_CONNS", 10),
		// 	MaxIdleConns:    envInt("DB_MAX_IDLE_CONNS", 5),
		// 	ConnMaxLifetime: envDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		// },
		// JWT: JWTConfig{
		// 	Secret: strings.TrimSpace(os.Getenv("JWT_SECRET")),
		// },
		// Log: LogConfig{
		// 	Level: envLogLevel("LOG_LEVEL", slog.LevelInfo),
		// },
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) validate() error {
	var errs []error

	if cfg.HTTP.Addr == "" {
		errs = append(errs, errors.New("HTTP_ADDR is required"))
	}
	// if cfg.Database.URL == "" {
	// 	errs = append(errs, errors.New("DATABASE_URL is required"))
	// }
	// if cfg.JWT.Secret == "" {
	// 	errs = append(errs, errors.New("JWT_SECRET is required"))
	// }
	// if cfg.Database.MaxOpenConns <= 0 {
	// 	errs = append(errs, errors.New("DB_MAX_OPEN_CONNS must be greater than 0"))
	// }
	// if cfg.Database.MaxIdleConns <= 0 {
	// 	errs = append(errs, errors.New("DB_MAX_IDLE_CONNS must be greater than 0"))
	// }
	// if cfg.Database.MaxIdleConns > cfg.Database.MaxOpenConns {
	// 	errs = append(errs, errors.New("DB_MAX_IDLE_CONNS must not exceed DB_MAX_OPEN_CONNS"))
	// }
	// if cfg.Database.ConnMaxLifetime <= 0 {
	// 	errs = append(errs, errors.New("DB_CONN_MAX_LIFETIME must be greater than 0"))
	// }

	return errors.Join(errs...)
}

func envString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func envLogLevel(key string, fallback slog.Level) slog.Level {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return fallback
	}

	switch value {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		fmt.Fprintf(os.Stderr, "unknown %s value %q, using default log level\n", key, value)
		return fallback
	}
}
