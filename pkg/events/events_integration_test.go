// +build integration

package events

import (
	"testing"

	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationEvents(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	event := struct {
		EventType string  `json:"eventType"`
		Amount    float64 `json:"amount"`
	}{
		EventType: "Purchase",
		Amount:    123.45,
	}

	// Test: Create
	err := client.CreateEvent(nr.TestAccountID, event)

	require.NoError(t, err)
}

func TestIntegrationEvents_Compression(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	event := struct {
		EventType string `json:"eventType"`
		Data      string `json:"data"`
	}{
		EventType: "Dessert",
		Data:      "Biscuit cotton candy candy canes dessert muffin ice cream carrot cake. Marzipan lemon drops lemon drops. Candy cake toffee. Powder pie wafer bonbon dessert powder. Jujubes sweet chocolate bar gummies jelly-o. Wafer biscuit candy. Oat cake cookie jelly liquorice cupcake cupcake dragée cupcake wafer. Wafer chocolate bar marzipan powder jujubes cake oat cake sweet roll. Chocolate bar caramels jelly sugar plum donut. Candy donut tiramisu candy canes icing macaroon. Gummies macaroon jujubes candy gummies cotton candy sesame snaps dragée. Sweet roll icing cake pie sweet. Candy topping toffee. Sesame snaps cookie lemon drops wafer jujubes powder fruitcake.",
	}

	// Test: Create
	err := client.CreateEvent(nr.TestAccountID, event)

	require.NoError(t, err)
}

func newIntegrationTestClient(t *testing.T) Events {
	tc := nr.NewIntegrationTestConfig(t)

	return New(tc)
}
