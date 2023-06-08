package godotenv

import (
	"testing"
)

var noopPresets = make(map[string]string)

func parseAndCompare(t *testing.T, rawEnvLine string, expectedKey string, expectedValue string) {
	result, err := Unmarshal(rawEnvLine)

	if err != nil {
		t.Errorf("Expected %q to parse as %q: %q, errored %q", rawEnvLine, expectedKey, expectedValue, err)
		return
	}
	if result[expectedKey] != expectedValue {
		t.Errorf("Expected '%v' to parse as '%v' => '%v', got %q instead", rawEnvLine, expectedKey, expectedValue, result)
	}
}

func TestParsing(t *testing.T) {
	// unquoted values
	parseAndCompare(t, "FOO=BAR", "FOO", "BAR")

	// parses values with spaces around equal sign
	parseAndCompare(t, "FOO =bar", "FOO", "bar")
	parseAndCompare(t, "FOO= bar", "FOO", "bar")

	// parses double quoted values
	parseAndCompare(t, `FOO="bar"`, "FOO", "bar")

	// parses single quoted values
	parseAndCompare(t, `FOO='bar'`, "FOO", "bar")

	// parses escaped double quoted
	parseAndCompare(t, `FOO="escaped\"bar"`, "FOO", `escaped"bar`)

	// parses single quotes inside double quotes
	parseAndCompare(t, `FOO="'d'"`, "FOO", `'d'`)

	// parses yaml style options
	parseAndCompare(t, "OPTION_A: 1", "OPTION_A", "1")

	// parses yaml values with equal signs
	parseAndCompare(t, "OPTION_A: FOO=bar", "OPTION_A", "FOO=bar")

	// parses non-yaml options with colons
	parseAndCompare(t, "OPTION_A=1:B", "OPTION_A", "1:B")

	// parses export keyword
	parseAndCompare(t, "export OPTION_A=2", "OPTION_A", "2")
	parseAndCompare(t, `export OPTION_B='\n'`, "OPTION_B", "\\n")
	parseAndCompare(t, "export exportFoo=2", "exportFoo", "2")
	parseAndCompare(t, "exportFoo=2", "exportFoo", "2")
	parseAndCompare(t, "export_Foo= 2", "export_Foo", "2")
	parseAndCompare(t, "export.Foo =2", "export.Foo", "2")
	parseAndCompare(t, "export\tOPTION_A=2", "OPTION_A", "2")
	parseAndCompare(t, "  export OPTION_A=2", "OPTION_A", "2")
	parseAndCompare(t, "\texport OPTION_A=2", "OPTION_A", "2")

	// it 'expands newline in quoted strings' do
	// expect(env('FOO="bar\nbaz"')).to eql('FOO' => "bar\nbaz")
	parseAndCompare(t, `FOO="bar\nbaz"`, "FOO", "bar\nbaz")

	// it 'parses variables with "." in the name' do
	// expect(env('FOO=foobar=')).to eql('FOO' => "foobar=")
	parseAndCompare(t, "FOO=foobar=", "FOO", "foobar=")

	// it 'strips unquoted values' do
	// expect(env('FOO=bar ')).to eql('FOO' => "bar")
	parseAndCompare(t, "FOO=bar ", "FOO", "bar")

	// unquoted internal whitespace is preserved
	parseAndCompare(t, "KEY=value value", "KEY", "value value")

	// it 'ignores inline comments' do
	// expect(env('foo=bar # this is foo')).to eql('foo' => 'bar')
	parseAndCompare(t, "foo=bar # this is foo", "foo", "bar")

	// it 'allows # in quoted value' do
	// expect(env('foo="bar#baz" # this is foo')).to eql('foo' => 'bar#baz')
	parseAndCompare(t, `foo="bar#baz" # comment`, "foo", "bar#baz")
	parseAndCompare(t, `foo='bar#baz' # comment`, "foo", "bar#baz")
	parseAndCompare(t, `foo='bar#baz#bang' # comment`, "foo", "bar#baz#bang")

	parseAndCompare(t, `="value"`, "", "value")

	// unquoted whitespace around keys should be ignored
	parseAndCompare(t, "KEY =value", "KEY", "value")
	parseAndCompare(t, "   KEY=value", "KEY", "value")
	parseAndCompare(t, "\tKEY=value", "KEY", "value")

	// it 'throws an error if line format is incorrect' do
	// expect(env('lol$wut')).to raise_error(Dotenv::FormatError)
	badlyFormattedLine := "lol$wut"
	_, err := Unmarshal(badlyFormattedLine)
	if err == nil {
		t.Errorf("Expected \"%v\" to return error, but it didn't", badlyFormattedLine)
	}
}
