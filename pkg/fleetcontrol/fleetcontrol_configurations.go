package fleetcontrol

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// CreateBlob creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlCreateBlob(
	requestBody interface{},
	customHeaders interface{},
	organizationID string,
) (*CreateBlobResponse, error) {
	return a.FleetControlCreateBlobWithContext(
		context.Background(),
		requestBody,
		customHeaders,
		organizationID,
	)
}

// CreatePolicyWithContext creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlCreateBlobWithContext(
	ctx context.Context,
	reqBody interface{},
	customHeaders interface{},
	organizationID string,
) (*CreateBlobResponse, error) {
	resp := CreateBlobResponse{}

	if organizationID == "" {
		return nil, fmt.Errorf("no organization ID specified")

	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	_, err = a.client.PostWithContext(
		ctx,
		a.config.Region().BlobServiceURL(fmt.Sprintf("/organizations/%s/AgentConfigurations", organizationID)),
		customHeaders,
		bytes.NewReader(reqBodyJSON),
		&resp,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type CreateBlobResponse struct {
	EntityGUID        string            `json:"entityGuid,omitempty"`
	BlobId            string            `json:"blobId,omitempty"`
	BlobVersionEntity BlobVersionEntity `json:"blobVersionEntity,omitempty"`
}

type BlobVersionEntity struct {
	BlobVersionEntityGUID string `json:"entityGuid,omitempty"`
	BlobVersion           int    `json:"version,omitempty"`
}
