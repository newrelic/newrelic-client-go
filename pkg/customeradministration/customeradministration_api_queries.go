package customeradministration

// Get Group(s) by organization ID
func (a *Customeradministration) GetGroupsByOrganizationID(
	organizationID string,
) (*MultiTenantIdentityGroupCollection, error) {
	return a.GetGroups(
		"",
		MultiTenantIdentityGroupFilterInput{
			OrganizationId: &MultiTenantIdentityOrganizationIdInput{
				Eq: organizationID,
			},
		},
		[]MultiTenantIdentityGroupSortInput{},
	)
}
