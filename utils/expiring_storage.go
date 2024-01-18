package utils

import (
	"context"
	"sync"
	"time"
)

type ExpiringStorage[K comparable, V any] struct {
	mux               sync.RWMutex
	values            map[K]expiringValue[V]
	expirationTimeout time.Duration
	cleanupInterval   time.Duration
}

type expiringValue[V any] struct {
	value          V
	expirationTime time.Time
}

func NewExpiringStorage[K comparable, V any](expirationTimeout time.Duration, cleanupInterval time.Duration) *ExpiringStorage[K, V] {
	return &ExpiringStorage[K, V]{
		values:            make(map[K]expiringValue[V]),
		expirationTimeout: expirationTimeout,
		cleanupInterval:   cleanupInterval,
	}
}

func (e *ExpiringStorage[K, V]) Add(key K, value V) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.values[key] = expiringValue[V]{
		value:          value,
		expirationTime: time.Now().Add(e.expirationTimeout),
	}
}

func (e *ExpiringStorage[K, V]) Get(key K) (value V) {
	e.mux.RLock()
	defer e.mux.RUnlock()

	return e.values[key].value
}

func (e *ExpiringStorage[K, V]) Len() int {
	e.mux.RLock()
	defer e.mux.RUnlock()

	return len(e.values)
}

func (e *ExpiringStorage[K, V]) cleanup() {
	now := time.Now()
	e.mux.Lock()
	defer e.mux.Unlock()

	for key := range e.values {
		if e.values[key].expirationTime.Before(now) {
			delete(e.values, key)
		}
	}
}

func (e *ExpiringStorage[K, V]) Run(ctx context.Context) {
	ticker := time.NewTicker(e.cleanupInterval)
	for {
		select {
		case <-ticker.C:
			e.cleanup()
		case <-ctx.Done():
			return
		}
	}
}
