package get

import (
	"bufio"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Get(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	val := server.Data[tokens[1]]
	pttl, ok := server.Expires[tokens[1]]
	ttl := pttl - time.Now().UnixMilli()

	if ok && ttl <= 0 {
		writer.WriteString(utils.NULL_BULK_STRING)
		writer.Flush()
	} else {
		utils.BuildBulkString(writer, val)
	}

}
