package apm

import (
	"strconv"
	"strings"
)

// ListApplicationsParams represents a set of filters to be
// used when querying New Relic applications.
type ListApplicationsParams struct {
	Name     *string
	Host     *string
	IDs      []int
	Language *string
}

type listApplicationsResponse struct {
	Applications []Application `json:"applications,omitempty"`
}

// ListApplications is used to retrieve New Relic applications.
func (apm *APM) ListApplications(params *ListApplicationsParams) ([]Application, error) {
	res := listApplicationsResponse{}
	paramsMap := buildListApplicationsParamsMap(params)
	responses, err := apm.client.GetMultiple("/applications.json", &paramsMap, &res)

	applications := []Application{}
	for _, r := range responses {
		if response, ok := r.(*listApplicationsResponse); ok {
			applications = append(applications, response.Applications...)
		}
	}

	if err != nil {
		return nil, err
	}

	return applications, nil
}

func buildListApplicationsParamsMap(params *ListApplicationsParams) map[string]string {
	paramsMap := map[string]string{}

	if params != nil {
		if params.Name != nil {
			paramsMap["filter[name]"] = *params.Name
		}

		if params.Host != nil {
			paramsMap["filter[host]"] = *params.Host
		}

		if params.IDs != nil {
			ids := []string{}
			for _, id := range params.IDs {
				ids = append(ids, strconv.Itoa(id))
			}
			paramsMap["filter[ids]"] = strings.Join(ids, ",")
		}

		if params.Language != nil {
			paramsMap["filter[language]"] = *params.Language
		}
	}

	return paramsMap
}
