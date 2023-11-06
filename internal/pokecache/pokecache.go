package pokecache

import (
	_"fmt"
	"log"
	"sync"
	"time"
)

// a cache entry
type CacheEntry struct {
	CreatedAt time.Time
	Val []byte
}

// cache type has a mutex to control access in goroutines
type Cache struct {
	Data map[string]CacheEntry
	sync.Mutex
}

// function to create cache
func NewCache(interval time.Duration) *Cache {
	nc :=  &Cache{
		Data: make(map[string]CacheEntry),
	}

	log.Println("Created cache, clearing it after", interval, "seconds")

	go nc.reapLoop(interval)

	return nc
}

// method to add a cacheEntry to cache
func (c *Cache) Add(key string, val []byte) {
	c.Lock()
	defer c.Unlock()

	entry := CacheEntry {
		CreatedAt: time.Now(),
		Val: val,
	}

	c.Data[key] = entry
}

// Get cache item
func (c *Cache) Get(key string) ([]byte, bool) {
	c.Lock()
	defer c.Unlock()

	if entry, present := c.Data[key]; present {
		return entry.Val, present
	}

	return nil, false
}

// clears cache based on time elapsed
func (c *Cache) reapLoop(i time.Duration) {
	ticker := time.NewTicker(i)

	go func(){
		for range ticker.C {
			c.Lock()
			for key, item := range c.Data {
				elapsed := time.Since(item.CreatedAt)

				if elapsed > i { delete(c.Data, key) }
			}
			c.Unlock()
		}
	}()
}



