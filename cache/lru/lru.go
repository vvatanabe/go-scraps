package lru

import (
	"sync"

	"github.com/golang/groupcache/lru"
)

// Cache is an LRU cache. It is safe for concurrent access. Just a wrapper for the groupcache/lru
type Cache struct {
	lru  lru.Cache
	lock sync.RWMutex
}

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key = lru.Key

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New(maxEntries int) *Cache {
	return &Cache{
		lru: *lru.New(maxEntries),
	}
}

// NewWithEvict creates a new Cache with onEvicted as callback.
func NewWithEvict(maxEntries int, onEvicted func(key Key, value interface{})) *Cache {
	c := New(maxEntries)
	c.lru.OnEvicted = onEvicted
	return c
}

// Add adds a value to the cache.
func (c *Cache) Add(key Key, value interface{}) {
	c.lock.Lock()
	c.lru.Add(key, value)
	c.lock.Unlock()
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	c.lock.Lock()
	c.lock.Unlock()
	return c.lru.Get(key)
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key Key) {
	c.lock.Lock()
	c.lru.Remove(key)
	c.lock.Unlock()
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache) RemoveOldest() {
	c.lock.Lock()
	c.lru.RemoveOldest()
	c.lock.Unlock()
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Len()
}

// Clear purges all stored items from the cache.
func (c *Cache) Clear() {
	c.lock.Lock()
	c.lru.Clear()
	c.lock.Unlock()
}
