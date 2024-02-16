// Code generated by tutone: DO NOT EDIT
package customeradministration

import "context"

// Accessible account shares
func (a *Customeradministration) GetAccountShares(
	cursor string,
	filter OrganizationAccountShareFilterInput,
	sort []OrganizationAccountShareSortInput,
) (*OrganizationAccountShareCollection, error) {
	return a.GetAccountSharesWithContext(context.Background(),
		cursor,
		filter,
		sort,
	)
}

// Accessible account shares
func (a *Customeradministration) GetAccountSharesWithContext(
	ctx context.Context,
	cursor string,
	filter OrganizationAccountShareFilterInput,
	sort []OrganizationAccountShareSortInput,
) (*OrganizationAccountShareCollection, error) {

	resp := accountSharesResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
		"sort":   sort,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getAccountSharesQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.AccountShares, nil
}

const getAccountSharesQuery = `query(
	$filter: OrganizationAccountShareFilterInput!,
	$sort: [OrganizationAccountShareSortInput!],
) { customerAdministration { accountShares(
	filter: $filter,
	sort: $sort,
) {
	items {
		accountId
		id
		limitingRole {
			id
		}
		name
		source {
			id
			name
		}
		target {
			id
			name
		}
	}
	nextCursor
} } }`

// accounts
func (a *Customeradministration) GetAccounts(
	cursor string,
	filter OrganizationAccountFilterInput,
	sort []OrganizationAccountSortInput,
) (*OrganizationAccountCollection, error) {
	return a.GetAccountsWithContext(context.Background(),
		cursor,
		filter,
		sort,
	)
}

// accounts
func (a *Customeradministration) GetAccountsWithContext(
	ctx context.Context,
	cursor string,
	filter OrganizationAccountFilterInput,
	sort []OrganizationAccountSortInput,
) (*OrganizationAccountCollection, error) {

	resp := accountsResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
		"sort":   sort,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getAccountsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Accounts, nil
}

const getAccountsQuery = `query(
	$filter: OrganizationAccountFilterInput!,
	$sort: [OrganizationAccountSortInput!],
) { customerAdministration { accounts(
	filter: $filter,
	sort: $sort,
) {
	items {
		id
		name
		regionCode
		status
	}
	nextCursor
	totalCount
} } }`

// Authentication domains
func (a *Customeradministration) GetAuthenticationDomains(
	cursor string,
	filter OrganizationAuthenticationDomainFilterInput,
	sort []OrganizationAuthenticationDomainSortInput,
) (*OrganizationAuthenticationDomainCollection, error) {
	return a.GetAuthenticationDomainsWithContext(context.Background(),
		cursor,
		filter,
		sort,
	)
}

// Authentication domains
func (a *Customeradministration) GetAuthenticationDomainsWithContext(
	ctx context.Context,
	cursor string,
	filter OrganizationAuthenticationDomainFilterInput,
	sort []OrganizationAuthenticationDomainSortInput,
) (*OrganizationAuthenticationDomainCollection, error) {

	resp := authenticationDomainsResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
		"sort":   sort,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getAuthenticationDomainsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.AuthenticationDomains, nil
}

const getAuthenticationDomainsQuery = `query(
	$filter: OrganizationAuthenticationDomainFilterInput!,
	$sort: [OrganizationAuthenticationDomainSortInput!],
) { customerAdministration { authenticationDomains(
	filter: $filter,
	sort: $sort,
) {
	items {
		authenticationType
		id
		name
		organizationId
		provisioningType
	}
	nextCursor
} } }`

// The `consumption` field is the entry point into a customer's consumption data that is scoped to the ID of the customer.
func (a *Customeradministration) GetConsumption(
	customerId string,
) (*Consumption, error) {
	return a.GetConsumptionWithContext(context.Background(),
		customerId,
	)
}

