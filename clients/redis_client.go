package clients

import (
	"context"
	"sample-web/configs"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg configs.RedisConfig) (*RedisClient, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeoutInSeconds)*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
		DB:   cfg.Database,
	})

	if cfg.AuthEnabled {
		client.Options().Username = cfg.Username
		client.Options().Password = cfg.Password
	}

	err := client.Ping(ctx).Err()

	if err != nil {
		return nil, err
	}

	return &RedisClient{
		Client: client,
	}, nil
}
