//go:build unit
// +build unit

package job

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testGetJobResponseJSON = `{
		"data": {
			"customerAdministration": {
				"jobs": {
					"organizationCreateAsyncResults": {
				  		"items": [
							{
					  			"customer": {
									"customerId": "` + customerId + `"
								},
					  			"job": {
									"errorMessage": null,
									"id": "` + jobId + `",
									"status": "SUCCEEDED"
								},
								"organization": {
									"id": "` + org1Id + `",
									"name": "` + org1Name + `"
					  			}
							}
						]
					}
			  	}
			}
		}
	}`
)

func TestUnitGetJobs(t *testing.T) {
	t.Parallel()

	job := newMockResponse(t, testGetJobResponseJSON, http.StatusCreated)

	expected := &OrganizationOrganizationCreateAsyncResultCollection{
		[]OrganizationOrganizationCreateAsyncResult{
			{
				Customer: OrganizationOrganizationCreateAsyncCustomerResult{
					CustomerId: customerId,
				},
				Job: OrganizationOrganizationCreateAsyncJobResult{
					ErrorMessage: "",
					ID:           jobId,
					Status:       "SUCCEEDED",
				},
				Organization: OrganizationOrganizationCreateAsyncOrganizationResult{
					ID:   org1Id,
					Name: org1Name,
				},
			},
		},
		"",
	}

	actual, err := job.JobGetOrganizationCreateResults(jobId)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
