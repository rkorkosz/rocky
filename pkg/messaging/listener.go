package messaging

import (
	"encoding/json"
	"log"
	"net"
)

type ListenUDP struct {
	pc net.PacketConn
}

func NewListener() *ListenUDP {
	pc, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	return &ListenUDP{pc}
}

func (l *ListenUDP) Close() error {
	return l.pc.Close()
}

func (l *ListenUDP) Listen(callback func(Msg) error) {
	buf := make([]byte, 1024)
	for {
		n, _, err := l.pc.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		if n > 0 {
			m := Msg{}
			err := json.Unmarshal(buf[:n], &m)
			if err != nil {
				log.Println(err)
				continue
			}
			err = callback(m)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
