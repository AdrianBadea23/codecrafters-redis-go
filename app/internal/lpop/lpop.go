package lpop

import (
	"bufio"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Lpop(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	slice, ok := server.Lists[tokens[1]].([]string)
	name := tokens[1]
	if !ok || len(slice) == 0 {
		writer.WriteString(utils.NULL_BULK_STRING)
		writer.Flush()
	}

	if len(tokens) == 2 {
		result := utils.LeftPop(server.Lists, name, 1)
		utils.BuildBulkString(writer, result[0])

	} else if len(tokens) == 3 {
		elements, _ := strconv.Atoi(tokens[2])
		result := utils.LeftPop(server.Lists, name, elements)
		message := utils.BuildArrayString(result)
		writer.WriteString(message)
		writer.Flush()
	}
}
