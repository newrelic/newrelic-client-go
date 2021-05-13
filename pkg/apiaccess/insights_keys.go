package apiaccess

import (
	"fmt"
	"time"
)

type InsightsKey struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	Key       string    `json:"key"`
	Notes     string    `json:"notes"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *APIAccess) ListInsightsInsertKeys(accountID int) ([]InsightsKey, error) {
	keys := []InsightsKey{}
	url := a.config.Region().InsightsKeysURL(accountID, "insert_keys?format=json")
	_, err := a.insightsKeysClient.Get(url, nil, &keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (a *APIAccess) CreateInsightsInsertKey(accountID int) (*InsightsKey, error) {
	key := InsightsKey{}
	url := a.config.Region().InsightsKeysURL(accountID, "insert_keys?format=json")
	_, err := a.insightsKeysClient.Post(url, nil, nil, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (a *APIAccess) GetInsightsInsertKey(accountID int, keyID int) (*InsightsKey, error) {
	key := InsightsKey{}
	url := a.config.Region().InsightsKeysURL(accountID, fmt.Sprintf("insert_keys/%d?format=json", keyID))
	_, err := a.insightsKeysClient.Get(url, nil, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (a *APIAccess) DeleteInsightsInsertKey(accountID int, keyID int) (*InsightsKey, error) {
	key := InsightsKey{}
	url := a.config.Region().InsightsKeysURL(accountID, fmt.Sprintf("insert_keys/%d?format=json", keyID))
	_, err := a.insightsKeysClient.Delete(url, nil, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (a *APIAccess) ListInsightsQueryKeys(accountID int) ([]InsightsKey, error) {
	keys := []InsightsKey{}
	url := a.config.Region().InsightsKeysURL(accountID, "query_keys?format=json")
	_, err := a.insightsKeysClient.Get(url, nil, &keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (a *APIAccess) CreateInsightsQueryKey(accountID int) (*InsightsKey, error) {
	key := InsightsKey{}
	url := a.config.Region().InsightsKeysURL(accountID, "query_keys?format=json")
	_, err := a.insightsKeysClient.Post(url, nil, nil, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (a *APIAccess) GetInsightsQueryKey(accountID int, keyID int) (*InsightsKey, error) {
	key := InsightsKey{}
	url := a.config.Region().InsightsKeysURL(accountID, fmt.Sprintf("query_keys/%d?format=json", keyID))
	_, err := a.insightsKeysClient.Get(url, nil, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (a *APIAccess) DeleteInsightsQueryKey(accountID int, keyID int) (*InsightsKey, error) {
	key := InsightsKey{}
	url := a.config.Region().InsightsKeysURL(accountID, fmt.Sprintf("query_keys/%d?format=json", keyID))
	_, err := a.insightsKeysClient.Delete(url, nil, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}
