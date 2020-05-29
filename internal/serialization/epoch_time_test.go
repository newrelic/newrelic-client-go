// +build unit

package serialization

import (
	"strconv"
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
	{
		Bytes:  []byte(`1587654321`), // Seconds
		Epoch:  EpochTime(time.Unix(1587654321, 0)),
		String: "2020-04-23 15:05:21 +0000 UTC",
		Err:    nil,
		Msg:    "Epoch: Seconds",
	},
	{
		Bytes:  []byte(`1587654321012`), // Milliseconds
		Epoch:  EpochTime(time.Unix(1587654321, 12*int64(time.Millisecond))),
		String: "2020-04-23 15:05:21.012 +0000 UTC",
		Err:    nil,
		Msg:    "Epoch: Millieconds",
	},
	{
		Bytes:  []byte(`1587654321012345678`), // Nanoseconds
		Epoch:  EpochTime(time.Unix(1587654321, 12345000)),
		String: "2020-04-23 15:05:21.000012345 +0000 UTC",
		Err:    nil,
		Msg:    "Epoch: Nanoseconds",
	},
	{
		Bytes:  []byte(`asdf`), // Invalid
		Epoch:  EpochTime(time.Unix(0, 0)),
		String: "0001-01-01 00:00:00 +0000 UTC",
		Err:    &strconv.NumError{},
		Msg:    "Epoch: invalid",
	},
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
		assert.Equal(t, v.String, time.Time(et).UTC().String(), v.Msg) // ensure to use UTC so tests work everywhere
	}
}

func TestEpochMarshalJSON(t *testing.T) {
	t.Parallel()

	for _, v := range testEpochValues {
		// MarhsalJSON never returns an error, so skip error tests
		if v.Err == nil {
			res, err := v.Epoch.MarshalJSON()

			assert.NoError(t, err)

			// Only check for millisecond resolition
			if len(v.Bytes) > 13 {
				assert.Equal(t, v.Bytes[0:13], []byte(res), v.Msg)
			} else {
				assert.Equal(t, v.Bytes, []byte(res), v.Msg)
			}
		}
	}
}
