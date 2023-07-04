package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func structToMap(s interface{}) map[string]interface{} {
	data, err := json.Marshal(s)
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

func kebabToPascalCase(i string) string {
	words := strings.Split(i, "-")
	for i, word := range words {
		words[i] = cases.Title(language.English).String(word)
	}
	return strings.Join(words, "")
}

func optionsCheckboxId(cmd string) string {
	return fmt.Sprintf("%sOptionsCheckbox", cmd)
}

func optionValueInputId(cmd string, option string) string {
	return fmt.Sprintf("Command%sOption%sInput", cmd, kebabToPascalCase(option))
}

func optionValuesFormId(cmd string) string {
	return fmt.Sprintf("Command%sOptionValuesForm", cmd)
}
