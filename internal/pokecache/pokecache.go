package pokecache

import (
	_"fmt"
	"log"
	"sync"
	"time"
)

// a cache entry
type cacheEntry struct {
	createdAt time.Time
	val []byte
}

// cache type has a mutex to control access in goroutines
type Cache struct {
	data map[string]cacheEntry
	sync.Mutex
}

// function to create cache
func NewCache(interval time.Duration) *Cache {
	nc :=  &Cache{
		data: make(map[string]cacheEntry),
	}

	log.Println("Created cache, clearing it after", interval, "seconds")

	go nc.reapLoop(interval)

	return nc
}

// method to add a cacheEntry to cache
func (c *Cache) Add(key string, val []byte) {
	c.Lock()
	defer c.Unlock()

	entry := cacheEntry {
		createdAt: time.Now(),
		val: val,
	}

	c.data[key] = entry
}

// Get cache item
func (c *Cache) Get(key string) ([]byte, bool) {
	c.Lock()
	defer c.Unlock()

	if entry, present := c.data[key]; present {
		return entry.val, present
	}

	return nil, false
}

// clears cache based on time elapsed
func (c *Cache) reapLoop(i time.Duration) {
	ticker := time.NewTicker(i)

	go func(){
		for range ticker.C {
			c.Lock()
			for key, item := range c.data {
				elapsed := time.Since(item.createdAt)

				if elapsed > i { delete(c.data, key) }
			}
			c.Unlock()
		}
	}()
}



