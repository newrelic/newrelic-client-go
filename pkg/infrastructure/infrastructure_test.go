// +build unit

package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/internal/region"
)

func TestBaseURLs(t *testing.T) {
	t.Parallel()

	pairs := map[region.Region]string{
		region.US:      "https://infra-api.newrelic.com/v2",
		region.EU:      "https://infra-api.eu.newrelic.com/v2",
		region.Staging: "https://staging-infra-api.newrelic.com/v2",
	}

	for k, v := range pairs {
		assert.Equal(t, v, BaseURLs[k])
	}
}
