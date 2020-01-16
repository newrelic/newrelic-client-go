// +build unit

package version

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var GitTag = "undefined"

func TestVersionTag(t *testing.T) {
	t.Parallel()

	log.Print("\n\n******")
	log.Printf("Current branch: %s", os.Getenv("CIRCLE_BRANCH"))
	log.Print("*******\n\n")

	assert.Equal(t, Version, GitTag)
}
