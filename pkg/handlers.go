package noiibat

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// UUID type for traceID
type UUID string

// UUIDKey refers to the key to use to store trace ID in context
const UUIDKey UUID = "user"

func (s *Noiibat) ApplyHandler(name string) func(http.Handler) http.Handler {
	var MapHandlers = map[string]func(http.Handler) http.Handler{
		"context": s.createContext,
		"timer":   s.timer,
		"respond": s.respond,
		"traceID": s.traceID,
	}

	return MapHandlers[name]
}

func (s *Noiibat) respond(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Context().Value(UUIDKey)
		log.Printf("%s responding to %s.", s.Name, uuid)
		next.ServeHTTP(w, r)
		w.Write([]byte("completed\n"))
	})
}

// createContext add a context to the request this can be used to propagate
// values down the request chain.
func (s *Noiibat) createContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.ContextTimeOut*time.Second))
		defer cancel()
		log.Printf("%s created a context.", s.Name)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// traceID get generated, this should be generated if one does not already exist
// this will allow monitoring to follow the request through the system.
func (s *Noiibat) traceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		ctx := context.WithValue(r.Context(), UUIDKey, uuid)
		log.Printf("%s Trace ID: %s", s.Name, uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// timer will log the amount of time the request takes.
func (s *Noiibat) timer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		uuid := r.Context().Value(UUIDKey)
		log.Printf("%s Took: %s on %s", s.Name, elapsed, uuid)
	})
}

func (s *Noiibat) FinalHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s Hit: %s, has uuid: %s", s.Name, r.URL, r.Context().Value(UUIDKey))
	time.Sleep(s.config.Delay * time.Millisecond)
}
