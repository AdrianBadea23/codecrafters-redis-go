package internal

import (
	"bufio"
	"net"
	"strconv"
	"strings"
)

const ECHO = "ECHO"
const PING = "PING"
const SET = "SET"
const GET = "GET"

const PONG = "+PONG\r\n"
const OK = "+OK\r\n"

func HandleConnection(conn net.Conn) {

	for {
		recv := bufio.NewReader(conn)
		firstToken, _ := recv.ReadByte()

		switch firstToken {
		case '*':
			tokens := arrayParser(recv)
			writer := bufio.NewWriter(conn)

			if strings.EqualFold(tokens[0], ECHO) {
				writer.WriteString("$")
				writer.WriteString(strconv.Itoa(len(tokens[1])))
				writer.WriteString("\r\n")
				writer.WriteString(tokens[1])
				writer.WriteString("\r\n")
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], PING) {
				writer.WriteString(PONG)
				writer.Flush()
			}

		case '+':
			myString := parseString(recv)
			writer := bufio.NewWriter(conn)

			if strings.EqualFold(myString, PING) {
				writer.WriteString(PONG)
				writer.Flush()
			}
		}

	}

}
