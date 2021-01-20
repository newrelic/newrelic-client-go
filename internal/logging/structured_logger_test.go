// +build unit

package logging

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewStructuredLogger_private(t *testing.T) {
	t.Parallel()

	logrusLevel := log.GetLevel()

	l := NewStructuredLogger()

	l.SetLevel("debug")

	// Make sure that the global logrus level is left unchanged.
	assert.Equal(t, log.GetLevel(), logrusLevel)

	loggerLevel := l.logger.GetLevel()
	assert.Equal(t, loggerLevel, log.DebugLevel)
}

func TestSetLevel(t *testing.T) {
	t.Parallel()
	l := NewStructuredLogger()

	defaultLevel, err := log.ParseLevel(defaultLogLevel)

	l.SetLevel("")
	loggerLevel := l.logger.GetLevel()

	assert.NoError(t, err)
	assert.Equal(t, loggerLevel, defaultLevel)

	l.SetLevel("notdebug")
	loggerLevel = l.logger.GetLevel()
	assert.Equal(t, loggerLevel, defaultLevel)

}

func TestStructuredLogger_interface(t *testing.T) {
	var l Logger = NewStructuredLogger()
	l.Info("testing")
}
