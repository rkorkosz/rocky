package messaging

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestListener(t *testing.T) {
	l := NewListener()
	out := make(chan []byte, 1)
	callback := func(m Msg) error {
		out <- m.Data
		return nil
	}
	go l.Listen(callback)
	b := NewBroadcast()
	t.Cleanup(func() {
		l.Close()
		b.Close()
		close(out)
	})
	expected := []byte(`{"a":1}`)
	msg := NewMsg("/test")
	_, err := msg.Write(expected)
	if err != nil {
		t.Error(err)
	}
	err = json.NewEncoder(b).Encode(msg)
	if err != nil {
		t.Error(err)
	}
	m := <-out
	if !bytes.Equal(expected, m) {
		t.Errorf("want: %s, got: %s", string(expected), string(m))
	}
}
