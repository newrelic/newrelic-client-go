package entities

import (
	"context"
	"errors"

	"github.com/newrelic/newrelic-client-go/internal/cache"
)

type entitiesCacheGroupID struct {
	queryBuilderWithoutName EntitySearchQueryBuilder
}

const entitiesCacheName = "entitiesByQueryBuilderWithoutName"

func (a *Entities) fetchEntitiesForCache(_ctx context.Context, groupID interface{}) (map[string]interface{}, error) {
	cacheGroupID, ok := groupID.(entitiesCacheGroupID)
	if !ok {
		return nil, errors.New("invalid cache group id type")
	}

	entitySearch, err := a.GetEntitySearch(
		EntitySearchOptions{},
		"",
		cacheGroupID.queryBuilderWithoutName,
		[]EntitySearchSortCriteria{},
	)
	if err != nil {
		return nil, err
	}

	itemMap := make(map[string]interface{})
	for _, entity := range entitySearch.Results.Entities {
		name := entity.GetName()
		var entities []*EntityOutlineInterface
		rawEntities, found := itemMap[name]
		if !found {
			entities = make([]*EntityOutlineInterface, 0)
		} else {
			entities, ok = rawEntities.([]*EntityOutlineInterface)
			if !ok {
				return nil, errors.New("invalid cache item type")
			}
		}
		entities = append(entities, &entity)
		itemMap[name] = entities
	}
	return itemMap, nil
}

func (a *Entities) CachedGetEntitiesByName(
	queryBuilderWithoutName EntitySearchQueryBuilder,
	name string,
) ([]*EntityOutlineInterface, error) {
	if queryBuilderWithoutName.Name != "" {
		return nil, errors.New("cannot supply Name in EntitySearchQueryBuilder to CachedGetEntitiesByName, use the 2nd argument instead")
	}

	c := cache.NamedMapGroupCache(&a.caches, entitiesCacheName, a.fetchEntitiesForCache)
	// Using context.TODO because the context isn't actually used in GetEntitySearch
	items, err := c.LookupOrFetchWithContext(context.TODO(), entitiesCacheGroupID{queryBuilderWithoutName}, name)
	if err != nil {
		return nil, err
	}
	typedItems, ok := items.([]*EntityOutlineInterface)
	if !ok {
		return nil, errors.New("invalid cache item type")
	}
	return typedItems, nil
}

func (a *Entities) InvalidateCachedEntitiesQuery(
	queryBuilderWithoutName EntitySearchQueryBuilder,
) {
	c, found := a.caches[entitiesCacheName]
	if found {
		c.InvalidateGroup(entitiesCacheGroupID{queryBuilderWithoutName})
	}
}
