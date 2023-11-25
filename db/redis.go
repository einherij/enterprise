package db

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
		DB:   cfg.DB,
	})
	return rdb, nil
}
