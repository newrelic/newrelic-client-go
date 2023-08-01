package changetracking

import (
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

/*
Since we can't yet generate the CustomAttributes type with Tutone,
defer ChangeTrackingDeploymentInput placement into types.go until CustomAttributes GA.
*/

// ChangeTrackingDeploymentInput - A deployment.
type ChangeTrackingDeploymentInput struct {
	// A URL for the changelog or list of changes if not linkable.
	Changelog string `json:"changelog,omitempty"`
	// The commit identifier, for example, a Git commit SHA.
	Commit string `json:"commit,omitempty"`
	// A list of key:value attribute pairs
	CustomAttributes *map[string]string `json:"customAttributes,omitempty"`
	// A link back to the system generating the deployment.
	DeepLink string `json:"deepLink,omitempty"`
	// The type of deployment, for example, ‘Blue green’ or ‘Rolling’.
	DeploymentType ChangeTrackingDeploymentType `json:"deploymentType,omitempty"`
	// A description of the deployment.
	Description string `json:"description,omitempty"`
	// The NR1 entity that was deployed.
	EntityGUID common.EntityGUID `json:"entityGuid"`
	// String that can be used to correlate two or more events.
	GroupId string `json:"groupId,omitempty"`
	// The start time of the deployment, the number of milliseconds since the Unix epoch.  Defaults to now
	Timestamp nrtime.EpochMilliseconds `json:"timestamp,omitempty"`
	// Username of the deployer or bot.
	User string `json:"user,omitempty"`
	// The version of the deployed software, for example, something like v1.1
	Version string `json:"version"`
}
