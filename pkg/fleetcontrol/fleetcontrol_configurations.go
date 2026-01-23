package fleetcontrol

import (
	"context"
	"fmt"
)

// CreateBlob creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlGetConfiguration(
	entityGUID string,
	organizationID string,
	getConfigurationMode GetConfigurationMode,
	version int,
) (*GetConfigurationResponse, error) {
	return a.FleetControlGetConfigurationWithContext(
		context.Background(),
		entityGUID,
		organizationID,
		getConfigurationMode,
		version,
	)
}

// CreatePolicyWithContext creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlGetConfigurationWithContext(
	ctx context.Context,
	entityGUID string,
	organizationID string,
	getConfigurationMode GetConfigurationMode,
	version int,
) (*GetConfigurationResponse, error) {
	var resp GetConfigurationResponse

	if organizationID == "" {
		return nil, fmt.Errorf("no organization ID specified")

	}

	versionQueryParameterAppender := ""
	if version >= 1 {
		versionQueryParameterAppender = fmt.Sprintf("?version=%d", version)
	}

	_, err := a.client.GetWithContext(
		ctx,
		a.config.Region().BlobServiceURL(
			fmt.Sprintf(
				"/organizations/%s/%s/%s%s",
				organizationID,
				string(getConfigurationMode),
				entityGUID,
				versionQueryParameterAppender,
			)),
		nil,
		&resp,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateBlob creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlCreateConfiguration(
	requestBody interface{},
	customHeaders interface{},
	organizationID string,
) (*CreateConfigurationResponse, error) {
	return a.FleetControlCreateConfigurationWithContext(
		context.Background(),
		requestBody,
		customHeaders,
		organizationID,
	)
}

// CreatePolicyWithContext creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlCreateConfigurationWithContext(
	ctx context.Context,
	reqBody interface{},
	customHeaders interface{},
	organizationID string,
) (*CreateConfigurationResponse, error) {
	resp := CreateConfigurationResponse{}

	if organizationID == "" {
		return nil, fmt.Errorf("no organization ID specified")

	}

	_, err := a.client.PostWithContext(
		ctx,
		a.config.Region().BlobServiceURL(fmt.Sprintf("/organizations/%s/AgentConfigurations", organizationID)),
		customHeaders,
		reqBody,
		&resp,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type CreateConfigurationResponse struct {
	BlobId                  string                     `json:"blobId,omitempty"`
	ConfigurationEntityGUID string                     `json:"entityGuid,omitempty"`
	ConfigurationVersion    ConfigurationVersionEntity `json:"blobVersionEntity,omitempty"`
}

type GetConfigurationResponse string

type DeleteBlobResponse struct {
	Response string `json:"response,omitempty"`
}

type ConfigurationVersionEntity struct {
	ConfigurationVersionEntityGUID string `json:"entityGuid,omitempty"`
	ConfigurationVersionNumber     int    `json:"version,omitempty"`
}

type GetConfigurationMode string

var GetConfigurationModeTypes = struct {
	ConfigEntity        GetConfigurationMode
	ConfigVersionEntity GetConfigurationMode
}{
	ConfigEntity:        "AgentConfigurations",
	ConfigVersionEntity: "AgentConfigurationVersions",
}

// CreateBlob creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlDeleteConfiguration(
	blobEntityGUID string,
	organizationID string,
) (*DeleteBlobResponse, error) {
	return a.FleetControlDeleteConfigurationWithContext(
		context.Background(),
		blobEntityGUID,
		organizationID,
	)
}

// CreatePolicyWithContext creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlDeleteConfigurationWithContext(
	ctx context.Context,
	blobEntityGUID string,
	organizationID string,
) (*DeleteBlobResponse, error) {
	resp := DeleteBlobResponse{}

	if organizationID == "" {
		return nil, fmt.Errorf("no organization ID specified")

	}

	x, err := a.client.DeleteWithContext(
		ctx,
		a.config.Region().BlobServiceURL(fmt.Sprintf("/organizations/%s/AgentConfigurations/%s", organizationID, blobEntityGUID)),
		nil,
		&resp,
	)

	if err != nil {
		fmt.Println(x)
		return nil, err
	}

	return &resp, nil
}

// CreateBlob creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlDeleteConfigurationVersion(
	configurationVersionGUID string,
	organizationID string,
) error {
	return a.FleetControlDeleteConfigurationVersionWithContext(
		context.Background(),
		configurationVersionGUID,
		organizationID,
	)
}

// CreatePolicyWithContext creates a new alert policy for a given account.
func (a *Fleetcontrol) FleetControlDeleteConfigurationVersionWithContext(
	ctx context.Context,
	configurationVersionGUID string,
	organizationID string,
) error {
	resp := DeleteBlobResponse{}

	if organizationID == "" {
		return fmt.Errorf("no organization ID specified")

	}

	x, err := a.client.DeleteWithContext(
		ctx,
		a.config.Region().BlobServiceURL(fmt.Sprintf("/organizations/%s/AgentConfigurationVersions/%s", organizationID, configurationVersionGUID)),
		nil,
		&resp,
	)

	if err != nil {
		fmt.Println(x)
		return err
	}

	return nil
}
