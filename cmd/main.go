package main

import (
	"fmt"
	"log"
	"os"

	netcat "netcat/internal"
)

func main() {
	var port string
	arg := os.Args
	arg = arg[1:]
	if len(arg) == 0 {
		port = "8888"
	} else if len(arg) == 1 {
		port = arg[0]
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	err := netcat.NewServer(port)
	if err != nil {
		log.Fatal(err)
	}
}
