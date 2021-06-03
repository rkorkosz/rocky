package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rkorkosz/rocky/pkg/messaging"
)

func main() {
	b := messaging.NewBroadcast()
	defer b.Close()
	log.Fatal(http.ListenAndServe(":8000", NewHandler(b)))
}

// Handler is http handler
type Handler struct {
	b *messaging.Broadcast
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
	m := NewMsg(r.URL.Path)
	io.Copy(&m, r.Body)
	log.Println(m)
	json.NewEncoder(h.b).Encode(m)
	w.WriteHeader(202)
}

// Msg carries data
type Msg struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// Write implements io.Writer interface
func (m *Msg) Write(p []byte) (n int, err error) {
	m.Data = append(m.Data, p...)
	return len(m.Data), nil
}

// NewMsg creates a new message with id
func NewMsg(typ string) Msg {
	return Msg{
		ID:        generateID(),
		Type:      typ,
		Timestamp: time.Now().UTC(),
	}
}

func generateID() string {
	b := make([]byte, 36)
	rand.Read(b)
	out := base64.RawStdEncoding.EncodeToString(b)
	return out
}
