package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Token    string
	Channels map[string]map[string]Channel `yaml:"channels"`
}

type Channel struct {
	ID   string `yaml:"id"`
	Link string `yaml:"link"`
}

func New(configPath string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("can't read config file: %w", err)
	}

	return cfg, nil
}
