package http

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

// detectCallingService attempts to auto-detect the calling service name
// using various strategies with fallback priority.
// This function is designed to never panic or return errors - it simply
// returns an empty string if detection fails, allowing the caller to use
// a default value. This ensures client initialization always succeeds.
func detectCallingService() string {
	// Priority 1: Go module path from build info
	// This returns the MAIN module that was built (the customer's app), not the library
	if modulePath := detectFromModulePath(); modulePath != "" {
		return modulePath
	}

	// Priority 2: Binary name from executable path
	if binaryName := detectFromBinaryName(); binaryName != "" {
		return binaryName
	}

	// Priority 3: Process name from Args (last resort)
	if processName := detectFromProcessArgs(); processName != "" {
		return processName
	}

	// If all detection methods fail, return empty string
	// The caller will use the default service name
	return ""
}

// detectFromModulePath attempts to detect the calling service from Go module build info
func detectFromModulePath() string {
	// Wrap in a recovery function to handle any panics from debug.ReadBuildInfo
	defer func() {
		if r := recover(); r != nil {
			// Silently recover - we don't want detection to crash the client
		}
	}()

	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Path == "" {
		return ""
	}

	// Filter out the newrelic-client-go library itself
	// This handles cases where someone is building/testing the library directly
	if info.Main.Path == "github.com/newrelic/newrelic-client-go/v2" ||
		info.Main.Path == "github.com/newrelic/newrelic-client-go" {
		return ""
	}

	return info.Main.Path
}

// detectFromBinaryName attempts to detect the calling service from the binary name
func detectFromBinaryName() string {
	defer func() {
		if r := recover(); r != nil {
			// Silently recover
		}
	}()

	execPath, err := os.Executable()
	if err != nil {
		return ""
	}

	binaryName := filepath.Base(execPath)

	// Filter out unhelpful binary names
	if binaryName == "" || binaryName == "main" || isTestBinary(binaryName) {
		return ""
	}

	return binaryName
}

// detectFromProcessArgs attempts to detect the calling service from process arguments
func detectFromProcessArgs() string {
	defer func() {
		if r := recover(); r != nil {
			// Silently recover
		}
	}()

	if len(os.Args) == 0 {
		return ""
	}

	processName := filepath.Base(os.Args[0])

	// Filter out unhelpful process names
	if processName == "" || processName == "main" || isTestBinary(processName) {
		return ""
	}

	return processName
}

// isTestBinary checks if the binary name suggests it's a test/temp binary
// Returns true for Go test binaries, temp binaries from go run, etc.
func isTestBinary(name string) bool {
	if name == "" {
		return false
	}

	// Check for leading underscore (common in test binaries)
	if name[0] == '_' {
		return true
	}

	// Check for common test/temp patterns
	testPatterns := []string{
		".test",    // Go test binaries: package.test
		"___",      // Some temp binary patterns
		"go-build", // Temp binaries from go run
		"/tmp/",    // Temp directory paths
		"T/",       // macOS temp directories
	}

	for _, pattern := range testPatterns {
		if strings.Contains(name, pattern) {
			return true
		}
	}

	return false
}
