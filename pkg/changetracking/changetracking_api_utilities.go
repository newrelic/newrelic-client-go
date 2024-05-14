package changetracking

import (
	"time"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

// DO NOT DELETE the following function - this is not covered by Tutone
// but is needed to ensure proper conversion of timestamps
func (input *ChangeTrackingDeploymentInput) CorrectTimestampMilliseconds() {
	inputTimestamp := input.Timestamp
	timestamp := time.Time(inputTimestamp)

	// since time.Time in Go does not have a milliseconds field, which is why the implementation
	// of unmarshaling time.Time into a Unix timestamp in the serialization package relies on
	// nanoseconds to produce a value of milliseconds, we try employing a similar logic below

	if timestamp.Nanosecond() < 100000000 {
		timestamp = timestamp.Add(time.Nanosecond * 100000000)
	}

	input.Timestamp = nrtime.EpochMilliseconds(timestamp)
}
