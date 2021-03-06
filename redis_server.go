package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func startRedis() {

	orderClient.redisClient = redis.NewClient(&redis.Options{
		Addr:     "order-sse-redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go func() {
		for v, r := pingRegis(); r != nil; v, r = pingRegis() {
			fmt.Println(v, r)
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

// executa teste de conexao com Redis
func pingRegis() (string, error) {
	return orderClient.redisClient.Ping().Result()
}
