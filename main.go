package main

// var OrderClient = struct {
// 	sync.RWMutex
// 	redisClient *redis.Client
// 	channels    map[string]*Channel
// }{channels: make(map[string]*Channel)}
var orderClient *OrderClient

func main() {

	orderClient = &OrderClient{
		channels: make(map[string]*Channel),
	}

	startRedis()

	startMux()
}
