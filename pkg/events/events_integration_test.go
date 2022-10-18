//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

const (
	testBatchTimeout = 10
	testBatchSize    = 1
)

// Used to create a variety of events for testing
type testEventData struct {
	Event   interface{}
	marshal []byte
	err     error
}

var testEvents = []testEventData{
	{
		Event: struct {
			EventType string  `json:"eventType"`
			Amount    float64 `json:"amount"`
		}{
			EventType: "Purchase",
			Amount:    123.45,
		},
		marshal: []byte(`{"eventType":"Purchase","amount":123.45}`),
		err:     nil,
	},
	{
		Event: struct {
			EventType string `json:"eventType"`
			Data      string `json:"data"`
		}{
			EventType: "Dessert",
			Data:      "Biscuit cotton candy candy canes dessert muffin ice cream carrot cake. Marzipan lemon drops lemon drops. Candy cake toffee. Powder pie wafer bonbon dessert powder. Jujubes sweet chocolate bar gummies jelly-o. Wafer biscuit candy. Oat cake cookie jelly liquorice cupcake cupcake dragée cupcake wafer. Wafer chocolate bar marzipan powder jujubes cake oat cake sweet roll. Chocolate bar caramels jelly sugar plum donut. Candy donut tiramisu candy canes icing macaroon. Gummies macaroon jujubes candy gummies cotton candy sesame snaps dragée. Sweet roll icing cake pie sweet. Candy topping toffee. Sesame snaps cookie lemon drops wafer jujubes powder fruitcake.",
		},
		marshal: []byte(`{"eventType":"Dessert","data":"Biscuit cotton candy candy canes dessert muffin ice cream carrot cake. Marzipan lemon drops lemon drops. Candy cake toffee. Powder pie wafer bonbon dessert powder. Jujubes sweet chocolate bar gummies jelly-o. Wafer biscuit candy. Oat cake cookie jelly liquorice cupcake cupcake dragée cupcake wafer. Wafer chocolate bar marzipan powder jujubes cake oat cake sweet roll. Chocolate bar caramels jelly sugar plum donut. Candy donut tiramisu candy canes icing macaroon. Gummies macaroon jujubes candy gummies cotton candy sesame snaps dragée. Sweet roll icing cake pie sweet. Candy topping toffee. Sesame snaps cookie lemon drops wafer jujubes powder fruitcake."}`),
		err:     nil,
	},
}

// Test: CreateEvent
func TestIntegrationEvents(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	for _, event := range testEvents {
		err := client.CreateEvent(testAccountID, event.Event)
		if event.err == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, event.err, err)
		}
	}
}

func TestIntegrationEventsLicenseKey(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	tc := mock.NewIntegrationTestConfig(t)
	tc.InsightsInsertKey = ""
	client := New(tc)

	for _, event := range testEvents {
		err := client.CreateEvent(testAccountID, event.Event)
		if event.err == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, event.err, err)
		}
	}
}

func TestIntegrationEvents_BatchMode_Timeout(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	err = client.BatchMode(context.Background(), testAccountID, BatchConfigTimeout(testBatchTimeout))
	require.NoError(t, err)

	for _, event := range testEvents {
		err := client.EnqueueEvent(context.Background(), event.Event)
		assert.NoError(t, err)
	}

	// Should of flushed
	time.Sleep(time.Duration(2*testBatchTimeout) * time.Second)
	assert.Equal(t, 0, len(client.eventQueue))
}

func TestIntegrationEvents_BatchMode_Size(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	err = client.BatchMode(context.Background(), testAccountID, BatchConfigQueueSize(testBatchSize))
	require.NoError(t, err)

	for _, event := range testEvents {
		err := client.EnqueueEvent(context.Background(), event.Event)
		assert.NoError(t, err)
	}

	// Should of flushed
	time.Sleep(time.Duration(2*testBatchTimeout) * time.Second)
	assert.Equal(t, 0, len(client.eventQueue))
}

func TestIntegrationEvents_marshalEvent(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	for _, event := range testEvents {
		data, err := client.marshalEvent(event.Event)
		if event.err == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, event.err, err)
		}

		assert.Equal(t, event.marshal, *data)
	}
}

func newIntegrationTestClient(t *testing.T) Events {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
