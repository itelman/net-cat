package netcat

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
)

func NewServer(port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	fmt.Printf("Listening to the port: %v\n", port)
	if err != nil {
		return err
	}
	defer ln.Close()

	chat := &Server{
		mu:      sync.Mutex{},
		clients: map[net.Conn]string{},
	}

	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				chat.mu.Lock()
				defer chat.mu.Unlock()
				os.Exit(0)
			}
		}
	}()

	go broadCast(chat)
	for {
		conn, err := ln.Accept()
		if err != nil {
			os.Exit(0)
		}
		chat.mu.Lock()

		if len(chat.clients) < 10 {
			Art(conn)
			go connect(chat, conn)
		} else {
			fmt.Fprint(conn, chatFullErrText)
			conn.Close()
		}
		chat.mu.Unlock()
	}
}
