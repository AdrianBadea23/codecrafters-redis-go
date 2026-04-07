package internal

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func HandleConnection(conn net.Conn, server *RedisServer) {

	for {
		recv := bufio.NewReader(conn)
		firstToken, _ := recv.ReadByte()

		switch firstToken {
		case '*':
			tokens := arrayParser(recv)
			writer := bufio.NewWriter(conn)

			if strings.EqualFold(tokens[0], ECHO) {
				buildBulkString(writer, tokens[1])
				// writer.WriteString(BULK_STRING)
				// writer.WriteString(strconv.Itoa(len(tokens[1])))
				// writer.WriteString(RESP_DELIMITER)
				// writer.WriteString(tokens[1])
				// writer.WriteString(RESP_DELIMITER)
				// writer.Flush()
			}

			if strings.EqualFold(tokens[0], LRANGE) {
				name := tokens[1]
				start, _ := strconv.Atoi(tokens[2])
				stop, _ := strconv.Atoi(tokens[3])
				slice := getRangeFromList(server.Lists, name, start, stop)
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
				length := addToListGrid(server.Lists, tokens, server)
				buildInteger(writer, length)
				// writer.WriteString(INTEGER)
				// writer.WriteString(strconv.Itoa(length))
				// writer.WriteString(RESP_DELIMITER)
				// writer.Flush()
			}

			if strings.EqualFold(tokens[0], LPUSH) {
				length := preAddToListGrid(server.Lists, tokens)
				fmt.Println(server.Lists[tokens[1]])
				buildInteger(writer, length)
				// writer.WriteString(INTEGER)
				// writer.WriteString(strconv.Itoa(length))
				// writer.WriteString(RESP_DELIMITER)
				// writer.Flush()
			}

			if strings.EqualFold(tokens[0], LLEN) {
				val, ok := server.Lists[tokens[1]].([]string)
				if !ok {
					writer.WriteString(ZERO)
					writer.Flush()
				} else {
					buildInteger(writer, len(val))
					// writer.WriteString(INTEGER)
					// writer.WriteString(strconv.Itoa(len(val)))
					// writer.WriteString(RESP_DELIMITER)
					// writer.Flush()
				}
			}

			if strings.EqualFold(tokens[0], TYPE) {
				name := tokens[1]
				val := getDataType(server, name)
				writer.WriteString(val)
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], XRANGE) {
				message := rangeOverStream(server.Streams, tokens)
				writer.WriteString(message)
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], XADD) {

				if tokens[2] == "*" {
					message := addStreamFullGen(server.Streams, tokens)
					buildBulkString(writer, message)
					// writer.WriteString(BULK_STRING)
					// writer.WriteString(strconv.Itoa(len(message)))
					// writer.WriteString(RESP_DELIMITER)
					// writer.WriteString(message)
					// writer.WriteString(RESP_DELIMITER)
					// writer.Flush()
				} else {
					splitString := strings.Split(tokens[2], "-")

					if splitString[1] == "*" {
						message := addStreamPartialGen(server.Streams, tokens)

						if strings.Contains(message, ERR) {
							writer.WriteString(message)
							writer.Flush()
						} else {
							buildBulkString(writer, message)
							// writer.WriteString(BULK_STRING)
							// writer.WriteString(strconv.Itoa(len(message)))
							// writer.WriteString(RESP_DELIMITER)
							// writer.WriteString(message)
							// writer.WriteString(RESP_DELIMITER)
							// writer.Flush()
						}
					} else {
						message, ok := validateStreamKey(server.Streams, tokens[1], tokens[2])
						if ok {
							message := addStream(server.Streams, tokens)
							buildBulkString(writer, message)
							// writer.WriteString(BULK_STRING)
							// writer.WriteString(strconv.Itoa(len(message)))
							// writer.WriteString(RESP_DELIMITER)
							// writer.WriteString(message)
							// writer.WriteString(RESP_DELIMITER)
							// writer.Flush()
						} else {
							writer.WriteString(message)
							writer.Flush()
						}
					}
				}

			}

			if strings.EqualFold(tokens[0], XREAD) {
				myMap := makeMapFromTokens(tokens)
				message := queryMultipleStreams(server.Streams, myMap)
				// fmt.Println(message)
				writer.WriteString(message)
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], LPOP) {
				slice, ok := server.Lists[tokens[1]].([]string)
				name := tokens[1]
				if !ok || len(slice) == 0 {
					writer.WriteString(NULL_BULK_STRING)
					writer.Flush()
				}

				if len(tokens) == 2 {
					result := leftPop(server.Lists, name, 1)
					// length := strconv.Itoa(len(result[0]))
					buildBulkString(writer, result[0])
					// writer.WriteString(BULK_STRING)
					// writer.WriteString(length)
					// writer.WriteString(RESP_DELIMITER)
					// writer.WriteString(result[0])
					// writer.WriteString(RESP_DELIMITER)
					// writer.Flush()
				} else if len(tokens) == 3 {
					elements, _ := strconv.Atoi(tokens[2])
					result := leftPop(server.Lists, name, elements)
					message := buildArrayString(result)
					writer.WriteString(message)
					writer.Flush()
				}
			}

			if strings.EqualFold(tokens[0], BLPOP) {
				key := tokens[1]
				lis, ok := server.Lists[key]

				if ok && len(lis.([]string)) > 0 {
					result := leftPop(server.Lists, key, 1)
					// length := strconv.Itoa(len(result[0]))
					buildBulkString(writer, result[0])
					// writer.WriteString(BULK_STRING)
					// writer.WriteString(length)
					// writer.WriteString(RESP_DELIMITER)
					// writer.WriteString(result[0])
					// writer.WriteString(RESP_DELIMITER)
					// writer.Flush()
				} else {

					timeInSeconds, _ := strconv.ParseFloat(tokens[2], 64)

					if timeInSeconds > 0 {
						server.Mu.Lock()
						channel := make(chan string, 1)
						server.Channels[key] = append(server.Channels[key], channel)
						server.Mu.Unlock()
						duration := time.Duration(timeInSeconds * float64(time.Second))
						select {
						case <-channel:
							result := leftPop(server.Lists, key, 1)
							result = append([]string{key}, result...)
							message := buildArrayString(result)
							writer.WriteString(message)
							writer.Flush()
						case <-time.After(duration):
							server.Mu.Lock()
							channelsSlice := server.Channels[key]
							channelsSlice = channelsSlice[1:]
							server.Channels[key] = append(server.Channels[key], channelsSlice...)
							server.Mu.Unlock()
							writer.WriteString(NULL_ARRAY)
							writer.Flush()
						}

					} else {
						server.Mu.Lock()
						channel := make(chan string, 1)
						server.Channels[key] = append(server.Channels[key], channel)
						server.Mu.Unlock()
						<-channel
						result := leftPop(server.Lists, key, 1)
						result = append([]string{key}, result...)
						message := buildArrayString(result)
						writer.WriteString(message)
						writer.Flush()
					}

				}
			}

			if strings.EqualFold(tokens[0], SET) {
				server.Data[tokens[1]] = tokens[2]

				if len(tokens) > 3 && strings.EqualFold(tokens[3], PX) {
					milisecondsToEnd, _ := strconv.ParseInt(tokens[4], 10, 64)
					server.Expires[tokens[1]] = time.Now().Add(time.Duration(milisecondsToEnd) * time.Millisecond).UnixMilli()
				}

				writer.WriteString(OK)
				writer.Flush()
			}

			if strings.EqualFold(tokens[0], GET) {
				val := server.Data[tokens[1]]
				pttl, ok := server.Expires[tokens[1]]
				ttl := pttl - time.Now().UnixMilli()

				if ok && ttl <= 0 {
					writer.WriteString(NULL_BULK_STRING)
					writer.Flush()
				} else {
					buildBulkString(writer, val)
					// writer.WriteString(BULK_STRING)
					// writer.WriteString(strconv.Itoa(len(val)))
					// writer.WriteString(RESP_DELIMITER)
					// writer.WriteString(val)
					// writer.WriteString(RESP_DELIMITER)
					// writer.Flush()
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
