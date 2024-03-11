package organization

import "context"

func (a *Organization) OrganizationGetOrganizations(ctx context.Context) (*OrganizationCustomerOrganizationWrapper, error) {
	return a.GetOrganizationsWithContext(ctx, "", OrganizationCustomerOrganizationFilterInput{})
}
