package notifications

import (
	"time"
)

type Destination struct {
	ID         UUID              `json:"id,omitempty"`
	Name       string            `json:"name,omitempty"`
	Type       DestinationType   `json:"type,omitempty"`
	CreatedAt  time.Time         `json:"createdAt,omitempty"`
	UpdatedAt  time.Time         `json:"updatedAt,omitempty"`
	UpdatedBy  int               `json:"updatedBy,omitempty"`
	AccountId  int               `json:"accountId,omitempty"`
	Status     DestinationStatus `json:"status,omitempty"`
	Active     bool              `json:"active,omitempty"`
	LastSent   time.Time         `json:"lastSent"`
	Auth       Auth              `json:"auth"`
	Properties []Property        `json:"properties,omitempty"`
}

type UUID string

type DestinationType string

type DestinationStatus string

var (
	DestinationTypes = struct {
		Email                       DestinationType
		EventBridge                 DestinationType
		Jira                        DestinationType
		ServiceNow                  DestinationType
		Webhook                     DestinationType
		Slack                       DestinationType
		PagerDutyServiceIntegration DestinationType
		PagerDutyAccountIntegration DestinationType
	}{
		Email:                       "EMAIL",
		EventBridge:                 "EVENT_BRIDGE",
		Jira:                        "JIRA",
		ServiceNow:                  "SERVICE_NOW",
		Webhook:                     "WEBHOOK",
		Slack:                       "SLACK",
		PagerDutyServiceIntegration: "PAGERDUTY_SERVICE_INTEGRATION",
		PagerDutyAccountIntegration: "PAGERDUTY_ACCOUNT_INTEGRATION",
	}
)

var (
	DestinationStatuses = struct {
		Draft                DestinationStatus
		Tested               DestinationStatus
		Throttled            DestinationStatus
		Error                DestinationStatus
		Default              DestinationStatus
		AuthenticationError  DestinationStatus
		AuthorizationError   DestinationStatus
		ConfigurationError   DestinationStatus
		ThrottlingWarning    DestinationStatus
		AuthorizationWarning DestinationStatus
		TemporaryWarning     DestinationStatus
		UnknownError         DestinationStatus
	}{
		Draft:                "draft",
		Tested:               "tested",
		Throttled:            "throttled",
		Error:                "error",
		Default:              "default",
		AuthenticationError:  "authentication_error",
		AuthorizationError:   "authorization_error",
		ConfigurationError:   "configuration_error",
		ThrottlingWarning:    "throttling_error",
		AuthorizationWarning: "authorization_warning",
		TemporaryWarning:     "temporary_warning",
		UnknownError:         "unknown_error",
	}
)

// SecureValue specifies a secure value, ie a password, an API key, etc
type SecureValue string

type BasicAuth struct {
	User     string      `json:"user"`
	Password SecureValue `json:"password"`
}

type TokenAuth struct {
	Prefix string      `json:"prefix"`
	Token  SecureValue `json:"token"`
}

type Auth struct {
	AuthType *AuthType `json:"authType,omitempty"`

	// Prefix exists ONLY for auth of type TOKEN.
	Prefix *string `json:"prefix"`

	// User exists ONLY for auth of type BASIC.
	User *string `json:"user"`
}

type AuthType string

var (
	AuthTypes = struct {
		Basic  AuthType
		OAuth2 AuthType
		Token  AuthType
	}{
		Basic:  "BASIC",
		OAuth2: "OAUTH2",
		Token:  "TOKEN",
	}
)

type Property struct {
	Key          string `json:"key,omitempty"`
	Label        string `json:"label,omitempty"`
	Value        string `json:"value,omitempty"`
	DisplayValue string `json:"display_value,omitempty"`
}

type PropertyInput struct {
	DisplayValue string `json:"display_value"`
	Key          string `json:"key,omitempty"`
	Label        string `json:"label"`
	Value        string `json:"value,omitempty"`
}

type PropertyFilter struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type AiNotificationsCredentialsInput struct {
	Basic BasicAuth `json:"basic,omitempty"`
	Token TokenAuth `json:"token,omitempty"`
	Type  AuthType  `json:"type,omitempty"`
}

type DestinationInput struct {
	Name       string                          `json:"name,omitempty"`
	Type       DestinationType                 `json:"type,omitempty"`
	Properties []PropertyInput                 `json:"properties,omitempty"`
	Auth       AiNotificationsCredentialsInput `json:"auth,omitempty"`
}

type AiNotificationsDestinationFilter struct {
	Id        UUID            `json:"id"`
	Active    bool            `json:"active"`
	AuthType  AuthType        `json:"auth_type"`
	Name      string          `json:"name"`
	Property  PropertyFilter  `json:"properties"`
	Type      DestinationType `json:"type"`
	UpdatedAt *time.Time      `json:"updated_at"`
}

type createDestinationResponse struct {
	AiNotificationsCreateDestination struct {
		Destination Destination `json:"destination,omitempty"`
	} `json:"aiNotificationsCreateDestination"`
}

type listDestinationsResponse struct {
	Actor struct {
		Account struct {
			AiNotifications struct {
				Destinations struct {
					Entities []*Destination `json:"entities"`
				}
			} `json:"aiNotifications"`
		} `json:"account"`
	} `json:"actor"`
}

type deleteDestinationResponse struct {
	AiNotificationsDeleteDestination struct {
		Ids []string `json:"ids,omitempty"`
	} `json:"aiNotificationsDeleteDestination"`
}
