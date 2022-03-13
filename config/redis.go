package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var CacheChannel chan string

type NewCacheConfig interface {
	ClearCache(keys ...string)
	RedisClient() *redis.Client
}

type CacheConfig struct {
	Host   string
	Port   int
	Client *redis.Client
}

func SetupRedis() NewCacheConfig {

	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	return &CacheConfig{
		Client: client,
	}
}

func (conf CacheConfig) SetupCacheChannel(c CacheConfig) {
	CacheChannel = make(chan string)
	go func(ch chan string, c CacheConfig) {
		for {
			time.Sleep(5 * time.Second)

			key := <-ch

			c.Client.Del(context.Background(), key)

			fmt.Println("Cache cleared " + key)
		}
	}(CacheChannel, c)
}

func (conf CacheConfig) ClearCache(keys ...string) {
	for _, key := range keys {
		CacheChannel <- key
	}
}

func (conf CacheConfig) RedisClient() *redis.Client {
	return conf.Client
}
