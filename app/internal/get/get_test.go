package get

import (
	"bufio"
	"bytes"
	"math"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/testingutils"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func TestGet(t *testing.T) {
	t.Run("testing get command with expected success", func(t *testing.T) {
		server := server.New()
		got := bytes.Buffer{}
		writer := bufio.NewWriter(&got)
		tokens := []string{"", "some_key"}

		server.Data["some_key"] = "something"
		server.Expires["some_key"] = math.MaxInt64

		Get(server, writer, tokens)

		want := bytes.Buffer{}
		wantWriter := bufio.NewWriter(&want)
		utils.BuildBulkString(wantWriter, "something")

		testingutils.AssertEquals(t, got.String(), want.String())

	})

	t.Run("testing get command with expected fail", func(t *testing.T) {
		server := server.New()
		got := bytes.Buffer{}
		writer := bufio.NewWriter(&got)
		tokens := []string{"", "some_key"}

		server.Data["some_key"] = "something"
		server.Expires["some_key"] = -1

		Get(server, writer, tokens)

		testingutils.AssertEquals(t, got.String(), utils.NULL_BULK_STRING)

	})
}
