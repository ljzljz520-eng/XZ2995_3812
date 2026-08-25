package store

import "encoding/json"

func marshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

func unmarshal(value []byte, target any) error {
	return json.Unmarshal(value, target)
}

func cloneBytes(value []byte) []byte {
	if value == nil {
		return nil
	}
	copyValue := make([]byte, len(value))
	copy(copyValue, value)
	return copyValue
}

func encodeKey(parts ...string) string {
	result := ""
	for index, part := range parts {
		if index > 0 {
			result += ":"
		}
		result += part
	}
	return result
}
