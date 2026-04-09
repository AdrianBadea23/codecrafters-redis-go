package typee

import (
	"bufio"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Typee(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	name := tokens[1]
	val := utils.GetDataType(server, name)
	writer.WriteString(val)
	writer.Flush()
}
