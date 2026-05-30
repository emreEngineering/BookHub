package redisdb

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, redisAddr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	fmt.Println("Redis bağlantısı başarılı")
	return client, nil
}
