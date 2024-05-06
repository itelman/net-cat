package netcat

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func connect(c *Server, conn net.Conn) {
	defer conn.Close()
	var userName string

	scan := bufio.NewScanner(conn)
	fmt.Fprint(conn, enterName)

	if scan.Scan() {
		userName = scan.Text()
		if IsValidNick(userName) {
			fmt.Fprint(conn, c.history)
			c.addUser(conn, userName)
		} else {
			fmt.Fprint(conn, "Invalid username\n")
			connect(c, conn)
		}

	}

	defer c.DelUser(conn)
	status := Status{
		username: userName,
		conn:     conn,
		time:     time.Now(),
		text:     join,
	}
	c.History(conn, fmt.Sprintf("%v %v", status.username, status.text))
	stsCh <- status
	defer func() {
		status := Status{
			username: userName,
			conn:     conn,
			time:     time.Now(),
			text:     left,
		}
		c.History(conn, fmt.Sprintf("%v %v", status.username, status.text))
		stsCh <- status
	}()
	write(c, conn, userName, scan)
}

func write(c *Server, conn net.Conn, userName string, scan *bufio.Scanner) {
	for scan.Scan() {
		go func() {
			msg := Message{
				conn:     conn,
				username: userName,
				time:     time.Now(),
			}

			msg.text = scan.Text()
			c.History(conn, fmt.Sprintf("[%v][%v]:%v", time.Now().Format(timeFormat), msg.username, msg.text))

			msgCh <- msg
		}()
	}
}

func broadCast(c *Server) {
	for {
		select {
		case msg := <-msgCh:

			c.mu.Lock()
			for conn, name := range c.clients {

				if msg.username != name {
					fmt.Fprint(conn, "\n", fmt.Sprintf("[%v][%v]:%v", time.Now().Format(timeFormat), msg.username, msg.text), "\n")
				}

				fmt.Fprintf(conn, "[%s][%s]:", time.Now().Format(timeFormat), name)
			}
			c.mu.Unlock()

		case stat := <-stsCh:

			c.mu.Lock()
			for conn, name := range c.clients {
				if stat.username != name {
					fmt.Fprint(conn, "\n", fmt.Sprintf("%v %v", stat.username, stat.text), "\n")
				}
				fmt.Fprintf(conn, "[%s][%s]:", time.Now().Format(timeFormat), name)
			}
			c.mu.Unlock()
		}
	}
}
