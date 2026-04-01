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

	number, _ := reader.ReadString(byte('\r'))
	number = number[:len(number)-1]

	length, _ := strconv.ParseInt(string(number), 10, 64)
	reader.ReadByte()
	reader.ReadByte()

	buff := make([]byte, length)
	reader.Read(buff)
	reader.ReadByte()
	reader.ReadByte()
	return string(buff)
}

func parseString(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	reader.ReadByte() // expect +
	fin, _ := reader.ReadString('\r')
	return fin[:len(fin)-1]
}

func arrayParser(input string) []string {
	reader := bufio.NewReader(strings.NewReader(input))
	reader.ReadByte() // expect *

	num, _ := reader.ReadString(byte('\r'))
	num = num[:len(num)-1]
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
