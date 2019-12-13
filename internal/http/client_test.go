// +build unit

package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

const testAPIKey string = "12345"
const testUserAgentHeader string = "go-newrelic/test"

func TestClientHeaders(t *testing.T) {
	cli := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != testAPIKey {
			t.Fatal("x-api-key was not correctly set")
		}
		if r.Header.Get("user-agent") != testUserAgentHeader {
			t.Fatal("user-agent was not correctly set")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))

	err := cli.Get("/path", nil, nil)

	if err != nil {
		t.Fatal(err)
	}
}

func TestClientDoPaging(t *testing.T) {
	for i, c := range []struct {
		expectedNext string
		linkHeader   string
		body         string
	}{
		{"", "", ""},
		{"", "", "{}"},
		{"", `<https://api.github.com/user/58276/repos?page=2>; rel="last"`, "{}"},
		//{"", "", `{"links":null}`},
		//{"", "", `{"links":{}}`},
		//{"", "", `{"links":{"last":"foo"}}`},

		{"https://api.github.com/user/58276/repos?page=2", `<https://api.github.com/user/58276/repos?page=2>; rel="next"`, "{}"},
		//{"https://api.github.com/user/58276/repos?page=2", "", `{"links":{"next":"https://api.github.com/user/58276/repos?page=2"}}`},
		//{"https://api.github.com/user/58276/repos?page=2", "", `{"links":{"next":"https://api.github.com/user/58276/repos?page=2"}}`},
		{"https://api.github.com/user/58276/repos?page=2", `<https://api.github.com/user/58276/repos?page=2>; rel="next"`, `{"links":{"next":"https://should-not-match"}}`},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expectedNext), func(t *testing.T) {
			cli := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if c.linkHeader != "" {
					w.Header().Set("Link", c.linkHeader)
				}
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(c.body))

				if err != nil {
					t.Fatal(err)
				}
			}))
			paging, err := cli.do("GET", "/path", nil)
			if err != nil {
				t.Fatal(err)
			}
			if paging.Next != c.expectedNext {
				t.Fatalf("expected %q but got %q", c.expectedNext, paging.Next)
			}
		})
	}
}

func TestErrNotFound(t *testing.T) {
	cli := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	err := cli.Get("/path", nil, nil)

	if err != ErrNotFound {
		t.Fatal(err)
	}
}

func TestInternalServerError(t *testing.T) {
	cli := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	err := cli.Get("/path", nil, nil)

	if err == nil {
		t.Fatal(err)
	}
}

func TestDefaultConfig(t *testing.T) {
	c := NewClient(config.Config{
		APIKey: testAPIKey,
	})

	expectedBaseURL := "https://api.newrelic.com/v2"
	if c.Client.HostURL != expectedBaseURL {
		t.Fatalf("expected baseURL: %s, received: %s", expectedBaseURL, c.Client.HostURL)
	}

	if c.Client.Debug {
		t.Fatalf("expected debug mode to be off")
	}
}

func TestSetProxyURL(t *testing.T) {
	expectedProxyURL := "http://proxy.url"
	c := NewClient(config.Config{
		APIKey:   testAPIKey,
		ProxyURL: expectedProxyURL,
	})

	if !c.Client.IsProxySet() {
		t.Fatalf("expected proxy to be set")
	}
}

func TestSetDebug(t *testing.T) {
	c := NewClient(config.Config{
		APIKey: testAPIKey,
		Debug:  true,
	})

	if !c.Client.Debug {
		t.Fatalf("expected debug mode to be on")
	}
}

func newTestAPIClient(handler http.Handler) *NewRelicClient {
	ts := httptest.NewServer(handler)

	c := NewClient(config.Config{
		APIKey:    testAPIKey,
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: testUserAgentHeader,
	})

	return &c
}
