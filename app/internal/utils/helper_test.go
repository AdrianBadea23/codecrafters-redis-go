package utils

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/testingutils"
)

func TestGetRangeFromList(t *testing.T) {
	t.Run("test list not found", func(t *testing.T) {
		someMap := make(map[string]any)
		got := GetRangeFromList(someMap, "something", -1, 1)
		want := []string{}

		for i := 0; i < len(want); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}

	})

	t.Run("test start > length", func(t *testing.T) {
		someMap := make(map[string]any)
		someMap["something"] = []string{"1", "2"}
		got := GetRangeFromList(someMap, "something", 10, 1)
		want := []string{}

		for i := 0; i < len(want); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}

	})

	t.Run("test start > stop", func(t *testing.T) {
		someMap := make(map[string]any)
		someMap["something"] = []string{"1", "2"}
		got := GetRangeFromList(someMap, "something", 1, 0)
		want := []string{}

		for i := 0; i < len(want); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}

	})

	t.Run("test start < length", func(t *testing.T) {
		someMap := make(map[string]any)
		someMap["something"] = []string{"1", "2"}
		got := GetRangeFromList(someMap, "something", -3, 1)
		want := []string{"1", "2"}

		for i := 0; i < len(want); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}

	})

	t.Run("test start fits", func(t *testing.T) {
		someMap := make(map[string]any)
		someMap["something"] = []string{"1", "2"}
		got := GetRangeFromList(someMap, "something", -2, 1)
		want := []string{"1", "2"}

		for i := 0; i < len(want); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}

	})

	t.Run("test stop < length", func(t *testing.T) {
		someMap := make(map[string]any)
		someMap["something"] = []string{"1", "2"}
		got := GetRangeFromList(someMap, "something", 0, -3)
		want := []string{"1"}

		for i := 0; i < len(want); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}

	})

	t.Run("test stop fits", func(t *testing.T) {
		someMap := make(map[string]any)
		someMap["something"] = []string{"1", "2"}
		got := GetRangeFromList(someMap, "something", -2, -2)
		want := []string{"1"}

		for i := 0; i < len(want); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}

	})

	t.Run("test stop fits", func(t *testing.T) {
		someMap := make(map[string]any)
		someMap["something"] = []string{"1", "2", "3"}
		got := GetRangeFromList(someMap, "something", 0, 2)
		want := []string{"1", "2"}

		for i := 0; i < len(want); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}

	})
}

func TestBuildBulkString(t *testing.T) {
	t.Run("test build bulk string", func(t *testing.T) {
		buffer := bytes.Buffer{}
		got := bufio.NewWriter(&buffer)

		val := "something"
		want := "$9\r\nsomething\r\n"

		BuildBulkString(got, val)

		testingutils.AssertEquals(t, buffer.String(), want)

	})
}

func TestStringBuildBulkString(t *testing.T) {
	t.Run("test build bulk string", func(t *testing.T) {
		var sb strings.Builder

		val := "something"
		want := "$9\r\nsomething\r\n"

		StringBuildBulkString(&sb, val)

		testingutils.AssertEquals(t, sb.String(), want)

	})
}

func TestBuildInteger(t *testing.T) {
	t.Run("test build integer", func(t *testing.T) {
		buffer := bytes.Buffer{}
		got := bufio.NewWriter(&buffer)

		val := 2
		want := ":2\r\n"

		BuildInteger(got, val)

		testingutils.AssertEquals(t, buffer.String(), want)

	})
}

func TestBuildArrayString(t *testing.T) {
	t.Run("test build bulk array string", func(t *testing.T) {
		slice := []string{"some", "thing"}
		got := BuildArrayString(slice)

		want := "*2\r\n$4\r\nsome\r\n$5\r\nthing\r\n"

		testingutils.AssertEquals(t, got, want)
	})
}

