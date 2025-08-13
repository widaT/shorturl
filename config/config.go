package config

import "time"

var config *Config

func Get() *Config {
	return config
}

func StartUp() {
	config = &Config{}
}

type Config struct {
	// Server
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}

	Scenes map[string]struct {
		TTL  time.Duration `yaml:"ttl"`
		Host string        `yaml:"host"`
	}
}
