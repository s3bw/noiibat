package noiibat

import (
	"fmt"
	"net/http"

	animalhash "github.com/foxyblue/animal-hash/animal-hash"
	"github.com/foxyblue/noiibat/config"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Noiibat is a listener on a single port
type Noiibat struct {
	Name   string
	Port   int
	Host   string
	router http.Handler
	config *config.Listener
	app    *config.Config
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

func (s *Noiibat) RegisterHandlers() {
	router := s.router
	for _, name := range s.config.Handlers {
		router = s.ApplyHandler(name, router)
	}
	s.router = router
}

func NewNoiibat(n string, config *config.Listener, app *config.Config) *Noiibat {
	router := mux.NewRouter()
	router.Handle("/target", http.HandlerFunc(FinalHandler))

	name := animalhash.Hash(n, app.HashSeed)
	return &Noiibat{
		Name:   name,
		Port:   config.Port,
		Host:   app.Address,
		config: config,
		router: router,
	}
}
