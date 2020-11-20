// +build unit

package region

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	pairs := map[string]Name{
		"us":      US,
		"Us":      US,
		"uS":      US,
		"US":      US,
		"eu":      EU,
		"Eu":      EU,
		"eU":      EU,
		"EU":      EU,
		"staging": Staging,
		"Staging": Staging,
		"STAGING": Staging,
		"local":   Local,
		"Local":   Local,
		"LOCAL":   Local,
	}

	for k, v := range pairs {
		result, err := Parse(k)
		assert.NoError(t, err)
		assert.Equal(t, v, result)
	}

	// Default is US
	result, err := Parse("")
	assert.Error(t, err)
	assert.IsType(t, UnknownError{}, err)
	assert.Equal(t, Name(""), result)
}

func TestRegionGet(t *testing.T) {
	t.Parallel()

	pairs := map[Name]*Region{
		US:      Regions[US],
		EU:      Regions[EU],
		Staging: Regions[Staging],
	}

	for k, v := range pairs {
		result, err := Get(k)
		assert.NoError(t, err)
		assert.Equal(t, v, result)
	}

	// Throws error, still returns the default
	var unk Name = "(unknown)"
	result, err := Get(unk)
	assert.Error(t, err)
	assert.IsType(t, UnknownUsingDefaultError{}, err)
	assert.Equal(t, Regions[Default], result)
}

func TestRegionString(t *testing.T) {
	t.Parallel()

	pairs := map[Name]string{
		US:      "US",
		EU:      "EU",
		Staging: "Staging",
		Local:   "Local",
	}

	for k, v := range pairs {
		result := Regions[k].String()
		assert.Equal(t, result, v)
	}

	// Verify that an uninitialized Region (should be 0) isn't known
	var unk Region
	result := unk.String()
	assert.Equal(t, result, "(Unknown)")
}

func TestInfrastructureURLs(t *testing.T) {
	t.Parallel()

	pairs := map[Name]string{
		US:      "https://infra-api.newrelic.com/v2",
		EU:      "https://infra-api.eu.newrelic.com/v2",
		Staging: "https://staging-infra-api.newrelic.com/v2",
		Local:   "http://localhost:3000/v2",
	}

	for k, v := range pairs {
		assert.Equal(t, v, Regions[k].InfrastructureURL())
	}
}

func TestSyntheticsURLs(t *testing.T) {
	t.Parallel()

	pairs := map[Name]string{
		US:      "https://synthetics.newrelic.com/synthetics/api",
		EU:      "https://synthetics.eu.newrelic.com/synthetics/api",
		Staging: "https://staging-synthetics.newrelic.com/synthetics/api",
		Local:   "http://localhost:3000/synthetics/api",
	}

	for k, v := range pairs {
		assert.Equal(t, v, Regions[k].SyntheticsURL())
	}
}

func TestLogsURLs(t *testing.T) {
	t.Parallel()

	pairs := map[Name]string{
		US:      "https://log-api.newrelic.com/log/v1",
		EU:      "https://log-api.eu.newrelic.com/log/v1",
		Staging: "https://staging-log-api.newrelic.com/log/v1",
		Local:   "http://localhost:3000/log/v1",
	}

	for k, v := range pairs {
		assert.Equal(t, v, Regions[k].LogsURL())
	}
}

func TestNerdgraphURLs(t *testing.T) {
	t.Parallel()

	pairs := map[Name]string{
		US:      "https://api.newrelic.com/graphql",
		EU:      "https://api.eu.newrelic.com/graphql",
		Staging: "https://staging-api.newrelic.com/graphql",
		Local:   "http://localhost:3000/graphql",
	}

	for k, v := range pairs {
		assert.Equal(t, v, Regions[k].NerdGraphURL())
	}
}

func TestRESTURLs(t *testing.T) {
	t.Parallel()

	pairs := map[Name]string{
		US:      "https://api.newrelic.com/v2",
		EU:      "https://api.eu.newrelic.com/v2",
		Staging: "https://staging-api.newrelic.com/v2",
		Local:   "http://localhost:3000/v2",
	}

	for k, v := range pairs {
		assert.Equal(t, v, Regions[k].RestURL())
	}
}

func TestConcatURLPaths(t *testing.T) {
	t.Parallel()

	res, err := concatURLPaths("http://localhost/", []string{"one", "/two", "//three", "four/", "five//"})

	assert.NoError(t, err)
	assert.Equal(t, "http://localhost/one/two/three/four/five", res)
}
