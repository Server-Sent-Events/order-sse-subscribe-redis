package main

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// OrderClient is an exported
type OrderClient struct {
	sync.RWMutex
	redisClient *redis.Client
	channels    map[string]*Channel
}

// Channel is an exported
type Channel struct {
	UUID       string `json:"id"`
	MerchantID string
	PIN        string `json:"pin"`
	Terminals  map[string]*Terminal
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	sync.Mutex
}

// NewChannel default
func NewChannel() *Channel {
	return &Channel{Terminals: make(map[string]*Terminal)}
}
