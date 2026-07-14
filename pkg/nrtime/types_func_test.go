//go:build unit
// +build unit

package nrtime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEpochMillisecondsJSONRoundTrip(t *testing.T) {
	t.Parallel()

	tm, err := time.Parse(time.RFC3339, "2021-05-17T21:28:04Z")
	require.NoError(t, err)

	ms := EpochMilliseconds(tm)
	out, err := json.Marshal(ms)
	require.NoError(t, err)
	assert.Equal(t, "1621286884000", string(out))

	var back EpochMilliseconds
	require.NoError(t, json.Unmarshal(out, &back))
	assert.True(t, time.Time(back).Equal(tm.UTC()))
}
