package internal

import (
	"bufio"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {

	for {
		command, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Error reading fron conn")
		}

		tokens := arrayParser(command)
		writer := bufio.NewWriter(conn)

		if tokens[0] == "ECHO" {
			writer.WriteString("$")
			writer.WriteString(string(len(tokens[1])))
			writer.WriteString("\r\n")
			writer.WriteString(tokens[1])
			writer.WriteString("\r\n")
			writer.Flush()
		}

		if tokens[0] == "PING" {
			writer.WriteString("+PONG\r\n")
			writer.Flush()
		}

	}

}
