package blpop

import (
	"bufio"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func Blpop(server *server.RedisServer, writer *bufio.Writer, tokens []string) {
	key := tokens[1]
	lis, ok := server.Lists[key]

	if ok && len(lis.([]string)) > 0 {
		result := utils.LeftPop(server.Lists, key, 1)
		utils.BuildBulkString(writer, result[0])
	} else {

		timeInSeconds, _ := strconv.ParseFloat(tokens[2], 64)

		if timeInSeconds > 0 {
			server.Mu.Lock()
			channel := make(chan string, 1)
			server.Channels[key] = append(server.Channels[key], channel)
			server.Mu.Unlock()
			duration := time.Duration(timeInSeconds * float64(time.Second))
			select {
			case <-channel:
				result := utils.LeftPop(server.Lists, key, 1)
				result = append([]string{key}, result...)
				message := utils.BuildArrayString(result)
				writer.WriteString(message)
				writer.Flush()
			case <-time.After(duration):
				server.Mu.Lock()
				channelsSlice := server.Channels[key]
				channelsSlice = channelsSlice[1:]
				server.Channels[key] = append(server.Channels[key], channelsSlice...)
				server.Mu.Unlock()
				writer.WriteString(utils.NULL_ARRAY)
				writer.Flush()
			}

		} else {
			server.Mu.Lock()
			channel := make(chan string, 1)
			server.Channels[key] = append(server.Channels[key], channel)
			server.Mu.Unlock()
			<-channel
			result := utils.LeftPop(server.Lists, key, 1)
			result = append([]string{key}, result...)
			message := utils.BuildArrayString(result)
			writer.WriteString(message)
			writer.Flush()
		}

	}
}
