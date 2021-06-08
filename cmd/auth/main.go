package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/rkorkosz/rocky/internal/auth"
	"github.com/rkorkosz/rocky/pkg/messaging"
)

type authMsg struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func main() {
	l := messaging.NewListener()
	l.Listen(authCallback)
}

func authCallback(m messaging.Msg) error {
	if m.Type == "/auth" {
		reply := messaging.NewMsg("/authenticated")
		m := authMsg{}
		err := auth.Authenticate(reply.Data)
		if err != nil {
			m.Error = err.Error()
		} else {
			m.Message = "ok"
		}
		return json.NewEncoder(&reply).Encode(&m)
	}
	return nil
}

func Listen(callback func(messaging.Msg) error) {
	pc, err := net.ListenPacket("udp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	buf := make([]byte, 1024)
	for {
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		if n > 0 {
			m := messaging.Msg{}
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
