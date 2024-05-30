package cache

import (
	"fmt"
	"go-proxy/modules/config"
	"go-proxy/modules/log"
)

// Cache is an interface defining the main functions for a cache system.
type Cache interface {
	// Set stores a value associated with the given key in the cache
	Set(key string, value string)

	// Get retrieves the value associated with the given key from the cache.
	// Returns the value and a boolean indicating whether the key was found.
	Get(key string) (string, bool)

	// Delete removes the value associated with the given key from the cache.
	Delete(key string)

	// Clear removes all entries from the cache.
	Clear()

	// Has checks if a given key exists in the cache.
	Has(key string) bool
}

var initializedCache Cache

// InitializeRedisCache initializes the Redis cache
func InitializeRedisCache(cfg config.Redis) (Cache, error) {
	log.Logger.Debug("Initializing Redis cache")
	return NewRedisCache(cfg.Host, cfg.Port, cfg.Password, cfg.Database)
}

// InitializeInMemoryCache initializes the in-memory cache
func InitializeInMemoryCache(cfg config.Memory) (Cache, error) {
	log.Logger.Debug("Initializing memory cache")
	return NewInMemoryCache(cfg.Capacity)
}

// InitCache initializes the cache based on the configuration
func InitCache() error {
	var err error
	switch config.Config.Proxy.Cache.Type {
	case "redis":
		initializedCache, err = InitializeRedisCache(config.Config.Proxy.Cache.Redis)
	case "memory":
		initializedCache, err = InitializeInMemoryCache(config.Config.Proxy.Cache.Memory)
	default:
		err = fmt.Errorf("unsupported cache type: %s", config.Config.Proxy.Cache.Type)
	}

	if err != nil {
		return err
	}

	return nil
}

// GetCache returns the initialized cache instance
func GetCache() Cache {
	return initializedCache
}
