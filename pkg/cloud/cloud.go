package cloud

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

type Cloud struct {
	client http.Client
	config config.Config
	logger logging.Logger
	pager  http.Pager
}

func New(config config.Config) Cloud {

	client := http.NewClient(config)
	client.SetAuthStrategy(&http.PersonalAPIKeyCapableV2Authorizer{})

	pkg := Cloud{
		client: client,
		config: config,
		logger: config.GetLogger(),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}

// CloudAuthenticateIntegration authenticates a cloud provider integration and returns
// an auth reference ID that can be used with CloudLinkAccount.
func (a *Cloud) CloudAuthenticateIntegration(
	accountID int,
	providerSlug string,
	authType string,
	payload string,
) (*CloudAuthenticateIntegrationPayload, error) {
	return a.CloudAuthenticateIntegrationWithContext(context.Background(), accountID, providerSlug, authType, payload)
}

// CloudAuthenticateIntegrationWithContext authenticates a cloud provider integration (context-aware).
// providerSlug should be "GCP", authType should be "WIF", and payload should be the JSON-encoded
// WIF credential (audience + service account email).
func (a *Cloud) CloudAuthenticateIntegrationWithContext(
	ctx context.Context,
	accountID int,
	providerSlug string,
	authType string,
	payload string,
) (*CloudAuthenticateIntegrationPayload, error) {
	resp := struct {
		CloudAuthenticateIntegration CloudAuthenticateIntegrationPayload `json:"cloudAuthenticateIntegration"`
	}{}
	vars := map[string]interface{}{
		"accountId":    accountID,
		"providerSlug": providerSlug,
		"authType":     authType,
		"payload":      payload,
	}
	if err := a.client.NerdGraphQueryWithContext(ctx, cloudAuthenticateIntegrationMutation, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.CloudAuthenticateIntegration, nil
}

const cloudAuthenticateIntegrationMutation = `mutation(
	$accountId: Int!,
	$providerSlug: CloudProviderType!,
	$authType: AuthenticationType!,
	$payload: String!,
) {
	cloudAuthenticateIntegration(
		accountId: $accountId
		providerSlug: $providerSlug
		authType: $authType
		payload: $payload
	) {
		authReferenceId
	}
}`
