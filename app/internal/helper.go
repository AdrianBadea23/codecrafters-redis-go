package internal

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	ECHO   = "ECHO"
	PING   = "PING"
	SET    = "SET"
	GET    = "GET"
	PX     = "PX"
	EX     = "EX"
	RPUSH  = "RPUSH"
	LPUSH  = "LPUSH"
	LRANGE = "LRANGE"
	LLEN   = "LLEN"
	LPOP   = "LPOP"
	BLPOP  = "BLPOP"
	TYPE   = "TYPE"
	XADD   = "XADD"
	XRANGE = "XRANGE"
	XREAD  = "XREAD"
	ERR    = "ERR"

	STREAM                = "+stream\r\n"
	STRING                = "+string\r\n"
	NONE                  = "+none\r\n"
	PONG                  = "+PONG\r\n"
	OK                    = "+OK\r\n"
	NULL_BULK_STRING      = "$-1\r\n"
	NULL_ARRAY            = "*-1\r\n"
	EMPTY_ARRAY           = "*0\r\n"
	ZERO                  = ":0\r\n"
	RESP_DELIMITER        = "\r\n"
	BULK_STRING           = "$"
	INTEGER               = ":"
	ARRAY                 = "*"
	XADD_ID_SMALLER_ERROR = "-ERR The ID specified in XADD is equal or smaller than the target stream top item\r\n"
	XADD_ID_GREATER_ZERO  = "-ERR The ID specified in XADD must be greater than 0-0\r\n"
)

func getRangeFromList(listGrid map[string]any, sliceName string, start, stop int) []string {
	slice, ok := listGrid[sliceName].([]string)

	if !ok {
		return []string{}
	}

	length := len(slice)

	if start < 0 {
		start += length
	}

	if start < 0 {
		start = 0
	}

	if stop < 0 {
		stop += length
	}

	if stop < 0 {
		stop = 0
	}

	if start > length {
		return []string{}
	}

	if start > stop {
		return []string{}
	}

	if stop >= length {
		return slice[start:]
	}

	return slice[start : stop+1]
}

func buildBulkString(writer *bufio.Writer, val string) {
	writer.WriteString(BULK_STRING)
	writer.WriteString(strconv.Itoa(len(val)))
	writer.WriteString(RESP_DELIMITER)
	writer.WriteString(val)
	writer.WriteString(RESP_DELIMITER)
	writer.Flush()
}

func stringBuildBulkString(sb *strings.Builder, val string) {
	sb.WriteString(BULK_STRING)
	sb.WriteString(strconv.Itoa(len(val)))
	sb.WriteString(RESP_DELIMITER)
	sb.WriteString(val)
	sb.WriteString(RESP_DELIMITER)
}

func buildInteger(writer *bufio.Writer, val int) {
	writer.WriteString(INTEGER)
	writer.WriteString(strconv.Itoa(val))
	writer.WriteString(RESP_DELIMITER)
	writer.Flush()
}

func buildArrayString(slice []string) string {
	var sb strings.Builder
	sb.WriteString(ARRAY)
	sb.WriteString(strconv.Itoa(len(slice)))
	sb.WriteString(RESP_DELIMITER)

	for _, val := range slice {
		sb.WriteString(BULK_STRING)
		sb.WriteString(strconv.Itoa(len(val)))
		sb.WriteString(RESP_DELIMITER)
		sb.WriteString(val)
		sb.WriteString(RESP_DELIMITER)
	}

	return sb.String()
}

func addToListGrid(listGrid map[string]any, tokens []string, server *RedisServer) int {
	_, ok := listGrid[tokens[1]]
	name := tokens[1]

	if !ok {
		slice := make([]string, 0)
		for i := 2; i < len(tokens); i++ {
			slice = append(slice, tokens[i])
		}
		listGrid[name] = slice // tokens[1] is the name given when inserting in the slice
		channelsSlice := server.Channels[name]
		if len(channelsSlice) > 0 {
			channelsSlice[0] <- slice[0]
			channelsSlice = channelsSlice[1:]
			server.Channels[name] = append(server.Channels[name], channelsSlice...)
		}

		return len(slice)
	} else {
		slice := listGrid[name].([]string)
		for i := 2; i < len(tokens); i++ {
			slice = append(slice, tokens[i])
		}
		listGrid[name] = slice
		return len(slice)
	}

}

func preAddToListGrid(listGrid map[string]any, tokens []string) int {
	_, ok := listGrid[tokens[1]]

	if !ok {
		slice := make([]string, 0)
		for i := 2; i < len(tokens); i++ {
			slice = append([]string{tokens[i]}, slice...)
		}
		listGrid[tokens[1]] = slice // tokens[1] is the name given when inserting in the slice
		return len(slice)
	} else {
		slice := listGrid[tokens[1]].([]string)
		for i := 2; i < len(tokens); i++ {
			slice = append([]string{tokens[i]}, slice...)
		}
		listGrid[tokens[1]] = slice
		return len(slice)
	}
}

