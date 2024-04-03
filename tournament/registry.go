package tournament

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Registry struct {
	connections  map[string]*websocket.Conn
	win_count    map[string]int
	disqualified map[string]bool
	lock         sync.RWMutex
}

func MakeRegistry() Registry {
	return Registry{win_count: make(map[string]int), disqualified: make(map[string]bool), lock: sync.RWMutex{}}
}

func (r *Registry) Close() {
	for _, conn := range r.connections {
		conn.Close()
	}
}

func (r *Registry) GetConnection(key string) *websocket.Conn {
	return r.connections[key]
}

func (r *Registry) RegisterPlayer(key string, conn *websocket.Conn) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.win_count[key] = 0
	r.connections[key] = conn
}

func (r *Registry) DeregisterPlayer(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.win_count, key)
	delete(r.connections, key)
}

func (r *Registry) IsRegistered(key string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.connections[key]

	return ok
}

func (r *Registry) DisqualifyPlayer(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.disqualified[key] = true
}

func (r *Registry) IsDisqualified(key string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	d, ok := r.disqualified[key]

	if !ok {
		return false
	}

	return d
}

func (r *Registry) GetRegistered() []string {
	keys := make([]string, len(r.win_count))

	i := 0
	for k := range r.connections {
		keys[i] = k
		i++
	}

	return keys
}

func (r *Registry) RecordWin(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.win_count[key] += 1
}
