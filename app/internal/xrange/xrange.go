package xrange

import (
	"bufio"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Xrange(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	message := utils.RangeOverStream(server.Streams, tokens)
	writer.WriteString(message)
	writer.Flush()
}
