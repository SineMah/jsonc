package jsonc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecognizeJsonc(t *testing.T) {

	t.Run("test json with no comments", func(t *testing.T) {
		bytes := []byte(`{"port": 1111, "port_as_string":"1110", "foo": {"bar": "baz"}, "is_true": true}`)

		assert.Equal(t, false, IsJsonc(bytes))
	})

	t.Run("json is no json with slashin string", func(t *testing.T) {
		bytes := []byte(`{"app//port": 1111, "port_as_string":"1110", "foo": {"bar": "baz"}, "is_true": true}`)

		assert.Equal(t, false, IsJsonc(bytes))
	})

	t.Run("jsonc with block comments recognized", func(t *testing.T) {
		bytes := []byte(`{"port": 1111, /*foo bar*/ "port_as_string":"1110", "foo": {"bar": "baz"}, "is_true": true}`)

		assert.Equal(t, true, IsJsonc(bytes))
	})

	t.Run("jsonc with line comments recognized", func(t *testing.T) {
		bytes := []byte(`
			{
				"port": 1111,
				// foo bar
				"port_as_string":"1110",
				"foo":
				{
					"bar": "baz"
				},
				"is_true": true
			}`)

		assert.Equal(t, true, IsJsonc(bytes))
	})
}

func TestRemoveCommentsFromJsonc(t *testing.T) {

	t.Run("test json with no comments", func(t *testing.T) {
		bytesNoComment := []byte(`{"port": 1111, "port_as_string":"1110", "foo": {"bar": "baz"}, "is_true": true}`)
		bytes := []byte(`{"port": 1111, "port_as_string":"1110", "foo": {"bar": "baz"}, "is_true": true}`)

		v, _ := RemoveComments(bytes)

		assert.Equal(t, bytesNoComment, v)
	})

	t.Run("test json with block comment", func(t *testing.T) {
		bytesWithComment := []byte(`{"port": 1111, /* 42 */"port_as_string":"1110", "foo": {"bar": "baz"}, "is_true": true}`)
		bytes := []byte(`{"port": 1111, "port_as_string":"1110", "foo": {"bar": "baz"}, "is_true": true}`)

		v, _ := RemoveComments(bytesWithComment)

		assert.Equal(t, bytes, v)
	})

	t.Run("test json with line comment", func(t *testing.T) {
		bytesWithComment := []byte(`{
			"port": "1111",
			// foo bar
			"port_as_string":"1110",
			"foo": {
				"bar": "baz"
			},
			"is_true": true
		}`)
		controlBytes := []byte(`{"port":"1111","port_as_string":"1110","foo":{"bar":"baz"},"is_true":true}`)

		v, _ := RemoveComments(bytesWithComment)
		raw := json.RawMessage(v)
		bytesNoComment, _ := json.Marshal(&raw)

		assert.Equal(t, controlBytes, bytesNoComment)
	})
}
