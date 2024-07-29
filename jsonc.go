package jsonc

import (
	"bytes"
	"encoding/json"
)

func RemoveComments(data []byte) []byte {
	var buffer bytes.Buffer
	inString := false
	inSingleLineComment := false
	inMultiLineComment := false

	for i := 0; i < len(data); i++ {
		switch {
		case inSingleLineComment:
			if data[i] == '\n' {
				inSingleLineComment = false
				buffer.WriteByte(data[i])
			}
		case inMultiLineComment:
			if iteratorInRange(data, i) && data[i] == '*' && data[i+1] == '/' {
				inMultiLineComment = false
				i++
			}
		case data[i] == '"':
			inString = !inString
			buffer.WriteByte(data[i])
		case !inString && data[i] == '/' && iteratorInRange(data, i):
			switch data[i+1] {
			case '/':
				inSingleLineComment = true
				i++
			case '*':
				inMultiLineComment = true
				i++
			default:
				buffer.WriteByte(data[i])
			}
		default:
			buffer.WriteByte(data[i])
		}
	}

	return buffer.Bytes()
}

func iteratorInRange(data []byte, i int) bool {

	return i+1 < len(data)
}

func Unmarshal(data []byte, v any) error {
	if IsJsonc(data) {
		data = RemoveComments(data)
	}

	return json.Unmarshal(data, v)
}

func IsJsonc(data []byte) bool {
	stringVal := 0

	for _, c := range data {

		if c == '"' {
			stringVal++
		}

		if stringVal%2 == 1 {
			continue
		}

		if c == '/' || c == '*' {
			return true
		}
	}

	return false
}
