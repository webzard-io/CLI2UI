package ui

import (
	"CLI2UI/pkg/config"
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

func parseOptions(c config.Command) ([]CheckboxOption, []string) {
	options := []CheckboxOption{}
	requiredOptions := []string{}

	for _, f := range c.Flags {
		options = append(options, CheckboxOption{
			Label:    f.Name,
			Value:    f.Default,
			Disabled: f.Required,
		})

		if f.Required {
			requiredOptions = append(requiredOptions, f.Name)
		}
	}

	for _, a := range c.Args {
		options = append(options, CheckboxOption{
			Label:    a.Name,
			Value:    a.Default,
			Disabled: a.Required,
		})

		if a.Required {
			requiredOptions = append(requiredOptions, a.Name)
		}
	}

	return options, requiredOptions
}
