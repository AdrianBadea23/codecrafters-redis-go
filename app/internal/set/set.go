package set

import (
	"bufio"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Set(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	server.Data[tokens[1]] = tokens[2]

	if len(tokens) > 3 && strings.EqualFold(tokens[3], utils.PX) {
		milisecondsToEnd, _ := strconv.ParseInt(tokens[4], 10, 64)
		server.Expires[tokens[1]] = time.Now().Add(time.Duration(milisecondsToEnd) * time.Millisecond).UnixMilli()
	}

	writer.WriteString(utils.OK)
	writer.Flush()
}
