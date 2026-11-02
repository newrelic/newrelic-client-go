//go:build unit
// +build unit

package logconfigurations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataPartitionRuleRetentionPolicyType_Values(t *testing.T) {
	t.Parallel()

	assert.Equal(t, LogConfigurationsDataPartitionRuleRetentionPolicyType("SECONDARY"), LogConfigurationsDataPartitionRuleRetentionPolicyTypeTypes.SECONDARY)
	assert.Equal(t, LogConfigurationsDataPartitionRuleRetentionPolicyType("STANDARD"), LogConfigurationsDataPartitionRuleRetentionPolicyTypeTypes.STANDARD)
	assert.Equal(t, LogConfigurationsDataPartitionRuleRetentionPolicyType("STANDARD_LIVE_ARCHIVE"), LogConfigurationsDataPartitionRuleRetentionPolicyTypeTypes.STANDARD_LIVE_ARCHIVE)
}
