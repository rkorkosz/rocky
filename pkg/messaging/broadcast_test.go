package messaging

import (
	"bytes"
	"net"
	"testing"
)

func TestBroadcast(t *testing.T) {
	out := make(chan []byte, 1)
	defer close(out)
	go listenUDP(t, out)
	b := NewBroadcast()
	defer b.Close()
	expected := []byte("tst")
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
	buf := make([]byte, 1024)
	for {
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			t.Log(err)
			return
		}
		out <- buf[:n]
	}
}
