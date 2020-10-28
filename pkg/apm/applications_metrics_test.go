package apm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalTimesliceValues(t *testing.T) {
	data := []byte(`{
		"as_percentage": 10.3,
		"average_time": 133.1,
		"calls_per_minute": 3029.409,
		"max_value": 500,
		"total_call_time_per_minute": 123.122,
		"standard_deviation": 100.155,
		"other_key": 22.11
	}`)

	expect := &MetricTimesliceValues{
		AsPercentage:           10.3,
		AverageTime:            133.1,
		CallsPerMinute:         3029.409,
		MaxValue:               500.0,
		TotalCallTimePerMinute: 123.122,
		Utilization:            0.0,
		Values: map[string]float64{
			"as_percentage":              10.3,
			"average_time":               133.1,
			"calls_per_minute":           3029.409,
			"max_value":                  500.0,
			"total_call_time_per_minute": 123.122,
			"standard_deviation":         100.155,
			"other_key":                  22.11,
		},
	}

	values := MetricTimesliceValues{}

	require.NoError(t, json.Unmarshal(data, &values))

	require.Equal(t, expect, &values)
}
