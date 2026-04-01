package internal

import (
	"bufio"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {

	_, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("Error reading fron conn")
	}

	conn.Write([]byte("+PONG\r\n"))
}
