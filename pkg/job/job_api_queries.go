package job

func (a *Job) JobGetOrganizationCreateResults(jobId string) (*OrganizationOrganizationCreateAsyncResultCollection, error) {
	return a.GetOrganizationCreateAsyncResults("", OrganizationOrganizationCreateAsyncResultFilterInput{
		JobID: &JobJobIDInput{
			Eq: jobId,
		},
	})
}
