package noiibat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/foxyblue/noiibat/config"
	"github.com/gorilla/mux"
)

type Noiibat struct {
	config  *config.Config
	servers []*http.Server
}

// ListenAndServe deploys noiibat
func (noii *Noiibat) ListenAndServe() error {
	for _, s := range noii.servers {
		log.Println("Started")
		log.Printf("Listening on %s", s.Addr)

		go func(s *http.Server) {
			s.ListenAndServe()
		}(s)
	}

	time.Sleep(50000 * time.Millisecond)
	log.Println("Session complete!")
	return nil
}

func timer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		elapsed := time.Since(start)
		log.Printf("Took: %s", elapsed)
	})
}

func TargetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Hit:", r.URL)
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte("completed\n"))
}

func Spawn(ctx context.Context, config *config.Config) (*Noiibat, error) {
	var servers []*http.Server

	router := mux.NewRouter()
	router.Handle("/target", timer(http.HandlerFunc(TargetHandler)))
	http.Handle("/", router)

	for _, port := range config.Listeners {
		address := fmt.Sprintf("%s%d", config.Address, port)
		servers = append(servers, &http.Server{
			Handler: router,
			Addr:    address,
		})
	}

	log.Println(servers)
	return &Noiibat{
		config:  config,
		servers: servers,
	}, nil
}
