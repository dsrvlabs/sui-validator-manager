package config

import (
	"errors"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// TODO: Handle singleton pattern for config instance.
var (
	config     *Config
	configLock = sync.Mutex{}
)

// Config contains all configurations of manager.
type Config struct {
	RPC []struct {
		Name     string `yaml:"name"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"rpc"`
}

func parseConfig(data []byte) (*Config, error) {
	newConfig := new(Config)
	err := yaml.Unmarshal(data, newConfig)
	if err != nil {
		return nil, err
	}

	return newConfig, nil
}

// Load loads file contents and parse YAML content into structure.
func Load(file string) (*Config, error) {
	configLock.Lock()
	defer configLock.Unlock()

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config, err = parseConfig(data)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Get returns parsed config
func Get() (*Config, error) {
	configLock.Lock()
	defer configLock.Unlock()

	if config == nil {
		return nil, errors.New("config not loaded")
	}

	return config, nil
}
