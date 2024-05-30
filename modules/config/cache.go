package config

import (
	"errors"
)

type Cache struct {
	Type   string `yaml:"type"`
	Redis  Redis  `yaml:"redis,omitempty"`
	Memory Memory `yaml:"memory,omitempty"`
}

type Redis struct {
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Password string `yaml:"password,omitempty"`
	Database int    `yaml:"database,omitempty"`
}

type Memory struct {
	Capacity int `yaml:"capacity"`
}

func GetDefaultCache() Cache {
	return Cache{
		Redis: Redis{
			Host:     "127.0.0.1", // default Redis host
			Port:     6379,        // default Redis port
			Password: "",          // default Redis password
			Database: 0,           // default Redis Database
		},
	}
}

func ValidateCacheConfiguration() error {
	if Config.Proxy.Cache.Type == "" {
		return errors.New("cache type is required")
	}

	if Config.Proxy.Cache.Type != "redis" && Config.Proxy.Cache.Type != "memory" {
		return errors.New("cache type is invalid")
	}

	if Config.Proxy.Cache.Type == "memory" && (Config.Proxy.Cache.Memory.Capacity == 0) {
		return errors.New("cache capacity is required or cannot be 0")
	}

	return nil
}
