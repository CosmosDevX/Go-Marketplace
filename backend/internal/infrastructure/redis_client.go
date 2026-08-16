package infrastructure

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(ctx context.Context, clientHost, clientPassword string) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:        clientHost,
		Password:    clientPassword,
		DB:          0,
		PoolSize:    10,
		PoolTimeout: 30 * time.Second,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	return RedisClient{
		client: client,
	}
}

func (c RedisClient) GetClient() *redis.Client {
	return c.client
}

func (c RedisClient) Shutdown() error {
	if err := c.client.Close(); err != nil {
		return err
	}

	return nil
}
