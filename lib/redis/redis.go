package redis

import (
	"github.com/redis/go-redis/v9"
)

func NewClient(addr string, pw string) *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       0, //use the default Redis DB
	})
	return c
}
