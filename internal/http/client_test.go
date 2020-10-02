// +build unit

package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

const (
	testServiceName = "serviceName"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	testRestURL := "https://www.mocky.io"
	testTimeout := time.Second * 5
	testTransport := http.DefaultTransport

	tc := config.New()
	tc.HTTPTransport = testTransport
	tc.Region().SetRestBaseURL(testRestURL)
	tc.ServiceName = testServiceName
	tc.Timeout = &testTimeout
	tc.UserAgent = mock.UserAgent

	c := NewClient(tc)

	assert.Equal(t, &testTimeout, c.config.Timeout)
	assert.Equal(t, testRestURL, c.config.Region().RestURL())
	assert.Equal(t, mock.UserAgent, c.config.UserAgent)
	assert.Equal(t, c.config.ServiceName, testServiceName+"|newrelic-client-go")

	assert.Same(t, testTransport, c.config.HTTPTransport)
}

func TestConfigDefaults(t *testing.T) {
	t.Parallel()
	tc := mock.NewTestConfig(t, nil)
	tc.ServiceName = testServiceName

	c := NewClient(tc)

	assert.Contains(t, c.config.UserAgent, "newrelic/newrelic-client-go")
	assert.Equal(t, c.config.ServiceName, testServiceName+"|newrelic-client-go")
}

func TestDefaultErrorValue(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"title":"error message"}}`))
	}))

	_, err := c.Get(c.config.Region().RestURL("path"), nil, nil)

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

func (c *CustomErrorResponse) IsNotFound() bool {
	return false
}

func TestCustomErrorValue(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"custom":"error message"}`))
	}))

	c.SetErrorValue(&CustomErrorResponse{})

	_, err := c.Get(c.config.Region().RestURL("path"), nil, nil)

	assert.Contains(t, err.(*errors.UnexpectedStatusCode).Error(), "error message")
}

type CustomResponseValue struct {
	Custom string `json:"custom"`
}

func TestResponseValue(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"custom":"custom response string"}`))
	}))

	v := &CustomResponseValue{}
	_, err := c.Get(c.config.Region().RestURL("path"), nil, v)

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

	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		a := r.URL.Query().Get("a")
		assert.Equal(t, "1", a)

		b := r.URL.Query().Get("b")
		assert.Equal(t, "2", b)
	}))

	_, _ = c.Get(c.config.Region().RestURL("path"), &queryParams, nil)
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

	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		actual := &TestRequestBody{}
		err := json.NewDecoder(r.Body).Decode(&actual)

		assert.NoError(t, err)
		assert.Equal(t, &expected, actual)
	}))

	_, _ = c.Post(c.config.Region().RestURL("path"), nil, expected, nil)
}

type TestInvalidRequestBody struct {
	Channel chan int `json:"a"`
}

func TestRequestBodyMarshalError(t *testing.T) {
	t.Parallel()
	b := TestInvalidRequestBody{
		Channel: make(chan int),
	}

	c := NewTestAPIClient(t, nil)

	_, err := c.Post(c.config.Region().RestURL("/path"), nil, b, nil)
	assert.Error(t, err)
}

func TestUrlParseError(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, nil)

	_, err := c.Get(c.config.Region().RestURL("\\"), nil, nil)
	assert.Error(t, err)
}

type TestInvalidReponseBody struct {
	Channel chan int `json:"channel"`
}

func TestResponseUnmarshalError(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"channel": "test"}`))
	}))

	_, err := c.Get(c.config.Region().RestURL("path"), nil, &TestInvalidReponseBody{})

	assert.Error(t, err)
}

func TestHeaders(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, mock.UserAgent, r.Header.Get("user-agent"))
		assert.Equal(t, "newrelic-client-go", r.Header.Get("newrelic-requesting-services"))
	}))

	_, err := c.Get(c.config.Region().RestURL("path"), nil, nil)

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

	tc := mock.NewTestConfig(t, ts)
	tc.UserAgent = "custom-user-agent"
	tc.ServiceName = "custom-requesting-service"

	c := NewClient(tc)

	_, err := c.Get(c.config.Region().RestURL("path"), nil, nil)

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

	tc := mock.NewTestConfig(t, ts)

	c := NewClient(tc)

	req, err := c.NewRequest("GET", c.config.Region().RestURL("path"), nil, nil, nil)

	req.SetHeader("user-agent", "custom-user-agent")
	req.SetServiceName("custom-requesting-service")

	_, err = c.Do(req)

	assert.Nil(t, err)
}

func TestErrNotFound(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	_, err := c.Get(c.config.Region().RestURL("path"), nil, nil)

	assert.IsType(t, &errors.NotFound{}, err)
}

func TestInternalServerError(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	_, err := c.Get(c.config.Region().RestURL("path"), nil, nil)

	assert.IsType(t, &errors.UnexpectedStatusCode{}, err)
}

func TestPost(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	// string
	_, err := c.Post(c.config.Region().RestURL("path"), &struct{}{}, "test string payload", &struct{}{})
	assert.NoError(t, err)

	// []byte
	_, err = c.Post(c.config.Region().RestURL("path"), &struct{}{}, []byte(`bytes`), &struct{}{})
	assert.NoError(t, err)

	// other data type
	_, err = c.Post(c.config.Region().RestURL("path"), &struct{}{}, &struct{}{}, &struct{}{})
	assert.NoError(t, err)
}

func TestPut(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	_, err := c.Put(c.config.Region().RestURL("path"), &struct{}{}, &struct{}{}, &struct{}{})

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{}`))
	}))

	_, err := c.Delete(c.config.Region().RestURL("path"), &struct{}{}, &struct{}{})

	assert.NoError(t, err)
}
