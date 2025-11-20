package config

import (
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
  DBUser:     getEnv("DB_USER", "postgres"),
  DBPassword: getEnv("DB_PASSWORD", "password"),
  DBName:     getEnv("DB_NAME", "PepeScale"),
  DBHost:     getEnv("DB_HOST", "172.17.0.2"),
  DBPort:     getEnv("DB_PORT", "5432"),
 }
}

func getEnv(key, fallback string) string {
 if value := os.Getenv(key); value != "" {
  return value
 }
 return fallback
}
