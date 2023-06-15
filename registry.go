package main

import "sync"

type registry struct {
	registered   map[string]int
	disqualified map[string]bool
	lock         sync.RWMutex
}

func makeRegistry() registry {
	return registry{registered: make(map[string]int), disqualified: make(map[string]bool), lock: sync.RWMutex{}}
}

func (r *registry) close() {
	// Not yet implemented
}

func (r *registry) registerPlayer(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.registered[key] = 0
}

func (r *registry) deregisterPlayer(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.registered, key)
}

func (r *registry) isRegistered(key string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.registered[key]

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
	keys := make([]string, len(r.registered))

	i := 0
	for k := range r.registered {
		keys[i] = k
		i++
	}

	return keys
}

func (r *registry) recordWin(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.registered[key] += 1
}
