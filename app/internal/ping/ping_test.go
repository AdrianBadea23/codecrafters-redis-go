package ping

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/testingutils"
)

func TestPing(t *testing.T) {
	t.Run("test ping function", func(t *testing.T) {
		buffer := bytes.Buffer{}
		writer := bufio.NewWriter(&buffer)

		Ping(writer)

		testingutils.AssertEquals(t, buffer.String(), "+PONG\r\n")
	})
}
