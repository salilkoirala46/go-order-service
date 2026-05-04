package config

import "os"

type Config struct {
	ServerPort  string
	DatabaseDSN string
	NATSUrl     string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8082"),
		DatabaseDSN: getEnv("DATABASE_DSN", "root:rootpassword@tcp(localhost:3306)/notificationdb?charset=utf8mb4&parseTime=True&loc=Local"),
		NATSUrl:     getEnv("NATS_URL", "nats://localhost:4222"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-jwt-key"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
