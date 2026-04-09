package rpush

import (
	"bufio"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Rpush(writer *bufio.Writer, server *server.RedisServer, tokens []string) {
	length := utils.AddToListGrid(server.Lists, tokens, server)
	utils.BuildInteger(writer, length)
}
