package config

import (
	"os"
)

type Config struct {
	Port           string
	DBPath         string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioUseSSL    bool
	MinioBucket    string
	JWTSecret      string
	AdminEmail     string
	AdminPassword  string
}

func LoadConfig() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBPath:         getEnv("DB_PATH", "data/photos.db"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "127.0.0.1:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		MinioBucket:    getEnv("MINIO_BUCKET", "deepphotos"),
		JWTSecret:      getEnv("JWT_SECRET", "deepphotos-secret-key-2026"),
		AdminEmail:     getEnv("ADMIN_EMAIL", "admin@deepphotos.local"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "deepphotos2026"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
