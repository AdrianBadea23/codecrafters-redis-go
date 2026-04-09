package parser

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/testingutils"
)

func TestParser(t *testing.T) {
	t.Run("testing the bulk string parser", func(t *testing.T) {
		buffer := bytes.Buffer{}
		buffer.WriteString("$10\r\nhellohello\r\n")
		reader := bufio.NewReader(&buffer)

		got := BulkString(reader)
		want := "hellohello"

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("testing the string parser", func(t *testing.T) {
		buffer := bytes.Buffer{}
		buffer.WriteString("+Wasup homie I'm Tony\r\n")
		reader := bufio.NewReader(&buffer)

		got := ParseString(reader)
		want := "Wasup homie I'm Tony"

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("testing the array parser", func(t *testing.T) {
		buffer := bytes.Buffer{}
		buffer.WriteString("2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
		reader := bufio.NewReader(&buffer)

		got := ArrayParser(reader)
		want := []string{"hello", "world"}

		for i := 0; i < len(got); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}

	})
}
