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
