//go:build integration
// +build integration

package newrelic

import (
	"context"
	"fmt"
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
	time.Sleep(3 * time.Second)
	fmt.Println("proxy", proxy)
	srv := &http.Server{
		Addr:    "localhost:1337",
		Handler: proxy,
	}
	fmt.Println("srv", srv)
	defer srv.Shutdown(context.Background())

	client, err := New(
		ConfigPersonalAPIKey(os.Getenv("NEW_RELIC_API_KEY")),
		ConfigRegion("US"),
		ConfigLogLevel("DEBUG"),
	)

	fmt.Println("client", client)
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
