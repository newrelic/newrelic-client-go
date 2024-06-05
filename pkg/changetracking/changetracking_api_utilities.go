package changetracking

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/utils"
)

// DO NOT DELETE the following function - this is not covered by Tutone
// but is needed to ensure proper conversion of timestamps
func (input *ChangeTrackingDeploymentInput) CorrectTimestampMilliseconds() {
	inputTimestamp := input.Timestamp
	timestamp := utils.GetSafeTimestampWithMilliseconds(inputTimestamp)
	input.Timestamp = timestamp
}
