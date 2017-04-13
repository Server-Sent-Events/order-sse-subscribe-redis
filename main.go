package main

import "github.com/go-redis/redis"

var (
	gRedisClient *redis.Client
	gChannels    map[string]*Channel
)

func main() {

	gChannels = make(map[string]*Channel)

	startRedis()

	startMux()
}
