package alerts

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func Example_policy() {
	// Initialize the client configuration. A Personal API key is required to
	// communicate with the backend API.
	// is deprecated.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Create a new alert policy.
	p := Policy{
		Name:               "Example alert policy",
		IncidentPreference: IncidentPreferenceTypes.PerCondition,
	}

	policy, err := client.CreatePolicy(p)
	if err != nil {
		log.Fatal("error creating policy:", err)
	}

	// Create a new alert notification channel.
	ec := Channel{
		Name: "Example email notification channel",
		Type: ChannelTypes.Email,
		Configuration: ChannelConfiguration{
			Recipients:            "test@newrelic.com",
			IncludeJSONAttachment: "1",
		},
	}

	emailChannel, err := client.CreateChannel(ec)
	if err != nil {
		log.Fatal("error creating notification channel:", err)
	}

	// Associate the new alert channel with the created policy.
	_, err = client.UpdatePolicyChannels(policy.ID, []int{emailChannel.ID})
	if err != nil {
		log.Fatal("error associating policy with channel:", err)
	}

	// Create a new NRQL alert condition.
	nc := &NrqlCondition{
		Name:       "Example NRQL condition",
		Type:       "static",
		RunbookURL: "https://www.example.com/myrunbook",
		Enabled:    true,
		Nrql: NrqlQuery{
			Query:      "FROM Transaction SELECT average(duration) WHERE appName = 'Example Application'",
			SinceValue: "3",
		},
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     OperatorTypes.Above,
				Priority:     PriorityTypes.Warning,
				Threshold:    3,
				TimeFunction: TimeFunctionTypes.All,
			},
			{
				Duration:     5,
				Operator:     OperatorTypes.Above,
				Priority:     PriorityTypes.Critical,
				Threshold:    1,
				TimeFunction: TimeFunctionTypes.All,
			},
		},
	}

	_, err = client.CreateNrqlCondition(policy.ID, *nc)
	if err != nil {
		log.Fatal("error creating NRQL condition:", err)
	}
}
