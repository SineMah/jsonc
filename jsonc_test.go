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

	t.Run("json is no json with slash in string", func(t *testing.T) {
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

		v := RemoveComments(bytes)

		assert.Equal(t, bytesNoComment, v)
	})

	t.Run("test json with block comment", func(t *testing.T) {
		bytesWithComment := []byte(`{"port": 1111, /* 42 */"port_as_string":"1110", "foo": {"bar": "baz"}, "is_true": true}`)
		bytes := []byte(`{"port": 1111, "port_as_string":"1110", "foo": {"bar": "baz"}, "is_true": true}`)

		v := RemoveComments(bytesWithComment)

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

		v := RemoveComments(bytesWithComment)
		raw := json.RawMessage(v)
		bytesNoComment, _ := json.Marshal(&raw)

		assert.Equal(t, controlBytes, bytesNoComment)
	})

	t.Run("test json with comments", func(t *testing.T) {
		bytesWithComment := []byte(`{
			"port": "1111",
			// foo bar
			"port_as_string":"1110",
			"foo": {
				"bar": "baz"
			},
			/* is 42? */
			"is_true": true
		}`)
		controlBytes := []byte(`{"port":"1111","port_as_string":"1110","foo":{"bar":"baz"},"is_true":true}`)

		v := RemoveComments(bytesWithComment)
		raw := json.RawMessage(v)
		bytesNoComment, _ := json.Marshal(&raw)

		assert.Equal(t, controlBytes, bytesNoComment)
	})
}

func TestUnmarshalJsonc(t *testing.T) {

	t.Run("unmarshal jsonc with comments", func(t *testing.T) {
		m := make(map[string]any)

		bytesWithComment := []byte(`{
			"port": "1111",
			// foo bar
			"port_as_string":"1110",
			"foo": {
				"bar": "baz"
			},
			/* is 42? */
			"is_true": true
		}`)
		controlBytes := []byte(`{"foo":{"bar":"baz"},"is_true":true,"port":"1111","port_as_string":"1110"}`)

		_ = Unmarshal(bytesWithComment, &m)
		bytesNoComment, _ := json.Marshal(&m)

		assert.Equal(t, controlBytes, bytesNoComment)
	})
}
