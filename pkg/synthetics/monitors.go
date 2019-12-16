package synthetics

import (
	"strconv"
)

const (
	listMonitorsLimit = 100
)

type listMonitorsResponse struct {
	Monitors []Monitor `json:"monitors,omitempty"`
}

// ListMonitors is used to retrieve New Relic Synthetics monitors.
func (s *Synthetics) ListMonitors() ([]Monitor, error) {
	res := listMonitorsResponse{}
	paramsMap := map[string]string{
		"limit": strconv.Itoa(listMonitorsLimit),
	}

	responses, err := s.client.GetMultiple("/monitors", &paramsMap, &res)

	monitors := []Monitor{}
	for _, r := range responses {
		if response, ok := r.(*listMonitorsResponse); ok {
			monitors = append(monitors, response.Monitors...)
		}
	}

	if err != nil {
		return nil, err
	}

	return monitors, nil
}