func leftPop(listGrid map[string]any, name string, elements int) []string {
	slice, _ := listGrid[name].([]string)

	sliceToKeep := slice[elements:]
	sliceToDiscard := slice[:elements]
	listGrid[name] = sliceToKeep

	return sliceToDiscard
}

func getDataType(server *RedisServer, name string) string {
	_, ok := server.Data[name]
	_, ok1 := server.Streams[name]

	if ok {
		return STRING
	}

	if ok1 {
		return STREAM
	}

	return NONE
}

func validateStreamKey(stream map[string][]streamStruct, streamKey, id string) (string, bool) {
	splitString := strings.Split(id, "-")
	mili, _ := strconv.ParseInt(splitString[0], 10, 64)
	seq, _ := strconv.ParseInt(splitString[1], 10, 64)
	fmt.Println(mili)
	fmt.Println(seq)

	_, ok := stream[streamKey]

	if !ok {
		if mili <= 0 && seq <= 0 {
			return XADD_ID_GREATER_ZERO, false
		}
	} else {

		aux := stream[streamKey][len(stream[streamKey])-1]
		auxSplitString := strings.Split(aux.ID, "-")
		auxMili, _ := strconv.ParseInt(auxSplitString[0], 10, 64)
		auxSeq, _ := strconv.ParseInt(auxSplitString[1], 10, 64)
		fmt.Println(auxMili)
		fmt.Println(auxSeq)

		if mili <= 0 && seq <= 0 {
			return XADD_ID_GREATER_ZERO, false
		}

		if mili < auxMili {
			return XADD_ID_SMALLER_ERROR, false
		}

		if mili == auxMili && seq <= auxSeq {
			return XADD_ID_SMALLER_ERROR, false
		}
	}

	return "OK", true

}

func addStreamFullGen(stream map[string][]streamStruct, tokens []string) string {
	streamKey := tokens[1]

	_, ok := stream[streamKey]

	if !ok {
		stream[streamKey] = make([]streamStruct, 0)
	}

	id := fmt.Sprintf("%d-0", time.Now().UnixMilli())

	aux := streamStruct{
		ID:     id,
		Fields: make(map[string]any),
	}

	for i := 3; i < len(tokens); i += 2 {
		aux.Fields[tokens[i]] = tokens[i+1]
	}

	stream[id] = append(stream[id], aux)

	return id
}

func addStreamPartialGen(stream map[string][]streamStruct, tokens []string) string {
	id := tokens[1]
	newId := ""
	var auxMili int64
	var auxSeq int64
	var auxSplitString []string
	_, ok := stream[id]

	if !ok {
		stream[id] = make([]streamStruct, 0)
		auxMili = 0
		auxSeq = 0
	} else {
		auxStruct := stream[id][len(stream[id])-1]
		auxSplitString = strings.Split(auxStruct.ID, "-")

		auxMili, _ = strconv.ParseInt(auxSplitString[0], 10, 64)
		auxSeq, _ = strconv.ParseInt(auxSplitString[1], 10, 64)
	}

	splitToken := strings.Split(tokens[2], "-")
	mili, _ := strconv.ParseInt(splitToken[0], 10, 64)
	fmt.Println(mili)
	fmt.Println(ok)

	if !ok && mili == 0 {
		newId = "0-1"
	} else if mili < auxMili {
		return XADD_ID_SMALLER_ERROR
	} else if mili > auxMili {
		newId = splitToken[0] + "-0"
	} else {
		newId = auxSplitString[0] + "-" + strconv.Itoa(int(auxSeq)+1)
	}

	aux := streamStruct{
		ID:     newId,
		Fields: make(map[string]any),
	}

	for i := 3; i < len(tokens); i += 2 {
		aux.Fields[tokens[i]] = tokens[i+1]
	}

	stream[id] = append(stream[id], aux)

	return newId
}

func addStream(stream map[string][]streamStruct, tokens []string) string {
	id := tokens[1]

	_, ok := stream[id]

	if !ok {
		stream[id] = make([]streamStruct, 0)
	}

	aux := streamStruct{
		ID:     tokens[2],
		Fields: make(map[string]any),
	}

	for i := 3; i < len(tokens); i += 2 {
		aux.Fields[tokens[i]] = tokens[i+1]
	}

	stream[id] = append(stream[id], aux)

	return tokens[2]
}

func isInRange(mili, seq, minMili, minSeq, maxMili, maxSeq int64) bool {

	if mili < minMili {
		return false
	}

	if mili == minMili && seq < minSeq {
		return false
	}

	if mili > maxMili {
		return false
	}

	if mili == maxMili && seq > maxSeq {
		return false
	}

	return true
}

