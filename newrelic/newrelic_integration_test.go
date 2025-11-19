//go:build integration
// +build integration

package newrelic

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/elazarl/goproxy"

	"github.com/stretchr/testify/require"
)

func TestNewRelic_TestEndpoints_withProxy(t *testing.T) {
	err := os.Setenv("HTTPS_PROXY", "localhost:1337")
	require.NoError(t, err)
	defer os.Unsetenv("HTTPS_PROXY")

	proxy := goproxy.NewProxyHttpServer()
	srv := &http.Server{
		Addr:    "localhost:1337",
		Handler: proxy,
	}
	defer srv.Shutdown(context.Background())

	client, err := New(
		ConfigPersonalAPIKey(os.Getenv("NEW_RELIC_API_KEY")),
		ConfigRegion("US"),
		ConfigLogLevel("DEBUG"),
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	err = client.TestEndpoints()
	require.Error(t, err)

	go func() {
		srv.ListenAndServe()
	}()

	err = client.TestEndpoints()

	require.NoError(t, err)
}
