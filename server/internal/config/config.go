package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Storage    StorageConf
	HTTP       HttpConf
	Cache      CacheConf
	WorkerPool WorkerPoolConf
}

type WorkerPoolConf struct {
	NumWorkers  int `json:"numWorkers"`
	ChannelSize int `json:"channelSize"`
}

type StorageConf struct {
	Type string `json:"type"`
	Dsn  string `json:"dsn"`
}

type HttpConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type CacheConf struct {
	DefaultExpiration int `json:"defaultExpiration"`
	CleanupInterval   int `json:"cleanupInterval"`
}

func NewConfig() Config {
	return Config{}
}

func LoadConfig(path string) (*Config, error) {
	resultConfig, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("invalid config %s: %w", path, err)
	}

	config := NewConfig()
	err = json.Unmarshal(resultConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("invalid unmarshal config %s:%w", path, err)
	}

	return &config, nil
}
