package logging

import "testing"

func TestMockLogger(t *testing.T) {
	var l Logger = MockLogger{t: t}
	l.Info("testing")
}
