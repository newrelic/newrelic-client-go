//go:build unit
// +build unit

package serialization

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalEpochMillisecondsJSON_CodyCase(t *testing.T) {
	t.Parallel()

	tm, err := time.Parse(time.RFC3339, "2021-05-17T21:28:04Z")
	require.NoError(t, err)

	b, err := MarshalEpochMillisecondsJSON(tm)
	require.NoError(t, err)
	assert.Equal(t, []byte("1621286884000"), b)
}

func TestUnmarshalEpochMillisecondsJSON_RFC3339(t *testing.T) {
	t.Parallel()

	want, err := time.Parse(time.RFC3339, "2021-05-17T21:28:04Z")
	require.NoError(t, err)

	for name, input := range map[string][]byte{
		"unquoted": []byte("2021-05-17T21:28:04Z"),
		"json":     []byte(`"2021-05-17T21:28:04Z"`),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var et EpochTime
			require.NoError(t, UnmarshalEpochMillisecondsJSON(input, &et))
			assert.True(t, time.Time(et).Equal(want.UTC()))
		})
	}
}

func TestUnmarshalEpochMillisecondsJSON_NumericDelegatesToEpochTime(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name  string
		input []byte
		want  time.Time
	}{
		{
			name:  "seconds",
			input: []byte(`1587654321`),
			want:  time.Unix(1587654321, 0).UTC(),
		},
		{
			name:  "milliseconds",
			input: []byte(`1587654321012`),
			want:  time.Unix(1587654321, 12*int64(time.Millisecond)).UTC(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var et EpochTime
			require.NoError(t, UnmarshalEpochMillisecondsJSON(tc.input, &et))
			assert.True(t, time.Time(et).Equal(tc.want))
		})
	}
}

func TestUnmarshalEpochMillisecondsJSON_QuotedNumeric(t *testing.T) {
	t.Parallel()

	var et EpochTime
	require.NoError(t, UnmarshalEpochMillisecondsJSON([]byte(`"1587654321012"`), &et))
	want := time.Unix(1587654321, 12*int64(time.Millisecond)).UTC()
	assert.True(t, time.Time(et).Equal(want))
}

func TestEpochSecondsMarshalJSON_StillTenDigitsWithoutSubsecond(t *testing.T) {
	t.Parallel()

	et := EpochTime(time.Unix(1587654321, 0).UTC())
	b, err := et.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte(`1587654321`), b)
}
