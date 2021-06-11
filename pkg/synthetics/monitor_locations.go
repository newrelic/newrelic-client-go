package synthetics

import "context"

// MonitorLocation represents a valid location for a New Relic Synthetics monitor.
type MonitorLocation struct {
	HighSecurityMode bool   `json:"highSecurityMode"`
	Private          bool   `json:"private"`
	Name             string `json:"name"`
	Label            string `json:"label"`
	Description      string `json:"description"`
}

// GetMonitorLocations is used to retrieve all valid locations for Synthetics monitors.
func (s *Synthetics) GetMonitorLocations() ([]*MonitorLocation, error) {
	return s.GetMonitorLocationsWithContext(context.Background())
}

// GetMonitorLocationsWithContext is used to retrieve all valid locations for Synthetics monitors.
func (s *Synthetics) GetMonitorLocationsWithContext(ctx context.Context) ([]*MonitorLocation, error) {
	url := "/v1/locations"

	resp := []*MonitorLocation{}

	_, err := s.client.GetWithContext(ctx, s.config.Region().SyntheticsURL(url), nil, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