func helperForArray(sb *strings.Builder, value streamStruct) {
	respTwo := "*2\r\n"
	sb.WriteString(respTwo)
	sb.WriteString(BULK_STRING)
	sb.WriteString(strconv.Itoa(len(value.ID)))
	sb.WriteString(RESP_DELIMITER)
	sb.WriteString(value.ID)
	sb.WriteString(RESP_DELIMITER)
	sb.WriteString(ARRAY)
	sb.WriteString(strconv.Itoa(len(value.Fields) * 2))
	sb.WriteString(RESP_DELIMITER)
	for key, value := range value.Fields {
		sb.WriteString(BULK_STRING)
		sb.WriteString(strconv.Itoa(len(key)))
		sb.WriteString(RESP_DELIMITER)
		sb.WriteString(key)
		sb.WriteString(RESP_DELIMITER)
		sb.WriteString(BULK_STRING)
		sb.WriteString(strconv.Itoa(len(value.(string))))
		sb.WriteString(RESP_DELIMITER)
		sb.WriteString(value.(string))
		sb.WriteString(RESP_DELIMITER)
	}
}

func rangeOverStream(stream map[string][]streamStruct, tokens []string) string {
	slice := stream[tokens[1]]
	numberOfEles := 0

	var minMili int64
	var maxMili int64
	var minSeq int64
	var maxSeq int64

	var sb strings.Builder
	var fsb strings.Builder

	if strings.Contains(tokens[2], "-") {

		if tokens[2] == "-" {
			minMili = 0
			minSeq = 1
		} else {
			spl := strings.Split(tokens[2], "-")
			minMili, _ = strconv.ParseInt(spl[0], 10, 64)
			minSeq, _ = strconv.ParseInt(spl[1], 10, 64)
		}

	} else {
		minMili, _ = strconv.ParseInt(tokens[2], 10, 64)
		minSeq = 0
	}

	if strings.Contains(tokens[3], "-") {
		spl := strings.Split(tokens[3], "-")
		maxMili, _ = strconv.ParseInt(spl[0], 10, 64)
		maxSeq, _ = strconv.ParseInt(spl[1], 10, 64)
	} else {

		if tokens[3] == "+" {
			maxMili = math.MaxInt64
			maxSeq = math.MaxInt64
		} else {
			maxMili, _ = strconv.ParseInt(tokens[3], 10, 64)
			maxSeq = math.MaxInt64
		}

	}

	for _, value := range slice {

		split := strings.Split(value.ID, "-")
		mili, _ := strconv.ParseInt(split[0], 10, 64)
		seq, _ := strconv.ParseInt(split[1], 10, 64)

		if isInRange(mili, seq, minMili, minSeq, maxMili, maxSeq) {
			numberOfEles++
			helperForArray(&sb, value)
		}
	}

	fsb.WriteString(ARRAY)
	fsb.WriteString(strconv.Itoa(numberOfEles))
	fsb.WriteString(RESP_DELIMITER)
	fsb.WriteString(sb.String())

	return fsb.String()

}

func isInRangeXread(milis, seq, auxMilis, auxSeq int64) bool {
	if auxMilis < milis {
		return false
	}

	if auxMilis == milis && auxSeq <= seq {
		return false
	}

	return true
}

func splitAndReturnInt(streamId string) (int64, int64) {
	splits := strings.Split(streamId, "-")
	milis, _ := strconv.ParseInt(splits[0], 10, 64)
	seq, _ := strconv.ParseInt(splits[1], 10, 64)

	return milis, seq
}

func preBuildString(sb *strings.Builder, streamKey string) {
	sb.WriteString(ARRAY)
	sb.WriteString("1")
	sb.WriteString(RESP_DELIMITER)
	sb.WriteString(ARRAY)
	sb.WriteString("2")
	sb.WriteString(RESP_DELIMITER)
	stringBuildBulkString(sb, streamKey)
}

func queryStream(stream map[string][]streamStruct, streamKey, streamId string) string {
	slice := stream[streamKey]
	var sb strings.Builder
	var fsb strings.Builder
	numOfEles := 0

	preBuildString(&fsb, streamKey)

	milis, seq := splitAndReturnInt(streamId)
	// fmt.Println(slice, streamKey, streamId)

	for _, val := range slice {
		auxMilis, auxSeq := splitAndReturnInt(val.ID)
		if isInRangeXread(milis, seq, auxMilis, auxSeq) {
			numOfEles++
			helperForArray(&sb, val)
		}
	}
	fsb.WriteString(ARRAY)
	fsb.WriteString(strconv.Itoa(numOfEles))
	fsb.WriteString(RESP_DELIMITER)
	fsb.WriteString(sb.String())

	return fsb.String()
}

func queryMultipleStreams(stream map[string][]streamStruct, keyIdPairs map[string]string) string {
	var sb strings.Builder
	slice := make([]string, 0)

	for key, val := range keyIdPairs {
		message := queryStream(stream, key, val)
		slice = append(slice, message)
	}

	sb.WriteString(ARRAY)
	sb.WriteString(strconv.Itoa(len(slice)))
	sb.WriteString(RESP_DELIMITER)

	for _, val := range slice {
		sb.WriteString(val)
	}

	sb.WriteString(RESP_DELIMITER)

	return sb.String()

}

func makeMapFromTokens(token []string) map[string]string {
	auxSlice := token[2:]
	length := len(auxSlice) / 2
	myMap := make(map[string]string)

	for i := 0; i <= length; i++ {
		myMap[auxSlice[i]] = auxSlice[i+length]
	}

	return myMap
}
