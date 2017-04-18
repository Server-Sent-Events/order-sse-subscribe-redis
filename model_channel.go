package main

import "sync"

// Channel is an exported
type Channel struct {
	UUID      string
	Terminals map[string]*Terminal
	sync.Mutex
}

// NewChannel default
func NewChannel() *Channel {
	return &Channel{Terminals: make(map[string]*Terminal)}
}
