package lrange

import (
	"bufio"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Lrange(conn net.Conn, tokens []string, writer *bufio.Writer, server *server.RedisServer) {
	name := tokens[1]
	start, _ := strconv.Atoi(tokens[2])
	stop, _ := strconv.Atoi(tokens[3])
	slice := utils.GetRangeFromList(server.Lists, name, start, stop)
	if len(slice) == 0 {
		writer.WriteString(utils.EMPTY_ARRAY)
		writer.Flush()
	} else {
		message := utils.BuildArrayString(slice)
		writer.WriteString(message)
		writer.Flush()
	}
}
