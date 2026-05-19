package redisdb

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func Connect() (*redis.Client, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env dosyası yüklenemedi, ortam değişkenleri kullanılacak")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	fmt.Println("Redis bağlantısı başarılı")
	return client, nil
}
