package config

import (
 "log"
 "os"
)

type Config struct {
 DBUser     string
 DBPassword string
 DBName     string
 DBHost     string
 DBPort     string
}

func LoadConfig() Config {
 return Config{
  DBUser:     getEnv("DB_USER", "admin"),
  DBPassword: getEnv("DB_PASSWORD", "secret"),
  DBName:     getEnv("DB_NAME", "demo"),
  DBHost:     getEnv("DB_HOST", "localhost"),
  DBPort:     getEnv("DB_PORT", "5432"),
 }
}

func getEnv(key, fallback string) string {
 if value := os.Getenv(key); value != "" {
  return value
 }
 return fallback
}
