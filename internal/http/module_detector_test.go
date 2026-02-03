//go:build unit
// +build unit

package http

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func TestIsTestBinary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Go test binary",
			input:    "http.test",
			expected: true,
		},
		{
			name:     "Leading underscore",
			input:    "_test_binary",
			expected: true,
		},
		{
			name:     "go-build pattern",
			input:    "/tmp/go-build123/exe/test",
			expected: true,
		},
		{
			name:     "Temp directory pattern",
			input:    "/tmp/something",
			expected: true,
		},
		{
			name:     "Normal binary name",
			input:    "my-service",
			expected: false,
		},
		{
			name:     "Normal executable",
			input:    "monitoring-app",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := isTestBinary(tt.input)
			assert.Equal(t, tt.expected, result, "isTestBinary(%q) = %v, want %v", tt.input, result, tt.expected)
		})
	}
}

func TestDetectCallingService(t *testing.T) {
	t.Parallel()

	// Test that detectCallingService never panics, even with edge cases
	t.Run("does not panic", func(t *testing.T) {
		t.Parallel()

		// Should not panic - just return empty or something valid
		assert.NotPanics(t, func() {
			result := detectCallingService()
			// Result can be empty or non-empty, we just care it doesn't panic
			_ = result
		})
	})

	// Test that it returns a string (empty or non-empty)
	t.Run("returns string", func(t *testing.T) {
		t.Parallel()

		result := detectCallingService()
		// Should return a string, even if empty
		assert.IsType(t, "", result)
	})

	// Test that when running in a test environment, it handles test binaries gracefully
	t.Run("handles test environment", func(t *testing.T) {
		t.Parallel()

		result := detectCallingService()
		// In test environment, it should either return empty or a non-test name
		// We don't assert specific values since it depends on how tests are run
		// but we verify it doesn't return test binary patterns if detected
		if result != "" {
			assert.NotContains(t, result, ".test", "should not return test binary name")
		}
	})
}

func TestDetectFromModulePath(t *testing.T) {
	t.Parallel()

	t.Run("does not panic", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() {
			result := detectFromModulePath()
			_ = result
		})
	})

	t.Run("filters out newrelic-client-go itself", func(t *testing.T) {
		t.Parallel()

		result := detectFromModulePath()
		// Should not return the client library's own module path
		assert.NotEqual(t, "github.com/newrelic/newrelic-client-go/v2", result)
		assert.NotEqual(t, "github.com/newrelic/newrelic-client-go", result)
	})
}

func TestDetectFromBinaryName(t *testing.T) {
	t.Parallel()

	t.Run("does not panic", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() {
			result := detectFromBinaryName()
			_ = result
		})
	})

	t.Run("filters out main", func(t *testing.T) {
		t.Parallel()

		result := detectFromBinaryName()
		// Should not return "main" as it's not useful
		assert.NotEqual(t, "main", result)
	})

	t.Run("filters out test binaries", func(t *testing.T) {
		t.Parallel()

		result := detectFromBinaryName()
		// Should not return test binary patterns
		if result != "" {
			assert.NotContains(t, result, ".test")
			assert.False(t, len(result) > 0 && result[0] == '_', "should not start with underscore")
		}
	})
}

func TestDetectFromProcessArgs(t *testing.T) {
	t.Parallel()

	t.Run("does not panic", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() {
			result := detectFromProcessArgs()
			_ = result
		})
	})

	t.Run("handles empty os.Args gracefully", func(t *testing.T) {
		t.Parallel()

		// Save original args
		originalArgs := os.Args

		// Test with empty args
		os.Args = []string{}
		result := detectFromProcessArgs()
		assert.Equal(t, "", result, "should return empty string when os.Args is empty")

		// Restore
		os.Args = originalArgs
	})

	t.Run("filters out main", func(t *testing.T) {
		t.Parallel()

		result := detectFromProcessArgs()
		// Should not return "main"
		assert.NotEqual(t, "main", result)
	})
}

func TestClientServiceNameWithAutoDetection(t *testing.T) {
	// Note: This test runs in a test environment where auto-detection
	// typically returns empty (due to test binary filtering), so we verify
	// the fallback behavior works correctly

	t.Run("empty service name falls back to default when detection fails", func(t *testing.T) {
		cfg := config.Config{}
		client := NewClient(cfg)

		// When auto-detection fails (common in tests), should use default
		// The serviceName is set during NewClient, we can verify by making a request
		// and checking the header would be set correctly
		assert.NotEmpty(t, client.config.ServiceName)
	})

	t.Run("custom service name is preserved", func(t *testing.T) {
		cfg := config.Config{
			ServiceName: "my-custom-app",
		}
		client := NewClient(cfg)

		// Custom service name should be preserved with library name appended
		assert.Contains(t, client.config.ServiceName, "my-custom-app")
		assert.Contains(t, client.config.ServiceName, "newrelic-client-go")
	})

	t.Run("env variable prepends to service name", func(t *testing.T) {
		// Set env variable
		originalEnv := os.Getenv("NEW_RELIC_SERVICE_NAME")
		os.Setenv("NEW_RELIC_SERVICE_NAME", "env-service")
		defer func() {
			if originalEnv != "" {
				os.Setenv("NEW_RELIC_SERVICE_NAME", originalEnv)
			} else {
				os.Unsetenv("NEW_RELIC_SERVICE_NAME")
			}
		}()

		cfg := config.Config{
			ServiceName: "my-app",
		}
		client := NewClient(cfg)

		// Should have env variable prepended
		assert.Contains(t, client.config.ServiceName, "env-service")
		assert.Contains(t, client.config.ServiceName, "my-app")
		assert.Contains(t, client.config.ServiceName, "newrelic-client-go")
	})

	t.Run("default service name is not duplicated", func(t *testing.T) {
		cfg := config.Config{
			// Empty ServiceName - will fall back to default after failed detection
		}
		client := NewClient(cfg)

		// Should not have "newrelic-client-go|newrelic-client-go"
		// Count occurrences - there should be exactly one
		serviceName := client.config.ServiceName
		count := 0
		searchStr := "newrelic-client-go"
		for i := 0; i <= len(serviceName)-len(searchStr); i++ {
			if serviceName[i:i+len(searchStr)] == searchStr {
				count++
			}
		}
		assert.Equal(t, 1, count, "should have exactly one occurrence of newrelic-client-go, got: %s", serviceName)
	})
}
