// build +unit

package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

func TestIntArrayToString(t *testing.T) {
	t.Parallel()

	var result string

	// empty
	result = IntArrayToString([]int{})
	assert.Equal(t, "", result)

	// single
	result = IntArrayToString([]int{1})
	assert.Equal(t, "1", result)

	// multiple
	result = IntArrayToString([]int{1, 2, 3, 4})
	assert.Equal(t, "1,2,3,4", result)
}

func TestGetSafeTimestampWithMilliseconds_ZeroNanoseconds(t *testing.T) {
	t.Parallel()

	now := time.Now()
	actualNanoseconds := 0
	expectedNanoseconds := actualNanoseconds + 100000000

	actualResult := GetSafeTimestampWithMilliseconds(nrtime.EpochMilliseconds(
		time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour()-3,
			now.Minute()-30,
			0,
			actualNanoseconds,
			time.Local,
		),
	))

	expectedResult := nrtime.EpochMilliseconds(
		time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour()-3,
			now.Minute()-30,
			0,
			expectedNanoseconds,
			time.Local,
		),
	)
	require.Equal(t, actualResult, expectedResult)
}

func TestGetSafeTimestampWithMilliseconds_NonZeroNanoseconds(t *testing.T) {
	t.Parallel()

	now := time.Now()
	actualNanoseconds := 123000000

	// equal, since actualNanoSeconds > 100000000
	expectedNanoseconds := actualNanoseconds

	actualResult := GetSafeTimestampWithMilliseconds(nrtime.EpochMilliseconds(
		time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour()-3,
			now.Minute()-30,
			0,
			actualNanoseconds,
			time.Local,
		),
	))

	expectedResult := nrtime.EpochMilliseconds(
		time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour()-3,
			now.Minute()-30,
			0,
			expectedNanoseconds,
			time.Local,
		),
	)

	require.Equal(t, actualResult, expectedResult)

}
