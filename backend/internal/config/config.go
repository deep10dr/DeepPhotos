package config

import (
	"log"
	"os"
	"strings"
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

// loadEnvFile searches for a .env file and sets variables in the process environment if not already set.
func loadEnvFile() {
	paths := []string{".env", "../.env", "../../.env", "backend/.env"}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				k := strings.TrimSpace(parts[0])
				v := strings.TrimSpace(parts[1])
				v = strings.Trim(v, `"'`)
				if _, set := os.LookupEnv(k); !set {
					os.Setenv(k, v)
				}
			}
		}
		log.Printf("Loaded environment settings from %s", p)
		break
	}
}

func LoadConfig() *Config {
	loadEnvFile()

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
