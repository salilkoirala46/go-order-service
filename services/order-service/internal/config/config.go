package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort  string
	DatabaseDSN string
	NATSUrl     string
	JWTSecret   string
	JWTExpiry   time.Duration
}

func Load() *Config {
	hours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8081"),
		DatabaseDSN: getEnv("DATABASE_DSN", "root:rootpassword@tcp(localhost:3306)/orderdb?charset=utf8mb4&parseTime=True&loc=Local"),
		NATSUrl:     getEnv("NATS_URL", "nats://localhost:4222"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-jwt-key"),
		JWTExpiry:   time.Duration(hours) * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
