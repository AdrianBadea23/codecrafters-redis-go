package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
)

// keyValue := make(map[string]string)
// keyTimetable := make(map[string]int64)
// listGrid := make(map[string]any, 100)

func main() {

	server := internal.New()

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go internal.HandleConnection(conn, server)
	}

}
