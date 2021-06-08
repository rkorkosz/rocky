package messaging

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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
		ID:        uuid.New().String(),
		Type:      typ,
		Timestamp: time.Now().UTC(),
	}
}
