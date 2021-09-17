//go:build unit
// +build unit

package logging

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogrusLogger_private(t *testing.T) {
	t.Parallel()

	logrusLevel := log.GetLevel()

	l := NewLogrusLogger()

	l.SetLevel("debug")

	// Make sure that the global logrus level is left unchanged.
	assert.Equal(t, log.GetLevel(), logrusLevel)

	loggerLevel := l.logger.GetLevel()
	assert.Equal(t, loggerLevel, log.DebugLevel)
}

func TestSetLevel(t *testing.T) {
	t.Parallel()
	l := NewLogrusLogger()

	defaultLevel, err := log.ParseLevel(defaultLogLevel)

	l.SetLevel("")
	loggerLevel := l.logger.GetLevel()

	assert.NoError(t, err)
	assert.Equal(t, loggerLevel, defaultLevel)

	l.SetLevel("notdebug")
	loggerLevel = l.logger.GetLevel()
	assert.Equal(t, loggerLevel, defaultLevel)

}

func TestLogrusLogger_interface(t *testing.T) {
	var l Logger = NewLogrusLogger()
	l.Info("testing")
}
