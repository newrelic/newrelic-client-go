package entities

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-client-go/pkg/accounts"
)

// Need Outlines to also implement Entity
func (x AlertableEntityOutline) ImplementsEntity()                       {}
func (x ApmApplicationEntityOutline) ImplementsEntity()                  {}
func (x ApmBrowserApplicationEntityOutline) ImplementsEntity()           {}
func (x ApmDatabaseInstanceEntityOutline) ImplementsEntity()             {}
func (x ApmExternalServiceEntityOutline) ImplementsEntity()              {}
func (x BrowserApplicationEntityOutline) ImplementsEntity()              {}
func (x DashboardEntityOutline) ImplementsEntity()                       {}
func (x EntityOutline) ImplementsEntity()                                {}
func (x GenericEntityOutline) ImplementsEntity()                         {}
func (x GenericInfrastructureEntityOutline) ImplementsEntity()           {}
func (x InfrastructureAwsLambdaFunctionEntityOutline) ImplementsEntity() {}
func (x InfrastructureHostEntityOutline) ImplementsEntity()              {}
func (x InfrastructureIntegrationEntityOutline) ImplementsEntity()       {}
func (x MobileApplicationEntityOutline) ImplementsEntity()               {}
func (x SecureCredentialEntityOutline) ImplementsEntity()                {}
func (x SyntheticMonitorEntityOutline) ImplementsEntity()                {}
func (x ThirdPartyServiceEntityOutline) ImplementsEntity()               {}
func (x UnavailableEntityOutline) ImplementsEntity()                     {}
func (x WorkloadEntityOutline) ImplementsEntity()                        {}

func (a *Actor) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	for k, v := range objMap {
		if v == nil {
			continue
		}

		switch k {
		case "accounts":
			var accts []accounts.AccountOutline
			err = json.Unmarshal(*v, &accts)
			if err != nil {
				return err
			}

			a.Accounts = accts
		case "entity":
			var e *EntityInterface
			e, err = UnmarshalEntityInterface([]byte(*v))
			if err != nil {
				return err
			}

			a.Entity = *e
		case "entities":
			var rawEntities []*json.RawMessage
			err = json.Unmarshal(*v, &rawEntities)
			if err != nil {
				return err
			}

			for _, m := range rawEntities {
				var e *EntityInterface
				e, err = UnmarshalEntityInterface(*m)
				if err != nil {
					return err
				}

				if e != nil {
					a.Entities = append(a.Entities, *e)
				}
			}
		case "entitySearch":
			var es EntitySearch
			err = json.Unmarshal(*v, &es)
			if err != nil {
				return err
			}

			a.EntitySearch = es
		default:
			log.Errorf("Unknown key '%s' value: %s", k, string(*v))
		}
	}

	return nil
}
