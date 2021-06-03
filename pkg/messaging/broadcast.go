package messaging

import (
	"log"
	"net"
)

// Broadcast sends a message to all services
type Broadcast struct {
	addr *net.UDPAddr
	pc   net.PacketConn
}

// Close closes all established connections for data broadcast
func (b *Broadcast) Close() error {
	return b.pc.Close()
}

// NewBroadcast initiates all broadcast connections
func NewBroadcast() *Broadcast {
	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:8080")
	if err != nil {
		log.Println(err)
	}
	pc, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Println(err)
	}
	return &Broadcast{addr, pc}
}

// Write implements io.Writer interface
func (b *Broadcast) Write(data []byte) (int, error) {
	return b.pc.WriteTo(data, b.addr)
}
