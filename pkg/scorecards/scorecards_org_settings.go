package scorecards

// scorecards_org_settings.go provides a purpose-specific helper for fetching
// the TeamsOrganizationSettings entity for the authenticated organisation.
//
// Callers (primarily the Terraform provider) need this to retrieve
// discovery.tagKeys so they can distinguish dynamically tag-assigned
// ownership entities from statically managed ones — the same lookup the
// Teams UI makes via useOrganizationSettingsQuery (EntityOwnership.tsx).
//
// The underlying query reuses the generated GetEntitySearch method with the
// type filter "type = 'TEAMS_ORGANIZATION_SETTINGS'" rather than introducing
// a separate query string, keeping the generated code surface clean.

import "context"

// GetTeamsOrganizationSettings fetches the singleton
// EntityManagementTeamsOrganizationSettingsEntity for the authenticated
// organisation and returns it typed. Returns nil (no error) when no settings
// entity exists yet (e.g. the org has never configured team settings).
func (a *Scorecards) GetTeamsOrganizationSettings() (*EntityManagementTeamsOrganizationSettingsEntity, error) {
	return a.GetTeamsOrganizationSettingsWithContext(context.Background())
}

// GetTeamsOrganizationSettingsWithContext is the context-aware form of
// GetTeamsOrganizationSettings.
func (a *Scorecards) GetTeamsOrganizationSettingsWithContext(
	ctx context.Context,
) (*EntityManagementTeamsOrganizationSettingsEntity, error) {

	result, err := a.GetEntitySearchWithContext(ctx, "", "type = 'TEAMS_ORGANIZATION_SETTINGS'")
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	for _, e := range result.Entities {
		if s, ok := e.(*EntityManagementTeamsOrganizationSettingsEntity); ok {
			return s, nil
		}
	}

	// Org has not configured team settings yet.
	return nil, nil
}
