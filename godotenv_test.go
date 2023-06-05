package godotenv

import (
	"testing"
)

var noopPresets = make(map[string]string)

func parseAndCompare(t *testing.T, rawEnvLine string, expectedKey string, expectedValue string) {
	result, err := Unmarshal(rawEnvLine)

	if err != nil {
		t.Error("Expected %q to parse as %q: %q, errored %q", rawEnvLine, expectedKey, expectedValue, err)
		return
	}
	if result[expectedKey] != expectedValue {
		t.Error("Expected '%v' to parse as '%v' => '%v', got %q instead", rawEnvLine, expectedKey, expectedValue, result)
	}
}

func TestParsing(t *testing.T) {
	// unquoted values
	parseAndCompare(t, "FOO=BAR", "FOO", "BAR")
}
