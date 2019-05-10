package noiibat

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func (s *Noiibat) ApplyHandler(name string, h http.Handler) http.Handler {
	var MapHandlers = map[string]func(*Noiibat, http.Handler) http.Handler{
		"timer":   (*Noiibat).timer,
		"traceID": (*Noiibat).traceID,
	}

	return MapHandlers[name](s, h)
}

func (s *Noiibat) traceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		log.Printf("%s Trace ID: %s", s.Name, uuid)
		next.ServeHTTP(w, r)
	})
}

func (s *Noiibat) timer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		log.Printf("%s Took: %s", s.Name, elapsed)
	})
}

func FinalHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Hit:", r.URL)
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte("completed\n"))
}
