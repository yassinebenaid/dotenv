package dotenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	input := []byte(`abcd`)

	parser := parser{input: input}

	assert.Equal(t, "a", string(parser.next()))
	assert.Equal(t, "b", string(parser.next()))
	assert.Equal(t, "c", string(parser.next()))
	assert.Equal(t, "d", string(parser.next()))
	assert.Equal(t, uint8(0x0), parser.next())
}

func TestParse(t *testing.T) {
	testCases := []struct {
		input    string
		expected map[string]string
	}{
		{"", map[string]string{}},
		{"key1=", map[string]string{"key1": ""}},
		{"key1=value-1", map[string]string{"key1": "value-1"}},
		{"key1= hello @ world 123 _+-;:, ", map[string]string{"key1": "hello @ world 123 _+-;:,"}},
		{"key1= value #comment", map[string]string{"key1": "value"}},
		{"key1=#comment", map[string]string{"key1": ""}},
		{"#comment", map[string]string{}},
		{`
	KEY_1 = value-1
	KEY_2 = value-2
	`, map[string]string{
			"KEY_1": "value-1",
			"KEY_2": "value-2",
		}},
		{`
	key1 = value-1
	key2 = value-2
	`, map[string]string{
			"key1": "value-1",
			"key2": "value-2",
		}},
		{"key1=value-1\nkey2=value-2", map[string]string{
			"key1": "value-1",
			"key2": "value-2",
		}},
	}

	for _, tc := range testCases {

		result, err := parse([]byte(tc.input))
		assert.NoError(t, err)

		assert.NotNil(t, result)

		assert.Equal(t, tc.expected, result)
	}
}
