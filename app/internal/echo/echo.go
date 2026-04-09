package echo

import (
	"bufio"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Echo(conn net.Conn, tokens []string, writer *bufio.Writer) {
	utils.BuildBulkString(writer, tokens[1])
}
