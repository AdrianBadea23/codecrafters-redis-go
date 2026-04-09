package echo

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/testingutils"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func TestEcho(t *testing.T) {

	t.Run("test echo for a basic usecase", func(t *testing.T) {
		var sb strings.Builder
		buffer := bytes.Buffer{}
		writer := bufio.NewWriter(&buffer)

		tokens := []string{"", "ECHO"}

		Echo(writer, tokens)

		sb.WriteString(utils.BULK_STRING)
		sb.WriteString(strconv.Itoa(len(tokens[1])))
		sb.WriteString(utils.RESP_DELIMITER)
		sb.WriteString(tokens[1])
		sb.WriteString(utils.RESP_DELIMITER)

		testingutils.AssertEquals(t, buffer.String(), sb.String())
	})

	t.Run("test echo for a negative usecase", func(t *testing.T) {
		var sb strings.Builder
		buffer := bytes.Buffer{}
		writer := bufio.NewWriter(&buffer)

		tokens := []string{"", "ECHO"}

		Echo(writer, tokens)

		sb.WriteString(utils.BULK_STRING)
		sb.WriteString(strconv.Itoa(len(tokens[1])))
		sb.WriteString(utils.RESP_DELIMITER)
		sb.WriteString("something")
		sb.WriteString(utils.RESP_DELIMITER)

		testingutils.AssertNotEquals(t, buffer.String(), sb.String())
	})

}
