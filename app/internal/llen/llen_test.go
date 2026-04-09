package llen

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/testingutils"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func TestLlen(t *testing.T) {
	t.Run("test len with nonzero return value", func(t *testing.T) {
		server := server.New()
		name := "something"
		slice := []string{"something", "something2"}
		server.Lists[name] = slice
		tokens := []string{"", name}

		got := bytes.Buffer{}
		writer := bufio.NewWriter(&got)
		Llen(server, writer, tokens)

		testingutils.AssertEquals(t, got.String(), ":2\r\n")

	})

	t.Run("test len with zero return value", func(t *testing.T) {
		server := server.New()
		name := "something"
		slice := []string{"something", "something2"}
		server.Lists[name] = slice
		tokens := []string{"", "name"}

		got := bytes.Buffer{}
		writer := bufio.NewWriter(&got)
		Llen(server, writer, tokens)

		testingutils.AssertEquals(t, got.String(), utils.ZERO)

	})
}
