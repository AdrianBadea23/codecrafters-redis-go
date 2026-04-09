package llen

import (
	"bufio"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Llen(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	val, ok := server.Lists[tokens[1]].([]string)
	if !ok {
		writer.WriteString(utils.ZERO)
		writer.Flush()
	} else {
		utils.BuildInteger(writer, len(val))
	}
}
