package util

import (
	"bytes"
	"encoding/json"
)

// MetadataToString Encode the map into a string
func MetadataToString(data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return bytes.NewBuffer(jsonData).String(), nil
}

// MetadataFromString Decode the JSON string into a map
func MetadataFromString(jsonString string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
