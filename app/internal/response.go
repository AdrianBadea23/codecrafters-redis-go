package internal

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	ECHO   = "ECHO"
	PING   = "PING"
	SET    = "SET"
	GET    = "GET"
	PX     = "PX"
	EX     = "EX"
	RPUSH  = "RPUSH"
	LPUSH  = "LPUSH"
	LRANGE = "LRANGE"
	LLEN   = "LLEN"
	LPOP   = "LPOP"

	PONG             = "+PONG\r\n"
	OK               = "+OK\r\n"
	NULL_BULK_STRING = "$-1\r\n"
	EMPTY_ARRAY      = "*0\r\n"
	ZERO             = ":0\r\n"
	RSVP_DELIMITER   = "\r\n"
	BULK_STRING      = "$"
	INTEGER          = ":"
	ARRAY            = "*"
)

func getRangeFromList(listGrid map[string]any, sliceName string, start, stop int) []string {
	slice, ok := listGrid[sliceName].([]string)

	if !ok {
		return []string{}
	}

	length := len(slice)

	if start < 0 {
		start += length
	}

	if start < 0 {
		start = 0
	}

	if stop < 0 {
		stop += length
	}

	if stop < 0 {
		stop = 0
	}

	if start > length {
		return []string{}
	}

	if start > stop {
		return []string{}
	}

	if stop >= length {
		return slice[start:]
	}

	return slice[start : stop+1]
}

func buildArrayString(slice []string) string {
	var sb strings.Builder
	sb.WriteString(ARRAY)
	sb.WriteString(strconv.Itoa(len(slice)))
	sb.WriteString(RSVP_DELIMITER)

	for _, val := range slice {
		sb.WriteString(BULK_STRING)
		sb.WriteString(strconv.Itoa(len(val)))
		sb.WriteString(RSVP_DELIMITER)
		sb.WriteString(val)
		sb.WriteString(RSVP_DELIMITER)
	}

	return sb.String()
}

func addToListGrid(listGrid map[string]any, tokens []string) int {
	_, ok := listGrid[tokens[1]]

	if !ok {
		slice := make([]string, 0)
		for i := 2; i < len(tokens); i++ {
			slice = append(slice, tokens[i])
		}
		listGrid[tokens[1]] = slice // tokens[1] is the name given when inserting in the slice
		return len(slice)
	} else {
		slice := listGrid[tokens[1]].([]string)
		for i := 2; i < len(tokens); i++ {
			slice = append(slice, tokens[i])
		}
		listGrid[tokens[1]] = slice
		return len(slice)
	}
}

func preAddToListGrid(listGrid map[string]any, tokens []string) int {
	_, ok := listGrid[tokens[1]]

	if !ok {
		slice := make([]string, 0)
		for i := 2; i < len(tokens); i++ {
			slice = append([]string{tokens[i]}, slice...)
		}
		listGrid[tokens[1]] = slice // tokens[1] is the name given when inserting in the slice
		return len(slice)
	} else {
		slice := listGrid[tokens[1]].([]string)
		for i := 2; i < len(tokens); i++ {
			slice = append([]string{tokens[i]}, slice...)
		}
		listGrid[tokens[1]] = slice
		return len(slice)
	}
}

func leftPop(listGrid map[string]any, name string, elements int) []string {
	slice, _ := listGrid[name].([]string)

	sliceToKeep := slice[elements:]
	sliceToDiscard := slice[:elements]
	listGrid[name] = sliceToKeep

	return sliceToDiscard
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

			if strings.EqualFold(tokens[0], LRANGE) {
				name := tokens[1]
				start, _ := strconv.Atoi(tokens[2])
				stop, _ := strconv.Atoi(tokens[3])
				slice := getRangeFromList(listGrid, name, start, stop)
				if len(slice) == 0 {
					writer.WriteString(EMPTY_ARRAY)
					writer.Flush()
				} else {
					message := buildArrayString(slice)
					writer.WriteString(message)
					writer.Flush()
				}

			}

			if strings.EqualFold(tokens[0], RPUSH) {
				length := addToListGrid(listGrid, tokens)
				writer.WriteString(INTEGER)
				writer.WriteString(strconv.Itoa(length))
				writer.WriteString(RSVP_DELIMITER)
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], LPUSH) {
				length := preAddToListGrid(listGrid, tokens)
				fmt.Println(listGrid[tokens[1]])
				writer.WriteString(INTEGER)
				writer.WriteString(strconv.Itoa(length))
				writer.WriteString(RSVP_DELIMITER)
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], LLEN) {
				val, ok := listGrid[tokens[1]].([]string)
				if !ok {
					writer.WriteString(ZERO)
					writer.Flush()
				} else {
					writer.WriteString(INTEGER)
					writer.WriteString(strconv.Itoa(len(val)))
					writer.WriteString(RSVP_DELIMITER)
					writer.Flush()
				}
			}

			if strings.EqualFold(tokens[0], LPOP) {
				slice, ok := listGrid[tokens[1]].([]string)
				name := tokens[1]
				if !ok || len(slice) == 0 {
					writer.WriteString(NULL_BULK_STRING)
					writer.Flush()
				}

				if len(tokens) == 2 {
					result := leftPop(listGrid, name, 1)
					length := strconv.Itoa(len(result[0]))
					writer.WriteString(BULK_STRING)
					writer.WriteString(length)
					writer.WriteString(RSVP_DELIMITER)
					writer.WriteString(result[0])
					writer.WriteString(RSVP_DELIMITER)
					writer.Flush()
				} else if len(tokens) == 3 {
					elements, _ := strconv.Atoi(tokens[2])
					result := leftPop(listGrid, name, elements)
					message := buildArrayString(result)
					writer.WriteString(message)
					writer.Flush()
				}
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
