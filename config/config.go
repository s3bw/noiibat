package config

import (
	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

var yamlExample = []byte(`
# listeners:
# - port: 8990
#   delay: 30
# - port: 8991
# - port: 8992
address: localhost
hash_seed: 17
`)

type Config struct {
	Listeners []*Listener
	Address   string
	HashSeed  int `yaml:"hash_seed"`
}

type Listener struct {
	Port     int
	Handlers []string
}

func New() *Config {
	// Provide a default config
	config := &Config{
		Listeners: []*Listener{
			{Port: 8990, Handlers: []string{"timer", "traceID"}},
		},
	}

	err := yaml.Unmarshal(yamlExample, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}
