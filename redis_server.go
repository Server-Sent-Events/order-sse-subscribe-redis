package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func startRedis() {

	gRedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
	return gRedisClient.Ping().Result()
}
