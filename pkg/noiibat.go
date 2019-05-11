package noiibat

import (
	"context"
	"fmt"
	"net/http"
	"time"

	animalhash "github.com/foxyblue/animal-hash/animal-hash"
	"github.com/foxyblue/noiibat/config"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Noiibat is a listener on a single port
type Noiibat struct {
	Name           string
	Port           int
	Host           string
	Context        *context.Context
	ContextTimeOut time.Duration
	router         http.Handler
	config         *config.Listener
	app            *config.Config
}

// ListenAndServe a single bat
func (s *Noiibat) ListenAndServe() {
	log.Printf("%s has started", s.Name)
	log.Printf("%s is listening on %s:%d", s.Name, s.Host, s.Port)

	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	server := &http.Server{
		Handler: s.router,
		Addr:    address,
	}

	go func(s *http.Server) {
		s.ListenAndServe()
	}(server)
}

// RegisterHandlers adds middleware as defined in config
func (s *Noiibat) RegisterHandlers() {
	router := mux.NewRouter()
	router.Handle("/target", http.HandlerFunc(s.FinalHandler))

	for _, name := range s.config.Handlers {
		middleware := s.ApplyHandler(name)
		router.Use(middleware)
	}
	s.router = router
}

func NewNoiibat(n string, config *config.Listener, app *config.Config) *Noiibat {
	name := animalhash.Hash(n, app.HashSeed)
	return &Noiibat{
		Name:   name,
		Port:   config.Port,
		Host:   app.Address,
		config: config,
	}
}
