package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/shordem/api.thryvo/lib/constants"
)

var ctx = context.Background()

type redisClient struct {
	client *redis.Client
}

type RedisClientInterface interface {
	Set(key string, value interface{}) error
	Get(key string, batchSize int64) ([]string, error)
}

func NewRedisClient(env constants.Env) RedisClientInterface {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.REDIS_SERVER,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	res, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	return &redisClient{
		client: rdb,
	}
}

func (c *redisClient) Set(key string, value interface{}) error {
	err := c.client.LPush(ctx, key, value).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *redisClient) Get(key string, batchSize int64) ([]string, error) {
	val, err := c.client.LRange(ctx, key, 0, batchSize-1).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}
