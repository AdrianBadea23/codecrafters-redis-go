package response

import (
	"bufio"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/blpop"
	"github.com/codecrafters-io/redis-starter-go/app/internal/echo"
	"github.com/codecrafters-io/redis-starter-go/app/internal/get"
	"github.com/codecrafters-io/redis-starter-go/app/internal/llen"
	"github.com/codecrafters-io/redis-starter-go/app/internal/lpop"
	"github.com/codecrafters-io/redis-starter-go/app/internal/lpush"
	"github.com/codecrafters-io/redis-starter-go/app/internal/lrange"
	"github.com/codecrafters-io/redis-starter-go/app/internal/parser"
	"github.com/codecrafters-io/redis-starter-go/app/internal/ping"
	"github.com/codecrafters-io/redis-starter-go/app/internal/rpush"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/set"
	typee "github.com/codecrafters-io/redis-starter-go/app/internal/type"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
	"github.com/codecrafters-io/redis-starter-go/app/internal/xadd"
	"github.com/codecrafters-io/redis-starter-go/app/internal/xrange"
	"github.com/codecrafters-io/redis-starter-go/app/internal/xread"
)

func HandleConnection(conn net.Conn, server *server.RedisServer) {

	for {
		recv := bufio.NewReader(conn)
		firstToken, _ := recv.ReadByte()

		if firstToken == '+' {
			myString := parser.ParseString(recv)
			writer := bufio.NewWriter(conn)

			if strings.EqualFold(myString, utils.PING) {
				ping.Ping(writer)
			}
			continue
		}
		tokens := parser.ArrayParser(recv)
		writer := bufio.NewWriter(conn)

		switch tokens[0] {
		case utils.ECHO:

			echo.Echo(writer, tokens)

		case utils.LRANGE:

			lrange.Lrange(writer, server, tokens)

		case utils.RPUSH:

			rpush.Rpush(writer, server, tokens)

		case utils.LPUSH:

			lpush.Lpush(server, writer, tokens)

		case utils.LLEN:

			llen.Llen(server, writer, tokens)

		case utils.TYPE:

			typee.Typee(server, writer, tokens)

		case utils.XRANGE:

			xrange.Xrange(server, writer, tokens)

		case utils.XADD:

			xadd.Xadd(server, writer, tokens)

		case utils.XREAD:

			xread.Xread(server, writer, tokens)

		case utils.LPOP:

			lpop.Lpop(server, writer, tokens)

		case utils.BLPOP:

			blpop.Blpop(server, writer, tokens)

		case utils.SET:

			set.Set(server, writer, tokens)

		case utils.GET:

			get.Get(server, writer, tokens)

		case utils.PING:

			ping.Ping(writer)

		}

	}

}
