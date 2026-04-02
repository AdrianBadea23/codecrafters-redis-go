package internal

import "sync"

type RedisServer struct {
	Mu sync.RWMutex

	Data    map[string]string
	Expires map[string]int64
	Lists   map[string]any
}

func New() *RedisServer {
	return &RedisServer{
		Data:    make(map[string]string),
		Expires: make(map[string]int64),
		Lists:   make(map[string]any, 100),
	}
}
