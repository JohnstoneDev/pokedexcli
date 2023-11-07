package pokecache

import (
	"log"
	"sync"
	"time"
	"fmt"

	"github.com/mtslzr/pokeapi-go/structs"
)

// a cache entry
type CacheEntry struct {
	CreatedAt time.Time
	Val structs.Resource
}

// cache type has a mutex to control access in goroutines
type Cache struct {
	Data map[string]CacheEntry
	sync.Mutex
}

// function to create cache, accepts a interval parameter
// that will be used to clear the cache
func NewCache(interval time.Duration) *Cache {
	newCache :=  &Cache{
		Data: make(map[string]CacheEntry),
	}

	log.Println("Created cache, clearing it after", interval)

	go newCache.reapLoop(interval)

	return newCache
}

// method to add a cacheEntry to cache
func (c *Cache) Add(key string, val structs.Resource) {
	c.Lock()
	defer c.Unlock()

	entry := CacheEntry {
		CreatedAt: time.Now(),
		Val: val,
	}

	c.Data[key] = entry
}

// Get cache item
func (c *Cache) Get(key string) (structs.Resource, bool) {
	c.Lock()
	defer c.Unlock()

	if entry, present := c.Data[key]; present {
		return entry.Val, present
	}

	return structs.Resource{}, false
}

// clears cache based on time elapsed
func (c *Cache) reapLoop(i time.Duration) {
	// ticker to create a channel that will
	// send the the interval to the channel after each tick
	ticker := time.NewTicker(i)

	// go routine that clears the cache
	go func(){
		for range ticker.C {
			c.Lock()
			for key, item := range c.Data {
				elapsed := time.Since(item.CreatedAt)

				if elapsed > i { delete(c.Data, key) }

				log.Println("cache cleared")

				fmt.Print("Pokedex > ")
			}
			c.Unlock()
		}
	}()
}
