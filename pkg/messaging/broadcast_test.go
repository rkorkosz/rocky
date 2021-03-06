package messaging

import (
	"bytes"
	"net"
	"testing"
)

func TestBroadcast(t *testing.T) {
	out := make(chan []byte, 1)
	go listenUDP(t, out)
	b := NewBroadcast()
	t.Cleanup(func() {
		close(out)
		b.Close()
	})
	expected := []byte("test")
	_, err := b.Write(expected)
	if err != nil {
		t.Error(err)
	}
	actual := <-out
	if !bytes.Equal(expected, actual) {
		t.Errorf("want: %s, got: %s", string(expected), string(actual))
	}
}

func listenUDP(t *testing.T, out chan []byte) {
	t.Helper()
	pc, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer pc.Close()
	buf := make([]byte, 8)
	for {
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			t.Log(err)
			return
		}
		if n > 0 {
			out <- buf[:n]
			return
		}
	}
}
