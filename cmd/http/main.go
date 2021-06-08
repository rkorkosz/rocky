package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rkorkosz/rocky/pkg/messaging"
)

func main() {
	b := messaging.NewBroadcast()
	defer b.Close()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	srv := &http.Server{
		Addr:    ":8000",
		Handler: NewHandler(b),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	ctxS, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()
	err := srv.Shutdown(ctxS)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Handler is http handler
type Handler struct {
	b io.Writer
}

// NewHandler creates Handler instance
func NewHandler(b *messaging.Broadcast) *Handler {
	return &Handler{b}
}

// ServeHTTP implements http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(400)
		return
	}
	defer r.Body.Close()
	m := messaging.NewMsg(r.URL.Path)
	io.Copy(&m, r.Body)
	log.Println(m)
	json.NewEncoder(h.b).Encode(m)
	w.WriteHeader(202)
}
