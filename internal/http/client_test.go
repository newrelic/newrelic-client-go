// +build unit

package http

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testBaseURL := "https://www.mocky.io"
	testTimeout := time.Second * 5
	testTransport := http.DefaultTransport

	c := NewClient(config.Config{
		APIKey:        testAPIKey,
		BaseURL:       testBaseURL,
		UserAgent:     testUserAgent,
		Timeout:       &testTimeout,
		HTTPTransport: &testTransport,
	})

	assert.Equal(t, &testTimeout, c.Config.Timeout)
	assert.Equal(t, testBaseURL, c.Config.BaseURL)
	assert.Same(t, &testTransport, c.Config.HTTPTransport)
}

func TestConfigDefaults(t *testing.T) {
	c := NewClient(config.Config{
		APIKey: testAPIKey,
	})

	assert.Equal(t, defaultBaseURLs[c.Config.Region], c.Config.BaseURL)
	assert.Contains(t, c.Config.UserAgent, "newrelic/newrelic-client-go/")
}

func TestDefaultErrorValue(t *testing.T) {
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"title":"error message"}}`))
	}))

	_, err := c.Get("/path", nil, nil, nil)

	assert.Equal(t, err.(*ErrorUnexpectedStatusCode).err, "error message")
}

type CustomErrorResponse struct {
	CustomError string `json:"custom"`
}

func (c *CustomErrorResponse) Error() string {
	return c.CustomError
}

func TestCustomErrorValue(t *testing.T) {
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"custom":"error message"}`))
	}))

	c.SetErrorValue(&CustomErrorResponse{})

	_, err := c.Get("/path", nil, nil, nil)

	assert.Equal(t, err.(*ErrorUnexpectedStatusCode).err, "error message")
}

type CustomResponseValue struct {
	Custom string `json:"custom"`
}

func TestResponseValue(t *testing.T) {
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"custom":"custom response string"}`))
	}))

	v := &CustomResponseValue{}
	_, err := c.Get("/path", nil, nil, v)

	assert.NoError(t, err)
	assert.Equal(t, &CustomResponseValue{Custom: "custom response string"}, v)
}

func TestQueryParams(t *testing.T) {
	queryParams := map[string]string{
		"a": "1",
		"b": "2",
	}

	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		a, ok := r.URL.Query()["a"]
		assert.True(t, ok)
		assert.Equal(t, "1", a[0])

		b, ok := r.URL.Query()["b"]
		assert.True(t, ok)
		assert.Equal(t, "2", b[0])
	}))

	_, _ = c.Get("/path", &queryParams, nil, nil)
}

type TestRequestBody struct {
	A string `json:"a"`
	B string `json:"b"`
}

func TestRequestBodyMarshal(t *testing.T) {
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

	_, _ = c.Get("/path", nil, expected, nil)
}

type TestInvalidRequestBody struct {
	Channel chan int `json:"a"`
}

func TestRequestBodyMarshalError(t *testing.T) {
	b := TestInvalidRequestBody{
		Channel: make(chan int),
	}

	c := NewTestAPIClient(nil)

	_, err := c.Get("/path", nil, b, nil)
	assert.Error(t, err)
}

func TestUrlParseError(t *testing.T) {
	c := NewTestAPIClient(nil)

	_, err := c.Get("\\", nil, nil, nil)
	assert.Error(t, err)
}

func TestPathOnlyUrl(t *testing.T) {
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, r.URL, "https://www.mocky.io/v2/path")
	}))

	c.Config.BaseURL = "https://www.mocky.io/v2"

	_, _ = c.Get("/path", nil, nil, nil)
}

func TestHostAndPathUrl(t *testing.T) {
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, r.URL, "https:/www.httpbin.org/path")
	}))

	c.Config.BaseURL = "https://www.mocky.io/v2"

	_, _ = c.Get("https://www.httpbin.org/path", nil, nil, nil)
}

type TestInvalidReponseBody struct {
	Channel chan int `json:"channel"`
}

func TestResponseUnmarshalError(t *testing.T) {
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"channel": "test"}`))
	}))

	_, err := c.Get("/path", nil, nil, &TestInvalidReponseBody{})

	assert.Error(t, err)
}

func TestHeaders(t *testing.T) {
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, testAPIKey, r.Header.Get("x-api-key"))
		assert.Equal(t, testUserAgent, r.Header.Get("user-agent"))
	}))

	_, err := c.Get("/path", nil, nil, nil)

	assert.Nil(t, err)
}

func TestErrNotFound(t *testing.T) {
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	_, err := c.Get("/path", nil, nil, nil)

	assert.IsType(t, &ErrorNotFound{}, err)
}

func TestInternalServerError(t *testing.T) {
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	_, err := c.Get("/path", nil, nil, nil)

	assert.IsType(t, &ErrorUnexpectedStatusCode{}, err)
}
