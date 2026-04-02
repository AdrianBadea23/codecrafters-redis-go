package internal

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	ECHO  = "ECHO"
	PING  = "PING"
	SET   = "SET"
	GET   = "GET"
	PX    = "PX"
	EX    = "EX"
	RPUSH = "RPUSH"

	PONG             = "+PONG\r\n"
	OK               = "+OK\r\n"
	NULL_BULK_STRING = "$-1\r\n"
	RSVP_DELIMITER   = "\r\n"
	BULK_STRING      = "$"
	INTEGER          = ":"
)

func addToListGrid(listGrid map[string]any, tokens []string) int {
	_, ok := listGrid[tokens[1]]

	if !ok {
		slice := make([]string, 0)
		for i := 2; i <= len(tokens); i++ {
			slice = append(slice, tokens[i])
		}
		listGrid[tokens[1]] = slice
		return len(slice)
	} else {
		slice := listGrid[tokens[1]].([]string)
		for i := 2; i <= len(tokens); i++ {
			slice = append(slice, tokens[i])
		}
		listGrid[tokens[1]] = slice
		return len(slice)
	}
}

func HandleConnection(conn net.Conn) {

	keyValue := make(map[string]string)
	keyTimetable := make(map[string]int64)
	listGrid := make(map[string]any, 100)

	for {
		recv := bufio.NewReader(conn)
		firstToken, _ := recv.ReadByte()

		switch firstToken {
		case '*':
			tokens := arrayParser(recv)
			writer := bufio.NewWriter(conn)

			if strings.EqualFold(tokens[0], ECHO) {
				writer.WriteString(BULK_STRING)
				writer.WriteString(strconv.Itoa(len(tokens[1])))
				writer.WriteString(RSVP_DELIMITER)
				writer.WriteString(tokens[1])
				writer.WriteString(RSVP_DELIMITER)
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], RPUSH) {
				length := addToListGrid(listGrid, tokens)
				writer.WriteString(INTEGER)
				writer.WriteString(strconv.Itoa(length))
				writer.WriteString(RSVP_DELIMITER)
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], SET) {
				keyValue[tokens[1]] = tokens[2]

				if len(tokens) > 3 && strings.EqualFold(tokens[3], PX) {
					milisecondsToEnd, _ := strconv.ParseInt(tokens[4], 10, 64)
					keyTimetable[tokens[1]] = time.Now().Add(time.Duration(milisecondsToEnd) * time.Millisecond).UnixMilli()
				}

				writer.WriteString(OK)
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], GET) {
				val := keyValue[tokens[1]]
				pttl, ok := keyTimetable[tokens[1]]
				ttl := pttl - time.Now().UnixMilli()

				if ok && ttl <= 0 {
					writer.WriteString(NULL_BULK_STRING)
					writer.Flush()
				} else {
					writer.WriteString(BULK_STRING)
					writer.WriteString(strconv.Itoa(len(val)))
					writer.WriteString(RSVP_DELIMITER)
					writer.WriteString(val)
					writer.WriteString(RSVP_DELIMITER)
					writer.Flush()
				}

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
