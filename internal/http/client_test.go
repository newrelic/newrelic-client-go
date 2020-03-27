// +build unit

package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

const (
	testServiceName = "serviceName"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	testBaseURL := "https://www.mocky.io"
	testTimeout := time.Second * 5
	testTransport := http.DefaultTransport

	c := NewClient(config.Config{
		PersonalAPIKey: testPersonalAPIKey,
		BaseURL:        testBaseURL,
		HTTPTransport:  testTransport,
		ServiceName:    testServiceName,
		Timeout:        &testTimeout,
		UserAgent:      testUserAgent,
	})

	assert.Equal(t, &testTimeout, c.config.Timeout)
	assert.Equal(t, testBaseURL, c.config.BaseURL)
	assert.Equal(t, testUserAgent, c.config.UserAgent)
	assert.Equal(t, c.config.ServiceName, testServiceName+"|newrelic-client-go")

	assert.Same(t, testTransport, c.config.HTTPTransport)
}

func TestConfigDefaults(t *testing.T) {
	t.Parallel()
	c := NewClient(config.Config{
		PersonalAPIKey: testPersonalAPIKey,
	})

	assert.Equal(t, region.DefaultBaseURLs[region.Parse(c.config.Region.String())], c.config.BaseURL)
	assert.Contains(t, c.config.UserAgent, "newrelic/newrelic-client-go/")
	assert.Equal(t, c.config.ServiceName, "newrelic-client-go")
}

func TestDefaultErrorValue(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"title":"error message"}}`))
	}))

	_, err := c.Get("/path", nil, nil)

	assert.Contains(t, err.(*errors.UnexpectedStatusCode).Error(), "error message")
}

type CustomErrorResponse struct {
	CustomError string `json:"custom"`
}

func (c *CustomErrorResponse) New() ErrorResponse {
	return &CustomErrorResponse{}
}

func (c *CustomErrorResponse) Error() string {
	return c.CustomError
}

func TestCustomErrorValue(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"custom":"error message"}`))
	}))

	c.SetErrorValue(&CustomErrorResponse{})

	_, err := c.Get("/path", nil, nil)

	assert.Contains(t, err.(*errors.UnexpectedStatusCode).Error(), "error message")
}

type CustomResponseValue struct {
	Custom string `json:"custom"`
}

func TestResponseValue(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"custom":"custom response string"}`))
	}))

	v := &CustomResponseValue{}
	_, err := c.Get("/path", nil, v)

	assert.NoError(t, err)
	assert.Equal(t, &CustomResponseValue{Custom: "custom response string"}, v)
}

func TestQueryParams(t *testing.T) {
	t.Parallel()
	queryParams := struct {
		A int `url:"a,omitempty"`
		B int `url:"b,omitempty"`
	}{
		A: 1,
		B: 2,
	}

	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		a := r.URL.Query().Get("a")
		assert.Equal(t, "1", a)

		b := r.URL.Query().Get("b")
		assert.Equal(t, "2", b)
	}))

	_, _ = c.Get("/path", &queryParams, nil)
}

type TestRequestBody struct {
	A string `json:"a"`
	B string `json:"b"`
}

func TestRequestBodyMarshal(t *testing.T) {
	t.Parallel()
	expected := TestRequestBody{
		A: "1",
		B: "2",
	}

	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		actual := &TestRequestBody{}
		err := json.NewDecoder(r.Body).Decode(&actual)

		assert.NoError(t, err)
		assert.Equal(t, &expected, actual)
	}))

	_, _ = c.Post("/path", nil, expected, nil)
}

type TestInvalidRequestBody struct {
	Channel chan int `json:"a"`
}

func TestRequestBodyMarshalError(t *testing.T) {
	t.Parallel()
	b := TestInvalidRequestBody{
		Channel: make(chan int),
	}

	c := NewTestAPIClient(nil)

	_, err := c.Post("/path", nil, b, nil)
	assert.Error(t, err)
}

func TestUrlParseError(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(nil)

	_, err := c.Get("\\", nil, nil)
	assert.Error(t, err)
}

func TestPathOnlyUrl(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, r.URL, "https://www.mocky.io/v2/path")
	}))

	c.config.BaseURL = "https://www.mocky.io/v2"

	_, _ = c.Get("/path", nil, nil)
}

func TestHostAndPathUrl(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, r.URL, "https:/www.httpbin.org/path")
	}))

	c.config.BaseURL = "https://www.mocky.io/v2"

	_, _ = c.Get("https://www.httpbin.org/path", nil, nil)
}

type TestInvalidReponseBody struct {
	Channel chan int `json:"channel"`
}

func TestResponseUnmarshalError(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"channel": "test"}`))
	}))

	_, err := c.Get("/path", nil, &TestInvalidReponseBody{})

	assert.Error(t, err)
}

