// Code generated by tutone: DO NOT EDIT
package authorizationmanagement

import "context"

// Grant access for a group
func (a *Authorizationmanagement) AuthorizationManagementGrantAccess(
	grantAccessOptions AuthorizationManagementGrantAccess,
) (*AuthorizationManagementGrantAccessPayload, error) {
	return a.AuthorizationManagementGrantAccessWithContext(context.Background(),
		grantAccessOptions,
	)
}

// Grant access for a group
func (a *Authorizationmanagement) AuthorizationManagementGrantAccessWithContext(
	ctx context.Context,
	grantAccessOptions AuthorizationManagementGrantAccess,
) (*AuthorizationManagementGrantAccessPayload, error) {

	resp := AuthorizationManagementGrantAccessQueryResponse{}
	vars := map[string]interface{}{
		"grantAccessOptions": grantAccessOptions,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AuthorizationManagementGrantAccessMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AuthorizationManagementGrantAccessPayload, nil
}

type AuthorizationManagementGrantAccessQueryResponse struct {
	AuthorizationManagementGrantAccessPayload AuthorizationManagementGrantAccessPayload `json:"AuthorizationManagementGrantAccess"`
}

const AuthorizationManagementGrantAccessMutation = `mutation(
	$grantAccessOptions: AuthorizationManagementGrantAccess,
) { authorizationManagementGrantAccess(
	grantAccessOptions: $grantAccessOptions,
) {
	roles {
		accountId
		displayName
		groupId
		id
		name
		organizationId
		roleId
		type
	}
} }`

// Revoke access for a group
func (a *Authorizationmanagement) AuthorizationManagementRevokeAccess(
	revokeAccessOptions AuthorizationManagementRevokeAccess,
) (*AuthorizationManagementRevokeAccessPayload, error) {
	return a.AuthorizationManagementRevokeAccessWithContext(context.Background(),
		revokeAccessOptions,
	)
}

// Revoke access for a group
func (a *Authorizationmanagement) AuthorizationManagementRevokeAccessWithContext(
	ctx context.Context,
	revokeAccessOptions AuthorizationManagementRevokeAccess,
) (*AuthorizationManagementRevokeAccessPayload, error) {

	resp := AuthorizationManagementRevokeAccessQueryResponse{}
	vars := map[string]interface{}{
		"revokeAccessOptions": revokeAccessOptions,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AuthorizationManagementRevokeAccessMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AuthorizationManagementRevokeAccessPayload, nil
}

type AuthorizationManagementRevokeAccessQueryResponse struct {
	AuthorizationManagementRevokeAccessPayload AuthorizationManagementRevokeAccessPayload `json:"AuthorizationManagementRevokeAccess"`
}

const AuthorizationManagementRevokeAccessMutation = `mutation(
	$revokeAccessOptions: AuthorizationManagementRevokeAccess,
) { authorizationManagementRevokeAccess(
	revokeAccessOptions: $revokeAccessOptions,
) {
	roles {
		accountId
		displayName
		groupId
		id
		name
		organizationId
		roleId
		type
	}
} }`