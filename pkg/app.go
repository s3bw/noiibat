package noiibat

import (
	"context"
	"time"

	"github.com/foxyblue/noiibat/config"
	log "github.com/sirupsen/logrus"
)

type App struct {
	config *config.Config
	bats   []*Noiibat
}

func (app *App) Start() error {
	for _, bat := range app.bats {
		bat.ListenAndServe()
	}
	time.Sleep(50000 * time.Millisecond)
	log.Println("Session complete!")
	return nil
}

func NewApp(ctx context.Context, config *config.Config) (*App, error) {
	var servers []*Noiibat

	for idx, listenConfig := range config.Listeners {
		name := string(65 + idx)
		noiibat := NewNoiibat(name, listenConfig, config)
		noiibat.RegisterHandlers()
		servers = append(servers, noiibat)
	}

	log.Printf("Starting %d servers.", len(servers))
	return &App{
		config: config,
		bats:   servers,
	}, nil

}
