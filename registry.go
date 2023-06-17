package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type registry struct {
	connections  map[string]*websocket.Conn
	win_count    map[string]int
	disqualified map[string]bool
	lock         sync.RWMutex
}

func makeRegistry() registry {
	return registry{win_count: make(map[string]int), disqualified: make(map[string]bool), lock: sync.RWMutex{}}
}

func (r *registry) close() {
	for _, conn := range r.connections {
		conn.Close()
	}
}

func (r *registry) getConnection(key string) *websocket.Conn {
	return r.connections[key]
}

func (r *registry) registerPlayer(key string, conn *websocket.Conn) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.win_count[key] = 0
	r.connections[key] = conn
}

func (r *registry) deregisterPlayer(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.win_count, key)
	delete(r.connections, key)
}

func (r *registry) isRegistered(key string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.connections[key]

	return ok
}

func (r *registry) disqualifyPlayer(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.disqualified[key] = true
}

func (r *registry) isDisqualified(key string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	d, ok := r.disqualified[key]

	if !ok {
		return false
	}

	return d
}

func (r *registry) getRegistered() []string {
	keys := make([]string, len(r.win_count))

	i := 0
	for k := range r.connections {
		keys[i] = k
		i++
	}

	return keys
}

func (r *registry) recordWin(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.win_count[key] += 1
}
