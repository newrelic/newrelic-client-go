# workflows
--
    import "github.com/newrelic/newrelic-client-go/v3/pkg/worklfows"

## Usage

#### type Workflows

```go
type Workflows struct {
}
```

Workflows is used to communicate with New Relic Workflows.

#### func  New

```go
func New(config config.Config) Workflows
```
New is used to create a new Workflows' client instance.

### Workflows

#### func  TestMutationWorkflow

```go
func TestMutationWorkflow(t *testing.T)
```

#### type AiWorkflowsWorkflow

```go
// AiWorkflowsWorkflow - Workflow object
type AiWorkflowsWorkflow struct {
AccountID                   int                                     `json:"accountId"`
CreatedAt                   nrtime.DateTime                         `json:"createdAt"`
DestinationConfigurations   []AiWorkflowsDestinationConfiguration   `json:"destinationConfigurations"`
DestinationsEnabled         bool                                    `json:"destinationsEnabled"`
Enrichments                 []AiWorkflowsEnrichment                 `json:"enrichments"`
EnrichmentsEnabled          bool                                    `json:"enrichmentsEnabled"`
ID                          string                                  `json:"id"`
IssuesFilter                AiWorkflowsFilter                       `json:"issuesFilter"`
LastRun                     nrtime.DateTime                         `json:"lastRun,omitempty"`
MutingRulesHandling         AiWorkflowsMutingRulesHandling          `json:"mutingRulesHandling"`
Name                        string                                  `json:"name"`
UpdatedAt                   nrtime.DateTime                         `json:"updatedAt"`
WorkflowEnabled             bool                                    `json:"workflowEnabled"`
}
```

AiWorkflowsWorkflow represents a New Relic AI workflow.

#### func (*Workflows) AiWorkflowsCreateWorkflow

```go
func (a *Workflows) AiWorkflowsCreateWorkflow(accountID int,createWorkflowData AiWorkflowsCreateWorkflowInput) (*AiWorkflowsCreateWorkflowResponse, error)
```
AiWorkflowsCreateWorkflow creates a new workflow for a given account.

#### func (*Workflows) GetWorkflows

```go
func (a *Workflows) GetWorkflows(accountID int,cursor string, filters ai.AiWorkflowsFilters) (*AiWorkflowsWorkflows, error)
```
GetWorkflows returns a list of workflow for a given account. You can filter by ID.

#### func (*Workflows) AiWorkflowsUpdateWorkflow

```go
func (a *Workflows) AiWorkflowsUpdateWorkflow(accountID int,updateWorkflowData AiWorkflowsUpdateWorkflowInput) (*AiWorkflowsUpdateWorkflowResponse, error)
```
AiWorkflowsUpdateWorkflow update a workflow for a given account.

#### type AiWorkflowsDeleteWorkflow

```go
func (a *Workflows) AiWorkflowsDeleteWorkflow(accountID int, iD string) (*AiWorkflowsDeleteWorkflowResponse, error)
```

AiWorkflowsDeleteWorkflow delete a workflow for a given account.
