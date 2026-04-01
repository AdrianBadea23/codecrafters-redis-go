package internal

import (
	"bufio"
	"io"
)

func bulkString(reader *bufio.Reader) string {
	length := 0
	b, _ := reader.ReadByte()

	if b != '$' {
		panic("problem with data type")
	}

	for {
		aux, _ := reader.ReadByte()

		if aux == '\r' {
			reader.ReadByte()
			break
		}

		length = length*10 + int(aux)
	}

	buff := make([]byte, length)
	io.ReadFull(reader, buff)
	reader.ReadByte()
	reader.ReadByte()
	return string(buff)
}

func parseString(reader *bufio.Reader) string {
	reader.ReadByte() // expect +
	fin, _ := reader.ReadString('\r')
	return fin[:len(fin)-1]
}

func arrayParser(reader *bufio.Reader) []string {
	numInteger := 0

	for {
		t, _ := reader.ReadByte()

		if t == '\r' {
			reader.ReadByte() // get rid of \n
			break
		}

		numInteger = numInteger*10 + int(t)
	}

	args := make([]string, numInteger)

	for i := range numInteger {
		aux := bulkString(reader)
		args[i] = aux
	}

	return args
}
