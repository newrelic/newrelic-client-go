package serialization

import (
	"fmt"
	"strconv"
	"time"
)

// MarshalEpochMillisecondsJSON returns JSON bytes for a Unix millisecond count
// (GraphQL EpochMilliseconds scalar), including when sub-second nanoseconds are zero.
func MarshalEpochMillisecondsJSON(t time.Time) ([]byte, error) {
	return []byte(strconv.FormatInt(t.UTC().UnixMilli(), 10)), nil
}

// UnmarshalEpochMillisecondsJSON parses JSON for EpochMilliseconds: null, numeric
// epoch forms handled by EpochTime, RFC3339 / RFC3339Nano strings (optionally JSON-quoted).
func UnmarshalEpochMillisecondsJSON(s []byte, e *EpochTime) error {
	if string(s) == "null" {
		return nil
	}

	inner := s
	if len(s) >= 2 && s[0] == '"' {
		out, err := strconv.Unquote(string(s))
		if err != nil {
			return fmt.Errorf("epoch milliseconds: unquote: %w", err)
		}
		inner = []byte(out)
	}

	if len(inner) == 0 {
		return fmt.Errorf("epoch milliseconds: empty value")
	}

	if isEpochNumericToken(inner) {
		return (*EpochTime)(e).UnmarshalJSON(inner)
	}

	if tm, err := time.Parse(time.RFC3339Nano, string(inner)); err == nil {
		*(*time.Time)(e) = tm.UTC()
		return nil
	}
	if tm, err := time.Parse(time.RFC3339, string(inner)); err == nil {
		*(*time.Time)(e) = tm.UTC()
		return nil
	}

	return fmt.Errorf("epoch milliseconds: unable to parse %q", inner)
}

func isEpochNumericToken(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	if string(b) == emptyTimeCase {
		return true
	}
	for i, c := range b {
		if i == 0 && c == '-' {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
