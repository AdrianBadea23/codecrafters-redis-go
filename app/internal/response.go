package internal

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

func HandleConnection(conn net.Conn) {

	for {
		recv, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Error reading fron conn")
		}

		var sb strings.Builder
		sb.WriteString(recv)

		if recv[0] == '*' {
			specLen, err := strconv.Atoi(recv[1 : len(recv)-2])

			if err != nil {
				log.Println("error with conversion")
			}

			for range specLen {
				chunk, _ := bufio.NewReader(conn).ReadString('\n')
				sb.WriteString(chunk)
			}
		}

		command := sb.String()
		reader := bufio.NewReader(strings.NewReader(command))

		firstToken, _ := reader.ReadByte()

		switch firstToken {
		case '*':
			tokens := arrayParser(command)
			writer := bufio.NewWriter(conn)

			if strings.EqualFold(tokens[0], "ECHO") {
				writer.WriteString("$")
				writer.WriteString(string(len(tokens[1])))
				writer.WriteString("\r\n")
				writer.WriteString(tokens[1])
				writer.WriteString("\r\n")
				writer.Flush()
			}

		case '+':
			myString := parseString(command)
			writer := bufio.NewWriter(conn)

			if strings.EqualFold(myString, "PING") {
				writer.WriteString("+PONG\r\n")
				writer.Flush()
			}
		}

	}

}
