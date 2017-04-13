package main

import (
	"sync"

	"github.com/go-redis/redis"
)

// Terminal is an exported
type Terminal struct {
	UUID      string
	Number    string
	Merchant  string
	ChannelID string
}

// Channel is an exported
type Channel struct {
	UUID      string
	Terminals map[string]*Terminal
	Sub       *redis.PubSub
	sync.Mutex
}
