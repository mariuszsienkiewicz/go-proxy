package cache

import (
	"container/list"
	"fmt"
	"go-proxy/modules/log"
	"sync"
)

type InMemoryCache struct {
	cache    map[string]*list.Element
	lruList  *list.List
	capacity int
	mu       sync.Mutex // Mutex to ensure thread safety
}

type Entry struct {
	key   string
	value string
}

func NewInMemoryCache(capacity int) (Cache, error) {
	if capacity == 0 {
		return nil, fmt.Errorf("NewInMemoryCache capacity can't be 0")
	}

	return &InMemoryCache{
		cache:    make(map[string]*list.Element),
		lruList:  list.New(),
		capacity: capacity,
	}, nil
}

// Set stores a value associated with the given key in the cache.
func (c *InMemoryCache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the key already exists
	if elem, ok := c.cache[key]; ok {
		c.lruList.MoveToFront(elem)
		elem.Value.(*Entry).value = value
	}

	// If the cache is at capacity, remove the least recently used item
	if c.lruList.Len() >= c.capacity {
		log.Logger.Debug("LRU capacity hit, removing last item")
		backElem := c.lruList.Back()
		if backElem != nil {
			c.lruList.Remove(backElem)
			delete(c.cache, backElem.Value.(*Entry).key)
		}
	}

	// Add the new entry to the front of the LRU list and to the cache map
	entry := &Entry{
		key:   key,
		value: value,
	}

	frontElem := c.lruList.PushFront(entry)
	c.cache[key] = frontElem
}

// Get retrieves the value associated with the given key from the cache.
func (c *InMemoryCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.lruList.MoveToFront(elem)
		return elem.Value.(*Entry).value, true
	}
	return "", false
}

// Delete removes the value associated with the given key from the cache.
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.lruList.Remove(elem)
		delete(c.cache, key)
	}
}

// Clear removes all entries from the cache.
func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*list.Element)
	c.lruList.Init()
}

// Has checks if a given key exists in the cache.
func (c *InMemoryCache) Has(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cache[key]
	return ok
}
