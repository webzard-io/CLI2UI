package ui

import (
	"encoding/json"
	"log"
)

func structToMap(inputStruct interface{}) map[string]interface{} {
	data, err := json.Marshal(inputStruct)
	if err != nil {
		log.Fatal("error marshaling struct to JSON", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal("error unmarshaling JSON to map", err)
	}

	return result
}
