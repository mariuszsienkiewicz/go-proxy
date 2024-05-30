package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-proxy/modules/log"
	"go.uber.org/zap"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(host string, port int, password string, database int) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     createAddr(host, port),
		Password: password,
		DB:       database,
	})

	return &RedisCache{
		client: rdb,
	}, nil
}

func (c *RedisCache) Set(key string, value string) {
	log.Logger.Debug("Set cache", zap.String("type", "redis"), zap.String("key", key), zap.String("value", value))
	err := c.client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		log.Logger.Error(
			"Set error",
			zap.String("type", "redis"),
			zap.String("key", key),
			zap.String("value", value),
			zap.NamedError("error", err),
		)
	}
}

func (c *RedisCache) Get(key string) (string, bool) {
	val, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Logger.Debug("Key not found", zap.String("type", "redis"), zap.String("key", key))
			return "", false
		}
		log.Logger.Warn("Get cache error", zap.String("type", "redis"), zap.String("key", key), zap.Error(err))
		return "", false
	}

	log.Logger.Debug("Get cache", zap.String("type", "redis"), zap.String("key", key), zap.String("value", val))
	return val, true
}

func (c *RedisCache) Delete(key string) {
	log.Logger.Debug("Delete cache", zap.String("type", "redis"), zap.String("key", key))
	err := c.client.Del(context.Background(), key).Err()
	if err != nil {
		log.Logger.Warn(
			"Delete error, key: %s, reason: %v",
			zap.String("type", "redis"),
			zap.String("key", key),
			zap.Error(err),
		)
	}
}

func (c *RedisCache) Clear() {
	log.Logger.Debug("Clear cache", zap.String("type", "redis"))
	err := c.client.FlushAll(context.Background()).Err()
	if err != nil {
		log.Logger.Warn(
			"Clear error",
			zap.String("type", "redis"),
			zap.Error(err),
		)
	}
}

func (c *RedisCache) Has(key string) bool {
	result, err := c.client.Exists(context.Background(), key).Result()
	if err != nil {
		log.Logger.Warn("Has error", zap.Error(err))
		return false
	}

	return result > 0
}

func createAddr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
