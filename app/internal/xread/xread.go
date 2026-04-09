package xread

import (
	"bufio"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Xread(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	myKeys, myIds := utils.MakeArraysFromTokens(tokens)
	message := utils.QueryMultipleStreams(server.Streams, myKeys, myIds)
	writer.WriteString(message)
	writer.Flush()
}
