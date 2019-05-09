package config

import (
	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Listeners []int
	Address   string
}

func New() *Config {
	var yamlExample = []byte(`
	listeners: [8990, 8991]
	address: localhost
	`)
	config := &Config{}

	err := yaml.Unmarshal(yamlExample, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}
