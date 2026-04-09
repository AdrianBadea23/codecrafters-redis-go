package rpush

import (
	"bufio"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Rpush(conn net.Conn, server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	length := utils.AddToListGrid(server.Lists, tokens, server)
	utils.BuildInteger(writer, length)
}
