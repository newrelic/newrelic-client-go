package metrics

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func TestMetricForwarding(t *testing.T) {
	// Initialize the client configuration. A New Relic License Key is required
	// to communicate with the backend API.
	cfg := config.New()
	cfg.LicenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
	cfg.Compression = config.Compression.None

	// Initialize the client.
	client := New(cfg)

	metric := struct {
		Name       string            `json:"name"`
		Type       string            `json:"type"`
		Value      int               `json:"value"`
		Attributes map[string]string `json:"attributes"`
	}{
		Name:  "service.errors.all",
		Type:  "gauge",
		Value: 9,
		Attributes: map[string]string{
			"service.response.statuscode": "400",
		},
	}

	// Post a Metric
	if err := client.CreateMetricEntry(metric); err != nil {
		log.Fatal("error posting metric: ", err)
	}

	fmt.Printf("success")
}
