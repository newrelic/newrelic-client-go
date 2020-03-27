// +build unit

package region

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	pairs := map[string]*Region{
		"us":      Regions[US],
		"Us":      Regions[US],
		"uS":      Regions[US],
		"US":      Regions[US],
		"eu":      Regions[EU],
		"Eu":      Regions[EU],
		"eU":      Regions[EU],
		"EU":      Regions[EU],
		"staging": Regions[Staging],
		"Staging": Regions[Staging],
		"STAGING": Regions[Staging],
	}

	for k, v := range pairs {
		result := Parse(k)
		assert.Equal(t, result, v)
	}

	// Default is US
	result := Parse("")
	assert.Equal(t, result, Regions[US])
}

func TestString(t *testing.T) {
	t.Parallel()

	pairs := map[Name]string{
		US:      "US",
		EU:      "EU",
		Staging: "Staging",
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

func TestInfrastructurURLs(t *testing.T) {
	t.Parallel()

	pairs := map[Name]string{
		US:      "https://infra-api.newrelic.com/v2",
		EU:      "https://infra-api.eu.newrelic.com/v2",
		Staging: "https://staging-infra-api.newrelic.com/v2",
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
	}

	for k, v := range pairs {
		assert.Equal(t, v, Regions[k].SyntheticsURL())
	}
}

func TestConcatURLPaths(t *testing.T) {
	t.Parallel()

	res, err := concatURLPaths("http://localhost/", []string{"one", "/two", "//three", "four/", "five//"})

	assert.NoError(t, err)
	assert.Equal(t, "http://localhost/one/two/three/four/five", res)
}
