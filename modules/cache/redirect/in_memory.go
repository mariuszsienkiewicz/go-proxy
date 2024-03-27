package redirect

import (
	"proxy/modules/config"
	"proxy/modules/log"
)

type InMemoryCache struct {
	cache map[string]config.Server
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string]config.Server),
	}
}

func (c InMemoryCache) Add(hash string, server config.Server) {
	log.Logger.Tracef("Add Redirect To Cache (Hash %s, Target: %v)", hash, server.Id)
	c.cache[hash] = server
}

func (c InMemoryCache) Find(hash string) (config.Server, bool) {
	s, ok := c.cache[hash]
	log.Logger.Tracef("Find Redirect In Cache (Hash: %s, Found: %v, Target: %v)", hash, ok, s.Id)
	return s, ok
}
