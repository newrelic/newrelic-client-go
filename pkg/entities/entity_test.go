//go:build unit

package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindTagByKey(t *testing.T) {
	t.Parallel()

	entityTags := []EntityTag{
		{
			Key:    "test",
			Values: []string{"someTag"},
		},
	}

	result := FindTagByKey(entityTags, "test")
	require.Equal(t, []string{"someTag"}, result)
}

func TestFindTagByKeyNotFound(t *testing.T) {
	t.Parallel()

	entityTags := []EntityTag{
		{
			Key:    "test",
			Values: []string{"someTag"},
		},
	}

	result := FindTagByKey(entityTags, "notFound")
	require.Empty(t, result)
}
