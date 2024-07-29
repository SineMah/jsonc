package jsonc

import (
	"encoding/json"
	"log"
)

func RemoveComments(data []byte) ([]byte, error) {
	stringVal := 0
	commentsMap := make(map[int][]int)

	for i, c := range data {

		if c == '"' {
			stringVal++
		}

		if stringVal%2 == 1 {
			continue
		}

		if i >= len(data)-1 {
			continue
		}

		cNext := data[i+1]

		if c == '/' && cNext == '/' {
			log.Print("found comment1")
			commentsMap[len(commentsMap)] = []int{i, findNextBytes(data, []byte{'\n'}, i)}
		}

		if c == '/' && cNext == '*' {
			log.Print("found comment2")
			commentsMap[len(commentsMap)] = []int{i, findNextBytes(data, []byte{'*', '/'}, i)}
		}
	}

	for i := len(commentsMap) - 1; i >= 0; i-- {
		start := commentsMap[i][0]
		end := commentsMap[i][1]
		data = append(data[:start], data[end:]...)
	}

	return data, nil
}

func findNextBytes(haystack []byte, search []byte, start int) int {
	for i, c := range haystack {
		if i < start {
			continue
		}

		if c == search[0] {
			found := true

			for j, s := range search {
				if haystack[i+j] != s {
					found = false
					break
				}
			}

			if found {
				return i + len(search)
			}
		}
	}

	return -1
}

func Unmarshal(data []byte, v any) error {
	if IsJsonc(data) {
		var err error
		data, err = RemoveComments(data)

		if err != nil {
			return err
		}
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
