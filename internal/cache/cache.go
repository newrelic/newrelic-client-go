package cache

import (
	"context"
	"fmt"
	"sync"
)

// GroupFetcherFunc is a function that accepts a group id and returns a map of item ids to items.
type GroupFetcherFunc = func(context.Context, interface{}) (map[string]interface{}, error)

// GroupCache represents a thread-safe cache for items that can be retrieved in named groups.
type GroupCache struct {
	mutex        *sync.Mutex
	groupFetcher GroupFetcherFunc
	cachedGroups map[interface{}]map[string]interface{}
}

// NewGroupCache creates an instance of GroupCache, with the given fetcher function.
func NewGroupCache(groupFetcher GroupFetcherFunc) GroupCache {
	return GroupCache{
		mutex:        &sync.Mutex{},
		groupFetcher: groupFetcher,
		cachedGroups: make(map[interface{}]map[string]interface{}),
	}
}

// NamedMapGroupCache is a convenience function for creating a cache as needed in a map of caches.
// This is useful when you have several different caches, with different fetcher functions, that all
// need to be stored in the same struct.
func NamedMapGroupCache(caches *map[string]*GroupCache, cacheName string, groupFetcher GroupFetcherFunc) *GroupCache {
	groupCachePtr, found := (*caches)[cacheName]
	if !found {
		groupCache := NewGroupCache(groupFetcher)
		groupCachePtr = &groupCache
		(*caches)[cacheName] = groupCachePtr
	}
	return groupCachePtr
}

// LookupOrFetch returns the given item, using the fetcher function if group not already fetched.
func (c *GroupCache) LookupOrFetch(groupID interface{}, itemID string) (interface{}, error) {
	return c.LookupOrFetchWithContext(context.Background(), groupID, itemID)
}

// LookupOrFetchWithContext returns the given item, using the fetcher function and given context
// if group not already fetched.
func (c *GroupCache) LookupOrFetchWithContext(ctx context.Context, groupID interface{}, itemID string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	group, found := c.cachedGroups[groupID]
	if !found {
		fetchedGroup, err := c.groupFetcher(ctx, groupID)
		if err != nil {
			return nil, err
		}
		group = fetchedGroup
		c.cachedGroups[groupID] = group
	}

	item, found := group[itemID]
	if !found {
		return nil, fmt.Errorf("item %v not found in %v", itemID, groupID)
	}
	return item, nil
}

// InvalidateGroup removes a cached group, if it was present.
func (c *GroupCache) InvalidateGroup(groupID interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cachedGroups, groupID)
}
