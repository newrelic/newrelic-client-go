// +build unit

package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myGroupId struct {
	s string
	i int
}

type itemContent struct {
	x, y int
}

func TestGroupCache(t *testing.T) {
	fetches := 0
	cache := NewGroupCache(func(_ context.Context, groupId interface{}) (map[string]interface{}, error) {
		fetches += 1

		myGroupId, ok := groupId.(myGroupId)
		if !ok {
			return nil, errors.New("Invalid group id")
		}

		if myGroupId.s == "no-such-group" {
			return nil, errors.New("Nope")
		}

		group := map[string]interface{}{
			"foo": itemContent{1, 2},
			"bar": itemContent{3, 4},
		}
		return group, nil
	})

	// Haven't fetched anything yet
	assert.Equal(t, 0, fetches)

	// Fetching group-1:1234 for the first time
	item, err := cache.LookupOrFetch(myGroupId{"group-1", 1234}, "foo")
	assert.Equal(t, 1, fetches)
	assert.Equal(t, itemContent{1, 2}, item)
	assert.Nil(t, err)

	// Already have group-1:1234 cached, so no need to fetch this other item
	item, err = cache.LookupOrFetch(myGroupId{"group-1", 1234}, "bar")
	assert.Equal(t, 1, fetches)
	assert.Equal(t, itemContent{3, 4}, item)
	assert.Nil(t, err)

	// Fetching group-2:5678 for the first time
	item, err = cache.LookupOrFetch(myGroupId{"group-2", 5678}, "foo")
	assert.Equal(t, 2, fetches)
	assert.Equal(t, itemContent{1, 2}, item)
	assert.Nil(t, err)

	// Already have group-2:5678 cached, but it doesn't have this item id
	item, err = cache.LookupOrFetch(myGroupId{"group-2", 5678}, "wat")
	assert.Equal(t, 2, fetches)
	assert.Nil(t, item)
	assert.Equal(t, errors.New("item wat not found in {group-2 5678}"), err)

	// Trouble fetching a group
	item, err = cache.LookupOrFetch(myGroupId{"no-such-group", 5}, "wat")
	assert.Equal(t, 3, fetches)
	assert.Nil(t, item)
	assert.Equal(t, errors.New("Nope"), err)

	// Negative results aren't cached
	item, err = cache.LookupOrFetch(myGroupId{"no-such-group", 5}, "wat")
	assert.Equal(t, 4, fetches)
	assert.Nil(t, item)
	assert.Equal(t, errors.New("Nope"), err)

	// Forcing a re-fetch for group-1:1234
	cache.InvalidateGroup(myGroupId{"group-1", 1234})
	item, err = cache.LookupOrFetch(myGroupId{"group-1", 1234}, "foo")
	assert.Equal(t, 5, fetches)
	assert.Equal(t, itemContent{1, 2}, item)
	assert.Nil(t, err)
}

func TestNamedMapGroupCache(t *testing.T) {
	caches := make(map[string]*GroupCache)

	makeFakeFetcher := func(errMsg string) GroupFetcherFunc {
		return func(context.Context, interface{}) (map[string]interface{}, error) { return nil, errors.New(errMsg) }
	}

	fooFetcher := makeFakeFetcher("foo")
	foo2Fetcher := makeFakeFetcher("foo2")
	bazFetcher := makeFakeFetcher("baz")

	// Creating the foo cache for first use
	cacheFoo := NamedMapGroupCache(&caches, "foo", fooFetcher)
	assert.Equal(t, 1, len(caches))
	_, err := cacheFoo.LookupOrFetch(myGroupId{"group-1", 1234}, "abc")
	assert.Equal(t, errors.New("foo"), err)

	// Don't need to create the foo cache again, it already exists
	cacheFoo2 := NamedMapGroupCache(&caches, "foo", foo2Fetcher)
	assert.Equal(t, 1, len(caches))
	assert.Equal(t, cacheFoo, cacheFoo2)
	_, err = cacheFoo2.LookupOrFetch(myGroupId{"group-1", 1234}, "abc")
	assert.Equal(t, errors.New("foo"), err) // Not foo2, foo2Fetcher should've been ignored above

	// Creating the baz cache for first use
	cacheBaz := NamedMapGroupCache(&caches, "baz", bazFetcher)
	assert.Equal(t, 2, len(caches))
	assert.NotEqual(t, cacheFoo, cacheBaz)
	_, err = cacheBaz.LookupOrFetch(myGroupId{"group-1", 1234}, "abc")
	assert.Equal(t, errors.New("baz"), err)
}