func TestAddToListGrid(t *testing.T) {
	t.Run("Test addition to list grid when listGrid is empty", func(t *testing.T) {
		listGrid := make(map[string]any)
		tokens := []string{"", "something", "ceva", "altceva"}
		server := server.New()

		got := AddToListGrid(listGrid, tokens, server)
		want := 2

		testingutils.AssertEquals(t, got, want)

	})

	t.Run("Test addition to list grid when listGrid is empty and channels are not empty", func(t *testing.T) {
		var wg sync.WaitGroup

		listGrid := make(map[string]any)
		tokens := []string{"", "something", "ceva", "altceva"}
		server := server.New()
		channel := make(chan string)
		server.Channels["something"] = []chan string{channel}

		wg.Go(func() {
			got := <-channel
			want := "ceva"

			testingutils.AssertEquals(t, got, want)
		})

		got := AddToListGrid(listGrid, tokens, server)
		want := 2

		testingutils.AssertEquals(t, got, want)

	})

	t.Run("Test addition to list when listgrid is not empty", func(t *testing.T) {
		listGrid := make(map[string]any)
		listGrid["something"] = []string{"", "something", "ceva", "altceva"}
		tokens := []string{"", "something", "ceva", "altceva"}
		server := server.New()

		got := AddToListGrid(listGrid, tokens, server)
		want := 6

		testingutils.AssertEquals(t, got, want)

	})
}

func TestPreAddToListGrid(t *testing.T) {
	t.Run("Test addition when list is not empty", func(t *testing.T) {
		listGrid := make(map[string]any)
		listGrid["something"] = []string{"", "something", "ceva", "altceva"}
		tokens := []string{"", "something", "ceva", "altceva"}

		got := PreAddToListGrid(listGrid, tokens)
		want := 6

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Test addition when list is empty", func(t *testing.T) {
		listGrid := make(map[string]any)
		tokens := []string{"", "something", "ceva", "altceva"}

		got := PreAddToListGrid(listGrid, tokens)
		want := 2

		testingutils.AssertEquals(t, got, want)
	})
}

func TestLeftPop(t *testing.T) {
	t.Run("Test left pop", func(t *testing.T) {
		listGrid := make(map[string]any)
		listGrid["something"] = []string{"", "something", "ceva", "altceva"}
		got := LeftPop(listGrid, "something", 2)
		want := []string{"", "something"}

		for i := 0; i < len(want); i++ {
			testingutils.AssertEquals(t, got[i], want[i])
		}
	})
}

func TestGetDataType(t *testing.T) {
	t.Run("Test Get Data Type for string", func(t *testing.T) {
		server := server.New()
		server.Data["something"] = "..."

		got := GetDataType(server, "something")
		want := STRING

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Test Get Data Type for stream", func(t *testing.T) {
		strct := server.StreamStruct{
			ID: "something",
		}

		q := server.New()

		q.Streams["something"] = []server.StreamStruct{strct}

		got := GetDataType(q, "something")
		want := STREAM

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Test Get Data Type for none", func(t *testing.T) {
		q := server.New()

		got := GetDataType(q, "something")
		want := NONE

		testingutils.AssertEquals(t, got, want)
	})
}

