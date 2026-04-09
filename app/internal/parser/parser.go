package parser

import (
	"bufio"
	"io"
)

func BulkString(reader *bufio.Reader) string {
	length := 0
	b, _ := reader.ReadByte()

	if b != '$' {
		panic("problem with data type")
	}

	for {
		t, _ := reader.ReadByte()

		if t == '\r' {
			reader.ReadByte()
			break
		}

		length = length*10 + int(t-'0')
	}

	buff := make([]byte, length)
	io.ReadFull(reader, buff)
	reader.ReadByte()
	reader.ReadByte()
	return string(buff)
}

func ParseString(reader *bufio.Reader) string {
	reader.ReadByte() // expect +
	fin, _ := reader.ReadString('\r')
	return fin[:len(fin)-1]
}

func ArrayParser(reader *bufio.Reader) []string {
	numInteger := 0

	for {
		t, _ := reader.ReadByte()

		if t == '\r' {
			reader.ReadByte() // get rid of \n
			break
		}

		numInteger = numInteger*10 + int(t-'0')
	}

	args := make([]string, numInteger)

	for i := range numInteger {
		aux := BulkString(reader)
		args[i] = aux
	}

	return args
}
