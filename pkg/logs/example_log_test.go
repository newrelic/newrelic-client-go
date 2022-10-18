//go:build integration
// +build integration

package logs

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func Example_log(t *testing.T) {
	// Initialize the client configuration.  A New Relic License Key is required
	// to communicate with the backend API.
	cfg := config.New()
	cfg.LicenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
	cfg.LogLevel = "trace"
	cfg.Compression = config.Compression.None

	// Initialize the client.
	client := New(cfg)

	logEntry := struct {
		Message string `json:"message"`
	}{
		Message: "INFO: From example_log_test.go",
	}

	// Post a Log entry
	if err := client.CreateLogEntry(logEntry); err != nil {
		log.Fatal("error posting Log entry: ", err)
	}

	fmt.Printf("success")
}
