package changetracking

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadCustomAttributesJS_FromString(t *testing.T) {
	js := `{foo: "bar", isProd: true, count: 42, pi: 3.14}`
	attrs, err := ReadCustomAttributesJS(js, false)
	require.NoError(t, err)
	require.Equal(t, "bar", attrs["foo"])
	require.Equal(t, true, attrs["isProd"])
	require.Equal(t, int64(42), attrs["count"])
	require.Equal(t, 3.14, attrs["pi"])
}

func TestReadCustomAttributesJS_FromFile(t *testing.T) {
	fileContent := `{region: "us-east-1", enabled: false, num: 7}`
	fileName := "test_custom_attrs.jsobj"
	err := os.WriteFile(fileName, []byte(fileContent), 0644)
	require.NoError(t, err)
	defer os.Remove(fileName)

	attrs, err := ReadCustomAttributesJS(fileName, true)
	require.NoError(t, err)
	require.Equal(t, "us-east-1", attrs["region"])
	require.Equal(t, false, attrs["enabled"])
	require.Equal(t, int64(7), attrs["num"])
}

func TestReadCustomAttributesJS_InvalidJS(t *testing.T) {
	js := `{foo: bar,}` // bar is not quoted, not a variable
	_, err := ReadCustomAttributesJS(js, false)
	require.Error(t, err)
}

func TestReadCustomAttributesJS_NotAMap(t *testing.T) {
	js := `42` // Not an object
	_, err := ReadCustomAttributesJS(js, false)
	require.Error(t, err)
}
