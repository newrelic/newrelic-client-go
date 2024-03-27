package jobs

func (a *Jobs) GetOrganizationCreateResults(jobId string) (*OrganizationOrganizationCreateAsyncResultCollection, error) {
	return a.GetOrganizationCreateAsyncResults("", OrganizationOrganizationCreateAsyncResultFilterInput{
		JobId: OrganizationOrganizationCreateJobIdInput{
			Eq: jobId,
		},
	})
}
