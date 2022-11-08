//go:build integration
// +build integration

package logs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

const (
	BatchTimeoutSeconds = 5
)

//func TestExample_log_batch(t *testing.T) {
func Example_basic() {
	// Initialize the client configuration.  A New Relic License Key is required to communicate with the backend API.
	cfg := config.New()
	cfg.LicenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
	cfg.LogLevel = "trace"
	cfg.Compression = config.Compression.None

	// Initialize the client.
	client := New(cfg)
	client.batchTimeout = BatchTimeoutSeconds * time.Second

	logEntries := []struct {
		Message string `json:"message"`
		LogType string `json:"logType"`
	}{
		{
			Message: "INFO: From example_log_batch_test.go message 1",
			LogType: "Teapot",
		},
		{
			Message: "INFO: From example_log_batch_test.go message 2",
			LogType: "Teapot",
		}}

	// Start batch mode, 12345 is a placeholder to pass PR compilation
	if err := client.BatchMode(context.Background(), 12345); err != nil {
		log.Fatal("error starting batch mode: ", err)
	}

	// Queue log entries.
	for _, e := range logEntries {
		if err := client.EnqueueLogEntry(context.Background(), e); err != nil {
			log.Fatal("error queuing log entry: ", err)
		}
	}

	// Hack to stop the test from ending BEFORE the queue flushes
	time.Sleep((BatchTimeoutSeconds * 2) * time.Second)
	// Force flush
	if err := client.Flush(); err != nil {
		log.Fatal("error flushing log queue: ", err)
	}
}
