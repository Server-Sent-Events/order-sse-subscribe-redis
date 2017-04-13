package main

import (
	"sync"

	"github.com/go-redis/redis"
)

// Terminal is an exported
type Terminal struct {
	UUID     string
	Number   string
	Merchant string
	Sub      *redis.PubSub
}

// Channel is an exported
type Channel struct {
	UUID      string
	Terminals map[string]*Terminal
	Sub       *redis.PubSub
	sync.Mutex
}

// NewChannel default
func NewChannel() *Channel {
	return &Channel{Terminals: make(map[string]*Terminal)}
}
