package config

import (
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

// DefaultHandlers will execute from right -> left
var DefaultHandlers = []string{"timer", "traceID", "context"}

var yamlExample = []byte(`
listeners:
- port: 8990
  handlers: ['respond']
  context_time_out: 3
  delay: 6
- port: 8991
# - port: 8992
#   delay: 30
address: localhost
hash_seed: 17
`)

type Config struct {
	Listeners []*Listener
	Address   string
	HashSeed  int `yaml:"hash_seed"`
}

type Listener struct {
	Port           int
	Delay          time.Duration
	ContextTimeOut time.Duration
	Handlers       []string
}

func New() *Config {
	// Provide a default config
	config := &Config{
		Listeners: []*Listener{
			{
				Port:  8990,
				Delay: 0,
			},
		},
	}

	err := yaml.Unmarshal(yamlExample, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Add default handlers to config
	for _, listener := range config.Listeners {
		for _, defaultH := range DefaultHandlers {
			listener.Handlers = append([]string{defaultH}, listener.Handlers...)
		}

		// If ContextTimeOut is set to `0` change this to 60sec
		if listener.ContextTimeOut == 0 {
			listener.ContextTimeOut = 60
		}
	}

	return config
}
