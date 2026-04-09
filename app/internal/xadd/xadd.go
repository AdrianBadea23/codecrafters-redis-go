package xadd

import (
	"bufio"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Xadd(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	if tokens[2] == "*" {
		message := utils.AddStreamFullGen(server.Streams, tokens)
		utils.BuildBulkString(writer, message)
	} else {
		splitString := strings.Split(tokens[2], "-")

		if splitString[1] == "*" {
			message := utils.AddStreamPartialGen(server.Streams, tokens)

			if strings.Contains(message, utils.ERR) {
				writer.WriteString(message)
				writer.Flush()
			} else {
				utils.BuildBulkString(writer, message)
			}
		} else {
			message, ok := utils.ValidateStreamKey(server.Streams, tokens[1], tokens[2])
			if ok {
				message := utils.AddStream(server.Streams, tokens)
				utils.BuildBulkString(writer, message)
			} else {
				writer.WriteString(message)
				writer.Flush()
			}
		}
	}
}
