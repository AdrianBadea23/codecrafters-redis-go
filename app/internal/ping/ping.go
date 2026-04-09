package ping

import (
	"bufio"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Ping(writer *bufio.Writer) {
	writer.WriteString(utils.PONG)
	writer.Flush()
}