func TestValidateStreamKey(t *testing.T) {
	t.Run("Validating stream key is ok", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		streamKey := "0-1"
		id := "0-1"

		got, gotBool := ValidateStreamKey(stream, streamKey, id)
		want := STRINGOK
		wantBool := true

		testingutils.AssertEquals(t, gotBool, wantBool)
		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Validating stream key is not ok", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		streamKey := "0-0"
		id := "0-0"

		got, gotBool := ValidateStreamKey(stream, streamKey, id)
		want := XADD_ID_GREATER_ZERO
		wantBool := false

		testingutils.AssertEquals(t, gotBool, wantBool)
		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Validating stream key is not ok in nonempty stream", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		streamKey := "0-1"
		id := "0-0"

		stream[streamKey] = append(stream[streamKey], server.StreamStruct{ID: "0-1"})

		got, gotBool := ValidateStreamKey(stream, streamKey, id)
		want := XADD_ID_GREATER_ZERO
		wantBool := false

		testingutils.AssertEquals(t, gotBool, wantBool)
		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Validating stream key is not ok in nonempty stream", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		streamKey := "0-1"
		id := "0-1"

		stream[streamKey] = append(stream[streamKey], server.StreamStruct{ID: "0-1"})

		got, gotBool := ValidateStreamKey(stream, streamKey, id)
		want := XADD_ID_SMALLER_ERROR
		wantBool := false

		testingutils.AssertEquals(t, gotBool, wantBool)
		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Validating stream key is not ok in nonempty stream", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		streamKey := "0-1"
		id := "0-1"

		stream[streamKey] = append(stream[streamKey], server.StreamStruct{ID: "1-1"})

		got, gotBool := ValidateStreamKey(stream, streamKey, id)
		want := XADD_ID_SMALLER_ERROR
		wantBool := false

		testingutils.AssertEquals(t, gotBool, wantBool)
		testingutils.AssertEquals(t, got, want)
	})

}

func TestAddStreamFullGen(t *testing.T) {
	t.Run("Test Add Stream function", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		tokens := []string{"", "something", "", "1", "2", "3", "2"}

		AddStreamFullGen(stream, tokens)
	})
}

func TestAddStreamPartialGen(t *testing.T) {
	t.Run("Test partial add with empty stream", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		tokens := []string{"", "something", "0-4"}

		got := AddStreamPartialGen(stream, tokens)
		want := "0-1"

		testingutils.AssertEquals(t, got, want)

	})

	t.Run("Test partial add with non empty stream", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		stream["something"] = []server.StreamStruct{{ID: "1-1"}}
		tokens := []string{"", "something", "1-1"}

		got := AddStreamPartialGen(stream, tokens)
		want := "1-2"

		testingutils.AssertEquals(t, got, want)

	})

	t.Run("Test partial add with non empty stream", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		stream["something"] = []server.StreamStruct{{ID: "1-1"}}
		tokens := []string{"", "something", "2-0"}

		got := AddStreamPartialGen(stream, tokens)
		want := "2-0"

		testingutils.AssertEquals(t, got, want)

	})

	t.Run("Test partial add with non empty stream", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		stream["something"] = []server.StreamStruct{{ID: "1-1"}}
		tokens := []string{"", "something", "0-1"}

		got := AddStreamPartialGen(stream, tokens)
		want := XADD_ID_SMALLER_ERROR

		testingutils.AssertEquals(t, got, want)

	})
}

func TestAddStream(t *testing.T) {
	t.Run("Add to stream", func(t *testing.T) {
		stream := make(map[string][]server.StreamStruct)
		tokens := []string{"", "", "1-1"}

		got := AddStream(stream, tokens)
		want := "1-1"
		testingutils.AssertEquals(t, got, want)
	})
}

