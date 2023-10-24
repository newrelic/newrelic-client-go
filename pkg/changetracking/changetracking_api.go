// Code generated by tutone: DO NOT EDIT
package changetracking

import "context"

// Creates a new deployment record in NRDB and its associated deployment marker.
func (a *Changetracking) ChangeTrackingCreateDeployment(
	dataHandlingRules ChangeTrackingDataHandlingRules,
	deployment ChangeTrackingDeploymentInput,
) (*ChangeTrackingDeployment, error) {
	return a.ChangeTrackingCreateDeploymentWithContext(context.Background(),
		dataHandlingRules,
		deployment,
	)
}

// Creates a new deployment record in NRDB and its associated deployment marker.
func (a *Changetracking) ChangeTrackingCreateDeploymentWithContext(
	ctx context.Context,
	dataHandlingRules ChangeTrackingDataHandlingRules,
	deployment ChangeTrackingDeploymentInput,
) (*ChangeTrackingDeployment, error) {

	resp := ChangeTrackingCreateDeploymentQueryResponse{}
	vars := map[string]interface{}{
		"dataHandlingRules": dataHandlingRules,
		"deployment":        deployment,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, ChangeTrackingCreateDeploymentMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.ChangeTrackingDeployment, nil
}

type ChangeTrackingCreateDeploymentQueryResponse struct {
	ChangeTrackingDeployment ChangeTrackingDeployment `json:"ChangeTrackingCreateDeployment"`
}

const ChangeTrackingCreateDeploymentMutation = `mutation(
	$dataHandlingRules: ChangeTrackingDataHandlingRules,
	$deployment: ChangeTrackingDeploymentInput!,
) { changeTrackingCreateDeployment(
	dataHandlingRules: $dataHandlingRules,
	deployment: $deployment,
) {
	changelog
	commit
	customAttributes
	deepLink
	deploymentId
	deploymentType
	description
	entityGuid
	groupId
	timestamp
	user
	version
} }`
