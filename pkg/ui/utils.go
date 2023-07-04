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

type Prefix string

func prefix(p Prefix, cmd string) Prefix {
	return Prefix(fmt.Sprintf("%s%s", p, kebabToPascalCase(cmd)))
}

func optionsCheckboxId(p Prefix) string {
	return fmt.Sprintf("%sOptionsCheckbox", p)
}

func optionValuesFormId(p Prefix) string {
	return fmt.Sprintf("%sOptionValuesForm", p)
}

func optionValueInputId(p Prefix, option string) string {
	return fmt.Sprintf("%s%sOptionValueInput", p, kebabToPascalCase(option))
}

func subcommandTabsId(p Prefix) string {
	return fmt.Sprintf("%sSubcommandTabs", p)
}

func commandStackId(p Prefix) string {
	return fmt.Sprintf("%sCommandStack", p)
}
