package echo

import (
	"bufio"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Echo(writer *bufio.Writer, tokens []string) {
	utils.BuildBulkString(writer, tokens[1])
}
