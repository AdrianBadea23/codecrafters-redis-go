package lpush

import (
	"bufio"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Lpush(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	length := utils.PreAddToListGrid(server.Lists, tokens)
	utils.BuildInteger(writer, length)
}