func TestHeaders(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, testUserAgent, r.Header.Get("user-agent"))
		assert.Equal(t, "newrelic-client-go", r.Header.Get("newrelic-requesting-services"))
	}))

	_, err := c.Get("/path", nil, nil)

	assert.Nil(t, err)
}

func TestCustomClientHeaders(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, "custom-user-agent", r.Header.Get("user-agent"))
		assert.Equal(t, "custom-requesting-service|newrelic-client-go", r.Header.Get("newrelic-requesting-services"))
	}))

	c := NewClient(config.Config{
		PersonalAPIKey: testPersonalAPIKey,
		AdminAPIKey:    testAdminAPIKey,
		BaseURL:        ts.URL,
		UserAgent:      "custom-user-agent",
		ServiceName:    "custom-requesting-service",
	})

	_, err := c.Get("/path", nil, nil)

	assert.Nil(t, err)
}

func TestCustomRequestHeaders(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, "custom-user-agent", r.Header.Get("user-agent"))
		assert.Equal(t, "custom-requesting-service|newrelic-client-go", r.Header.Get("newrelic-requesting-services"))
	}))

	c := NewClient(config.Config{
		PersonalAPIKey: testPersonalAPIKey,
		AdminAPIKey:    testAdminAPIKey,
		BaseURL:        ts.URL,
	})

	req, err := NewRequest(c, "GET", "/path", nil, nil, nil)

	req.SetHeader("user-agent", "custom-user-agent")
	req.SetServiceName("custom-requesting-service")

	_, err = c.Do(req)

	assert.Nil(t, err)
}

func TestAdminAPIKeyHeader(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, testAdminAPIKey, r.Header.Get("x-api-key"))
	}))

	c := NewClient(config.Config{
		AdminAPIKey: testAdminAPIKey,
		BaseURL:     ts.URL,
		UserAgent:   testUserAgent,
	})

	_, err := c.Get("/path", nil, nil)

	assert.Nil(t, err)
}

func TestErrNotFound(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	_, err := c.Get("/path", nil, nil)

	assert.IsType(t, &errors.NotFound{}, err)
}

func TestInternalServerError(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	_, err := c.Get("/path", nil, nil)

	assert.IsType(t, &errors.UnexpectedStatusCode{}, err)
}

func TestPost(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	_, err := c.Post("/path", &struct{}{}, &struct{}{}, &struct{}{})

	assert.NoError(t, err)
}

func TestRawPost(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	// string
	_, err := c.RawPost("/path", &struct{}{}, "test string payload", &struct{}{})
	assert.NoError(t, err)

	// []byte
	_, err = c.RawPost("/path", &struct{}{}, []byte(`bytes`), &struct{}{})
	assert.NoError(t, err)

	// invalid
	_, err = c.RawPost("/path", &struct{}{}, &struct{}{}, &struct{}{})
	assert.Error(t, err)
}

func TestPut(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	_, err := c.Put("/path", &struct{}{}, &struct{}{}, &struct{}{})

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{}`))
	}))

	_, err := c.Delete("/path", &struct{}{}, &struct{}{})

	assert.NoError(t, err)
}
