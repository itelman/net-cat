package netcat

import (
	"fmt"
	"net"
	"os"
)

func IsValidNick(nickName string) bool {
	if nickName == "" || len(nickName) <= 3 {
		return false
	}
	return true
}

func Art(conn net.Conn) {
	file, _ := os.ReadFile("asset/welcome.txt")
	fmt.Fprint(conn, "Welcome to Chat\n")
	for _, ch := range file {
		fmt.Fprint(conn, string(ch))
	}
	fmt.Fprint(conn, "\n")
}

func (c *Server) addUser(conn net.Conn, username string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[conn] = username
}

func (c *Server) DelUser(conn net.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.clients, conn)
}

func (c *Server) History(conn net.Conn, text string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.history += text + "\n"
}
