// +build unit

package serialization

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testEpochValues = []struct {
	Bytes  []byte
	Epoch  EpochTime
	String string
	Err    error
	Msg    string
}{
	// {
	// 	Bytes:  []byte(`1587654321`), // Seconds
	// 	Epoch:  EpochTime(time.Unix(1587654321, 0).UTC()),
	// 	String: "2020-04-23 15:05:21 +0000 UTC",
	// 	Err:    nil,
	// 	Msg:    "Epoch: Seconds",
	// },
	// {
	// 	Bytes:  []byte(`1587654321012`), // Milliseconds
	// 	Epoch:  EpochTime(time.Unix(1587654321, 12*int64(time.Millisecond)).UTC()),
	// 	String: "2020-04-23 15:05:21.012 +0000 UTC",
	// 	Err:    nil,
	// 	Msg:    "Epoch: Milliseconds",
	// },
	{
		Bytes:  []byte(`1587654321012345678`), // Nanoseconds
		Epoch:  EpochTime(time.Unix(1587654321, 12345).UTC()),
		String: "2020-04-23 15:05:21.000012345 +0000 UTC",
		Err:    nil,
		Msg:    "Epoch: Nanoseconds",
	},
	// {
	// 	Bytes:  []byte(`asdf`), // Invalid
	// 	Epoch:  EpochTime{},
	// 	String: "0001-01-01 00:00:00 +0000 UTC",
	// 	Err:    &strconv.NumError{},
	// 	Msg:    "Epoch: invalid",
	// },
}

func TestEpochUnmarshal(t *testing.T) {
	t.Parallel()

	for _, v := range testEpochValues {
		var et EpochTime
		err := et.UnmarshalJSON(v.Bytes)

		if v.Err != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, v.Epoch, et, v.Msg)
	}
}

func TestEpochMarshalJSON(t *testing.T) {
	t.Parallel()

	for _, v := range testEpochValues {
		result, err := v.Epoch.MarshalJSON()

		fmt.Print("\n\n **************************** \n")
		fmt.Printf("\n THING:  %+v - %+v\n", result, string(result))
		fmt.Print("\n **************************** \n\n")
		time.Sleep(3 * time.Second) // This line isn't necessary, but it helps to read the output

		if v.Err != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, v.Bytes, result, v.Msg)

	}
}

func TestEpochString(t *testing.T) {
	t.Parallel()

	for _, v := range testEpochValues {
		if v.Err == nil {
			res := v.Epoch.String()
			assert.Equal(t, v.String, res, v.Msg)
		}
	}
}
