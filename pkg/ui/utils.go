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

type Path []string

func (p Path) String() string {
	for i, s := range p {
		p[i] = kebabToPascalCase(s)
	}

	return strings.Join(p, "")
}

func (p Path) append(s string) Path {
	return append(p, s)
}

func (p Path) optionsCheckboxId() string {
	return fmt.Sprintf("%sOptionsCheckbox", p)
}

func (p Path) optionValuesFormId() string {
	return fmt.Sprintf("%sOptionValuesForm", p)
}

func (p Path) optionValueInputId(option string) string {
	return fmt.Sprintf("%s%sOptionValueInput", p, kebabToPascalCase(option))
}

func (p Path) subcommandTabsId() string {
	return fmt.Sprintf("%sSubcommandTabs", p)
}

func (p Path) commandStackId() string {
	return fmt.Sprintf("%sCommandStack", p)
}
