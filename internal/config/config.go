package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr  string
	DatabaseURL string
	RedisAddr   string
	MongoURI    string
	MongoDBName string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env dosyası yüklenemedi, ortam değişkenleri kullanılacak")
	}

	return Config{
		ServerAddr:  envOrDefault("SERVER_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisAddr:   envOrDefault("REDIS_ADDR", "localhost:6379"),
		MongoURI:    envOrDefault("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: envOrDefault("MONGO_DATABASE", "bookhub"),
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