func TestIsInRange(t *testing.T) {
	t.Run("Mili < MinMili", func(t *testing.T) {
		mili := 0
		seq := 1
		minMili := 1
		minSeq := 0
		maxMili := 1
		maxSeq := 0

		got := IsInRange(int64(mili), int64(seq), int64(minMili), int64(minSeq), int64(maxMili), int64(maxSeq))
		want := false

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Mili == MinMili, Seq < MinSeq", func(t *testing.T) {
		mili := 1
		seq := 1
		minMili := 1
		minSeq := 2
		maxMili := 1
		maxSeq := 0

		got := IsInRange(int64(mili), int64(seq), int64(minMili), int64(minSeq), int64(maxMili), int64(maxSeq))
		want := false

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Mili > MaxMili", func(t *testing.T) {
		mili := 3
		seq := 1
		minMili := 1
		minSeq := 2
		maxMili := 1
		maxSeq := 6

		got := IsInRange(int64(mili), int64(seq), int64(minMili), int64(minSeq), int64(maxMili), int64(maxSeq))
		want := false

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Mili > MaxMili, Seq > MaxSeq", func(t *testing.T) {
		mili := 1
		seq := 7
		minMili := 1
		minSeq := 2
		maxMili := 1
		maxSeq := 6

		got := IsInRange(int64(mili), int64(seq), int64(minMili), int64(minSeq), int64(maxMili), int64(maxSeq))
		want := false

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("Good range", func(t *testing.T) {
		mili := 1
		seq := 3
		minMili := 1
		minSeq := 2
		maxMili := 1
		maxSeq := 6

		got := IsInRange(int64(mili), int64(seq), int64(minMili), int64(minSeq), int64(maxMili), int64(maxSeq))
		want := true

		testingutils.AssertEquals(t, got, want)
	})
}

func TestHelperForArray(t *testing.T) {
	t.Run("Test helper for array", func(t *testing.T) {
		var sb strings.Builder
		mapy := make(map[string]any)
		mapy["something"] = "123"
		value := server.StreamStruct{ID: "0-1", Fields: mapy}

		HelperForArray(&sb, value)
		want := "*2\r\n$3\r\n0-1\r\n*2\r\n$9\r\nsomething\r\n$3\r\n123\r\n"

		testingutils.AssertEquals(t, sb.String(), want)
	})
}

func TestIsInRangeXread(t *testing.T) {
	t.Run("AuxMilis < Milis", func(t *testing.T) {
		milis := 1
		seq := 1

		auxMilis := 0
		auxSeq := 0

		got := isInRangeXread(int64(milis), int64(seq), int64(auxMilis), int64(auxSeq))
		want := false

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("AuxMilis == Milis, auxSeq <= seq", func(t *testing.T) {
		milis := 1
		seq := 1

		auxMilis := 1
		auxSeq := 0

		got := isInRangeXread(int64(milis), int64(seq), int64(auxMilis), int64(auxSeq))
		want := false

		testingutils.AssertEquals(t, got, want)
	})

	t.Run("True statement", func(t *testing.T) {
		milis := 1
		seq := 1

		auxMilis := 1
		auxSeq := 2

		got := isInRangeXread(int64(milis), int64(seq), int64(auxMilis), int64(auxSeq))
		want := true

		testingutils.AssertEquals(t, got, want)
	})
}

func TestSplitAndReturnInt(t *testing.T) {
	t.Run("Split and return int64s", func(t *testing.T) {
		gotMilis, gotSeq := splitAndReturnInt("0-1")
		wantMilis := int64(0)
		wantSeq := int64(1)

		testingutils.AssertEquals(t, gotMilis, wantMilis)
		testingutils.AssertEquals(t, gotSeq, wantSeq)
	})
}

func TestPreBuildString(t *testing.T) {
	t.Run("Prebuilt string test", func(t *testing.T) {
		var sb strings.Builder
		streamkey := "0-1"

		preBuildString(&sb, streamkey)
		want := "*2\r\n" + BULK_STRING + strconv.Itoa(len(streamkey)) + RESP_DELIMITER + streamkey + RESP_DELIMITER

		testingutils.AssertEquals(t, sb.String(), want)
	})
}

func TestMakeArraysFromTokens(t *testing.T) {
	t.Run("Check value of SUT", func(t *testing.T) {
		tokens := []string{"", "", "something", "something-else", "0-1", "0-2"}

		gotKeys, gotIds := MakeArraysFromTokens(tokens)
		wantKeys := []string{"something", "something-else"}
		wantIds := []string{"0-1", "0-2"}

		for i := 0; i < len(wantKeys); i++ {
			testingutils.AssertEquals(t, gotKeys[i], wantKeys[i])
		}

		for i := 0; i < len(wantIds); i++ {
			testingutils.AssertEquals(t, gotIds[i], wantIds[i])
		}
	})
}
