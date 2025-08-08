//go:build unit
// +build unit

package serialization

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testEpochValuesSuccess = []struct {
	Name   string
	Bytes  []byte
	Epoch  EpochTime
	String string
	Unix   int64
	Err    error
	Msg    string
}{
	{
		Name:   "success seconds",
		Bytes:  []byte(`1587654321`), // Seconds
		Epoch:  EpochTime(time.Unix(1587654321, 0).UTC()),
		String: "2020-04-23 15:05:21 +0000 UTC",
		Unix:   1587654321,
		Err:    nil,
		Msg:    "Epoch: Seconds",
	},
	{
		Name:   "success milliseconds",
		Bytes:  []byte(`1587654321012`), // Milliseconds
		Epoch:  EpochTime(time.Unix(1587654321, 12*int64(time.Millisecond)).UTC()),
		String: "2020-04-23 15:05:21.012 +0000 UTC",
		Unix:   1587654321,
		Err:    nil,
		Msg:    "Epoch: Millieconds",
	},
	{
		Name:   "success nanoseconds",
		Bytes:  []byte(`1587654321000012345`), // Nanoseconds
		Epoch:  EpochTime(time.Unix(1587654321, 12345).UTC()),
		String: "2020-04-23 15:05:21.000012345 +0000 UTC",
		Unix:   1587654321,
		Err:    nil,
		Msg:    "Epoch: Nanoseconds",
	},
	{
		Name:   "success nanoseconds",
		Bytes:  []byte(`1587654321000012345`), // Nanoseconds
		Epoch:  EpochTime(time.Unix(1587654321, 12345).UTC()),
		String: "2020-04-23 15:05:21.000012345 +0000 UTC",
		Unix:   1587654321,
		Err:    nil,
		Msg:    "Epoch: Nanoseconds",
	},
	{
		Name:   "success empty time object",
		Bytes:  []byte(emptyTimeCase), // Invalid
		Epoch:  EpochTime{},
		String: "0001-01-01 00:00:00 +0000 UTC",
		Unix:   -62135596800,
		Err:    nil,
		Msg:    "Epoch: invalid",
	},
}

func TestEpochUnmarshal(t *testing.T) {
	t.Parallel()

	for _, v := range testEpochValuesSuccess {
		t.Run(v.Name, func(t *testing.T) {
			var et EpochTime
			err := et.UnmarshalJSON(v.Bytes)

			if v.Err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, v.Epoch, et, v.Msg)
		})
	}

	t.Run("failed wrong input", func(t *testing.T) {
		var et EpochTime
		err := et.UnmarshalJSON([]byte(`asdf`))

		assert.Error(t, err)
		assert.Equal(t, EpochTime{}, et)
	})
}

func TestEpochMarshalJSON(t *testing.T) {
	t.Parallel()

	for _, v := range testEpochValuesSuccess {
		t.Run(v.Name, func(t *testing.T) {
			bytes, err := v.Epoch.MarshalJSON()

			if v.Err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, v.Bytes, bytes, v.Msg)
		})
	}
}

func TestEpochString(t *testing.T) {
	t.Parallel()

	for _, v := range testEpochValuesSuccess {
		if v.Err == nil {
			res := v.Epoch.String()
			assert.Equal(t, v.String, res, v.Msg)
		}
	}
}

func TestEpochUnix(t *testing.T) {
	t.Parallel()

	for _, v := range testEpochValuesSuccess {
		if v.Err == nil {
			res := v.Epoch.Unix()
			assert.Equal(t, v.Unix, res, v.Msg)
		}
	}
}
