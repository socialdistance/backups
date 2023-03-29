package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	HTTP HTTPConf
	File FileName
}

type HTTPConf struct {
	TargetUrl string `json:"target_url"`
}

type FileName struct {
	FileNameBackup string `json:"file_name_backup"`
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
