package redirect

import (
	"proxy/modules/log"
)

type InMemoryCache struct {
	cache map[string]string
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string]string),
	}
}

func (c InMemoryCache) Add(hash string, target string) {
	log.Logger.Tracef("Add Redirect To Cache (Hash %s, Target: %v)", hash, target)
	c.cache[hash] = target
}

func (c InMemoryCache) Find(hash string) (string, bool) {
	t, ok := c.cache[hash]
	log.Logger.Tracef("Find Redirect In Cache (Hash: %s, Found: %v, Target: %v)", hash, ok, t)
	return t, ok
}

func (c InMemoryCache) Clear() {
	c.cache = make(map[string]string)
	log.Logger.Tracef("Clear Redirect In Cache")
}
