package messaging

import (
	"bytes"
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
	expected := []byte(`{"data": {"a": 1}}`)
	_, err := b.Write(expected)
	if err != nil {
		t.Error(err)
	}
	m := <-out
	if !bytes.Equal(expected, m) {
		t.Errorf("want: %s, got: %s", string(expected), string(m))
	}
}
