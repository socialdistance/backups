package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	HTTP HTTPConf
}

type HTTPConf struct {
	TargetUrl string `json:"targetUrl"`
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
