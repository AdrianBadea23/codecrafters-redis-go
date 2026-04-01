package internal

import (
	"bufio"
	"strconv"
	"strings"
)

func bulkString(reader *bufio.Reader) string {
	b, _ := reader.ReadByte()

	if b != '$' {
		panic("problem with data type")
	}

	number, _ := reader.ReadByte()

	length, _ := strconv.ParseInt(string(number), 10, 64)
	reader.ReadByte()
	reader.ReadByte()

	buff := make([]byte, length)
	reader.Read(buff)
	reader.ReadByte()
	reader.ReadByte()
	return string(buff)
}

func arrayParser(input string) []string {
	reader := bufio.NewReader(strings.NewReader(input))
	reader.ReadByte() // expect *

	num, _ := reader.ReadByte()
	numInteger, _ := strconv.Atoi(string(num))
	args := make([]string, numInteger)

	reader.ReadByte() // expect \r
	reader.ReadByte() // expect \n

	for i := range numInteger {
		aux := bulkString(reader)
		args[i] = aux
	}

	return args
}
