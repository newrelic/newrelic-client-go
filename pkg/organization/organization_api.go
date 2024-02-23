// Code generated by tutone: DO NOT EDIT
package organization

import "context"

// The new organization to create.
// The mutation actually only accepts one between newManagedAccount and sharedAccount, but the schema allows both so updating this manually.
// Do not accept changes from tutone about this.
func (a *Organization) OrganizationCreate(
	customerId string,
	newManagedAccount *OrganizationNewManagedAccountInput,
	organization OrganizationCreateOrganizationInput,
	sharedAccount *OrganizationSharedAccountInput,
) (*OrganizationCreateOrganizationResponse, error) {
	return a.OrganizationCreateWithContext(context.Background(),
		customerId,
		newManagedAccount,
		organization,
		sharedAccount,
	)
}

// The new organization to create.
// The mutation actually only accepts one between newManagedAccount and sharedAccount, but the schema allows both so updating this manually.
// Do not accept changes from tutone about this.
func (a *Organization) OrganizationCreateWithContext(
	ctx context.Context,
	customerId string,
	newManagedAccount *OrganizationNewManagedAccountInput,
	organization OrganizationCreateOrganizationInput,
	sharedAccount *OrganizationSharedAccountInput,
) (*OrganizationCreateOrganizationResponse, error) {

	resp := OrganizationCreateQueryResponse{}
	vars := map[string]interface{}{
		"customerId":        customerId,
		"newManagedAccount": newManagedAccount,
		"organization":      organization,
		"sharedAccount":     sharedAccount,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, OrganizationCreateMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.OrganizationCreateOrganizationResponse, nil
}

type OrganizationCreateQueryResponse struct {
	OrganizationCreateOrganizationResponse OrganizationCreateOrganizationResponse `json:"OrganizationCreate"`
}

const OrganizationCreateMutation = `mutation(
	$customerId: ID,
	$newManagedAccount: OrganizationNewManagedAccountInput,
	$organization: OrganizationCreateOrganizationInput!,
	$sharedAccount: OrganizationSharedAccountInput,
) { organizationCreate(
	customerId: $customerId,
	newManagedAccount: $newManagedAccount,
	organization: $organization,
	sharedAccount: $sharedAccount,
) {
	jobId
} }`

// Accessible organizations
func (a *Organization) GetOrganizations(
	cursor string,
	filter OrganizationCustomerOrganizationFilterInput,
) (*OrganizationCustomerOrganizationWrapper, error) {
	return a.GetOrganizationsWithContext(context.Background(),
		cursor,
		filter,
	)
}

// Accessible organizations
func (a *Organization) GetOrganizationsWithContext(
	ctx context.Context,
	cursor string,
	filter OrganizationCustomerOrganizationFilterInput,
) (*OrganizationCustomerOrganizationWrapper, error) {

	resp := organizationsResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getOrganizationsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Organizations, nil
}

const getOrganizationsQuery = `query { customerAdministration { organizations {
	items {
		contractId
		customerId
		id
		name
	}
	nextCursor
} } }`
