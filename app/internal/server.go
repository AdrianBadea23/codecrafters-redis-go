package internal

import "sync"

type streamStruct struct {
	ID     string
	Fields map[string]any
}

type RedisServer struct {
	Mu sync.Mutex

	Data     map[string]string
	Expires  map[string]int64
	Lists    map[string]any
	Channels map[string][]chan string
	Streams  map[string][]streamStruct
}

func New() *RedisServer {
	return &RedisServer{
		Mu: sync.Mutex{},

		Data:     make(map[string]string),
		Expires:  make(map[string]int64),
		Lists:    make(map[string]any, 100),
		Channels: make(map[string][]chan string),
		Streams:  make(map[string][]streamStruct),
	}
}
