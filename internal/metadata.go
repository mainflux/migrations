package util

import (
	"encoding/json"
)

// MetadataToString Encode the map into a string.
func MetadataToString(data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// MetadataFromString Decode the JSON string into a map.
func MetadataFromString(jsonString string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return nil, err
	}

	return data, nil
}