// The `consumption` field is the entry point into a customer's consumption data that is scoped to the ID of the customer.
func (a *Customeradministration) GetConsumptionWithContext(
	ctx context.Context,
	customerId string,
) (*Consumption, error) {

	resp := consumptionResponse{}
	vars := map[string]interface{}{
		"customerId": customerId,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getConsumptionQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Consumption, nil
}

const getConsumptionQuery = `query(
	$customerId: ID!,
) { customerAdministration { consumption(
	customerId: $customerId,
) {
	customerId
} } }`

// Accessible contracts
func (a *Customeradministration) GetContracts(
	cursor string,
	filter OrganizationCustomerContractFilterInput,
) (*OrganizationCustomerContractWrapper, error) {
	return a.GetContractsWithContext(context.Background(),
		cursor,
		filter,
	)
}

// Accessible contracts
func (a *Customeradministration) GetContractsWithContext(
	ctx context.Context,
	cursor string,
	filter OrganizationCustomerContractFilterInput,
) (*OrganizationCustomerContractWrapper, error) {

	resp := contractsResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getContractsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Contracts, nil
}

const getContractsQuery = `query { customerAdministration { contracts {
	items {
		billingStructure
		customerId
		id
		organizationGroups {
			items {
				id
				name
			}
			nextCursor
		}
		telemetryId
	}
	nextCursor
} } }`

// list of grants
func (a *Customeradministration) GetGrants(
	cursor string,
	filter MultiTenantAuthorizationGrantFilterInputExpression,
	sort []MultiTenantAuthorizationGrantSortInput,
) (*MultiTenantAuthorizationGrantCollection, error) {
	return a.GetGrantsWithContext(context.Background(),
		cursor,
		filter,
		sort,
	)
}

// list of grants
func (a *Customeradministration) GetGrantsWithContext(
	ctx context.Context,
	cursor string,
	filter MultiTenantAuthorizationGrantFilterInputExpression,
	sort []MultiTenantAuthorizationGrantSortInput,
) (*MultiTenantAuthorizationGrantCollection, error) {

	resp := grantsResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
		"sort":   sort,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getGrantsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Grants, nil
}

const getGrantsQuery = `query(
	$filter: MultiTenantAuthorizationGrantFilterInputExpression!,
	$sort: [MultiTenantAuthorizationGrantSortInput!],
) { customerAdministration { grants(
	filter: $filter,
	sort: $sort,
) {
	items {
		group {
			id
		}
		id
		role {
			id
			name
		}
		scope {
			id
			type
		}
	}
	nextCursor
} } }`

// Named sets of New Relic users within an authentication domain
func (a *Customeradministration) GetGroups(
	cursor string,
	filter MultiTenantIdentityGroupFilterInput,
	sort []MultiTenantIdentityGroupSortInput,
) (*MultiTenantIdentityGroupCollection, error) {
	return a.GetGroupsWithContext(context.Background(),
		cursor,
		filter,
		sort,
	)
}

// Named sets of New Relic users within an authentication domain
func (a *Customeradministration) GetGroupsWithContext(
	ctx context.Context,
	cursor string,
	filter MultiTenantIdentityGroupFilterInput,
	sort []MultiTenantIdentityGroupSortInput,
) (*MultiTenantIdentityGroupCollection, error) {

	resp := groupsResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
		"sort":   sort,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getGroupsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Groups, nil
}

const getGroupsQuery = `query(
	$filter: MultiTenantIdentityGroupFilterInput!,
	$sort: [MultiTenantIdentityGroupSortInput!],
) { customerAdministration { groups(
	filter: $filter,
	sort: $sort,
) {
	items {
		authenticationDomainId
		id
		name
		users {
			items {
				email
				id
				name
				timeZone
			}
			nextCursor
			totalCount
		}
	}
	nextCursor
	totalCount
} } }`

// This provides access to fields you can use to check the status of asynchronous jobs related to customer administration.
func (a *Customeradministration) GetJobs() (*CustomerAdministrationJobs, error) {
	return a.GetJobsWithContext(context.Background())
}

// This provides access to fields you can use to check the status of asynchronous jobs related to customer administration.
func (a *Customeradministration) GetJobsWithContext(
	ctx context.Context,
) (*CustomerAdministrationJobs, error) {

	resp := jobsResponse{}
	vars := map[string]interface{}{}

	if err := a.client.NerdGraphQueryWithContext(ctx, getJobsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Jobs, nil
}

const getJobsQuery = `query { customerAdministration { jobs {
	
} } }`

// Accessible organizations
func (a *Customeradministration) GetOrganizations(
	cursor string,
	filter OrganizationCustomerOrganizationFilterInput,
) (*OrganizationCustomerOrganizationWrapper, error) {
	return a.GetOrganizationsWithContext(context.Background(),
		cursor,
		filter,
	)
}

// Accessible organizations
func (a *Customeradministration) GetOrganizationsWithContext(
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

// list of permissions
func (a *Customeradministration) GetPermissions(
	cursor string,
	filter MultiTenantAuthorizationPermissionFilter,
) (*MultiTenantAuthorizationPermissionCollection, error) {
	return a.GetPermissionsWithContext(context.Background(),
		cursor,
		filter,
	)
}

// list of permissions
func (a *Customeradministration) GetPermissionsWithContext(
	ctx context.Context,
	cursor string,
	filter MultiTenantAuthorizationPermissionFilter,
) (*MultiTenantAuthorizationPermissionCollection, error) {

	resp := permissionsResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getPermissionsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Permissions, nil
}

const getPermissionsQuery = `query { customerAdministration { permissions {
	items {
		category
		feature
		id
		name
		product
	}
	nextCursor
} } }`

// list of roles
func (a *Customeradministration) GetRoles(
	cursor string,
	filter MultiTenantAuthorizationRoleFilterInputExpression,
	sort []MultiTenantAuthorizationRoleSortInput,
) (*MultiTenantAuthorizationRoleCollection, error) {
	return a.GetRolesWithContext(context.Background(),
		cursor,
		filter,
		sort,
	)
}

// list of roles
func (a *Customeradministration) GetRolesWithContext(
	ctx context.Context,
	cursor string,
	filter MultiTenantAuthorizationRoleFilterInputExpression,
	sort []MultiTenantAuthorizationRoleSortInput,
) (*MultiTenantAuthorizationRoleCollection, error) {

	resp := rolesResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
		"sort":   sort,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getRolesQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Roles, nil
}

const getRolesQuery = `query(
	$filter: MultiTenantAuthorizationRoleFilterInputExpression!,
	$sort: [MultiTenantAuthorizationRoleSortInput!],
) { customerAdministration { roles(
	filter: $filter,
	sort: $sort,
) {
	items {
		id
		name
		scope
		type
	}
	nextCursor
	totalCount
} } }`

// The authenticated `User` who made this request.
func (a *Customeradministration) GetUser() (*User, error) {
	return a.GetUserWithContext(context.Background())
}

// The authenticated `User` who made this request.
func (a *Customeradministration) GetUserWithContext(
	ctx context.Context,
) (*User, error) {

	resp := userResponse{}
	vars := map[string]interface{}{}

	if err := a.client.NerdGraphQueryWithContext(ctx, getUserQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.User, nil
}

const getUserQuery = `query { customerAdministration { user {
	email
	id
	name
} } }`

// A collection of New Relic users
func (a *Customeradministration) GetUsers(
	cursor string,
	filter MultiTenantIdentityUserFilterInput,
	sort []MultiTenantIdentityUserSortInput,
) (*MultiTenantIdentityUserCollection, error) {
	return a.GetUsersWithContext(context.Background(),
		cursor,
		filter,
		sort,
	)
}

// A collection of New Relic users
func (a *Customeradministration) GetUsersWithContext(
	ctx context.Context,
	cursor string,
	filter MultiTenantIdentityUserFilterInput,
	sort []MultiTenantIdentityUserSortInput,
) (*MultiTenantIdentityUserCollection, error) {

	resp := usersResponse{}
	vars := map[string]interface{}{
		"cursor": cursor,
		"filter": filter,
		"sort":   sort,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getUsersQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CustomerAdministration.Users, nil
}

const getUsersQuery = `query(
	$filter: MultiTenantIdentityUserFilterInput!,
	$sort: [MultiTenantIdentityUserSortInput!],
) { customerAdministration { users(
	filter: $filter,
	sort: $sort,
) {
	items {
		authenticationDomainId
		email
		emailVerificationState
		groups {
			items {
				id
				name
			}
			nextCursor
			totalCount
		}
		id
		lastActive
		name
		pendingUpgradeRequest {
			id
			message
			requestedUserType {
				id
				name
			}
		}
		timeZone
		type {
			id
			name
		}
	}
	nextCursor
	totalCount
} } }`
