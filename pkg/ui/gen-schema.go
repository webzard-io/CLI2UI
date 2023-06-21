package ui

import (
	"encoding/json"
	"fmt"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"strconv"
	"strings"
)

func kebabToPascalCase(input string) string {
	words := strings.Split(input, "-")
	for i, word := range words {
		words[i] = cases.Title(language.English).String(word)
	}
	return strings.Join(words, "")
}

func getFlagValidationId(cmdName string, flagName string) string {
	return fmt.Sprintf("CLI2UI_%s_Form_%s_Validation", cmdName, kebabToPascalCase(flagName))
}

func getFlagFieldId(cmdName string, flagName string) string {
	return fmt.Sprintf("CLI2UI_%s_Form_%s_Field", cmdName, kebabToPascalCase(flagName))
}

func getFlagInputId(cmdName string, flagName string) string {
	return fmt.Sprintf("CLI2UI_%s_Form_%s_Input", cmdName, kebabToPascalCase(flagName))
}

func getFlagSelectorId(cmdName string) string {
	return fmt.Sprintf("%sFormFlagSelector", cmdName)
}

type FlagArgType string

const (
	FlagArgTypeString  FlagArgType = "string"
	FlagArgTypeNumber  FlagArgType = "number"
	FlagArgTypeArray   FlagArgType = "array"
	FlagArgTypeBoolean FlagArgType = "boolean"
	FlagArgTypeEnum    FlagArgType = "enum"
)

type CLIJson struct {
	Name        string    `json:"name" yaml:"name"` // CLI name, e.g., "ping", "kubectl", etc.
	Description string    `json:"description" yaml:"description"`
	Help        string    `json:"help" yaml:"help"`                           // CLI help message
	Version     string    `json:"version,omitempty" yaml:"version,omitempty"` // CLI version number, optional field
	Commands    []Command `json:"actions" yaml:"actions"`                     // List of supported actions by the CLI
}

type Command struct {
	Name        string      `json:"name" yaml:"name"` // Name of the action, e.g., "apply", "delete", etc.
	Description string      `json:"description" yaml:"description"`
	Flags       []FlagOrArg `json:"flags" yaml:"flags"`                           // List of flag parameters for the action
	Examples    []string    `json:"examples,omitempty" yaml:"examples,omitempty"` // List of example usages for the action, optional field
	Args        []FlagOrArg `json:"args,omitempty" yaml:"args,omitempty"`         // List of positional arguments for the action, optional field
}

type FlagOrArg struct {
	Name        string      `json:"name" yaml:"name"` // Name of the flag parameter
	Description string      `json:"description" yaml:"description"`
	Short       string      `json:"short,omitempty" yaml:"short,omitempty"`       // Short option for the flag parameter, optional field
	Type        FlagArgType `json:"type,omitempty" yaml:"type,omitempty"`         // Type of the flag parameter, optional field
	Required    bool        `json:"required,omitempty" yaml:"required,omitempty"` // Whether the flag parameter is required, optional field
	Default     string      `json:"default,omitempty" yaml:"default,omitempty"`   // Default value for the flag parameter, optional field
	Options     []string    `json:"options,omitempty" yaml:"options,omitempty"`   // Enum options when flag type is enum
}

func structToMap(inputStruct interface{}) map[string]interface{} {
	data, err := json.Marshal(inputStruct)
	if err != nil {
		log.Fatal("Error marshaling struct to JSON", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal("Error unmarshal JSON to map", err)
	}

	return result
}

type CSSProperties struct {
	PaddingLeft     string `json:"paddingLeft"`
	PaddingRight    string `json:"paddingRight"`
	PaddingTop      string `json:"paddingTop"`
	PaddingBottom   string `json:"paddingBottom"`
	MarginRight     string `json:"marginRight"`
	BackgroundColor string `json:"backgroundColor"`
	MarginBottom    string `json:"marginBottom"`
	Width           string `json:"width"`
	FontSize        string `json:"fontSize"`
	FontWeight      string `json:"fontWeight"`
	Color           string `json:"color"`
}

type Style struct {
	StyleSlot     string        `json:"styleSlot"`
	Style         string        `json:"style"`
	CSSProperties CSSProperties `json:"cssProperties"`
}

type LayoutProperty struct {
	ShowHeader              bool `json:"showHeader"`
	ShowSideBar             bool `json:"showSideBar"`
	SidebarCollapsible      bool `json:"sidebarCollapsible"`
	SidebarDefaultCollapsed bool `json:"sidebarDefaultCollapsed"`
	ShowFooter              bool `json:"showFooter"`
}

func genLayout() []sunmao.ComponentSchema {
	return []sunmao.ComponentSchema{
		{
			Id:   "Layout",
			Type: "arco/v1/layout",
			Properties: structToMap(LayoutProperty{
				ShowHeader:              true,
				ShowSideBar:             false,
				SidebarCollapsible:      false,
				SidebarDefaultCollapsed: false,
				ShowFooter:              false,
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v1/style",
					Properties: structToMap(struct {
						Styles []Style `json:"styles"`
					}{
						Styles: []Style{
							{
								StyleSlot: "layout",
								Style:     "",
								CSSProperties: CSSProperties{
									PaddingLeft:     "32px",
									PaddingRight:    "32px",
									PaddingTop:      "16px",
									PaddingBottom:   "16px",
									MarginRight:     "",
									BackgroundColor: "",
									MarginBottom:    "",
								},
							},
						},
					}),
				},
			},
		},
	}
}

type Container struct {
	ID   string `json:"id"`
	Slot string `json:"slot"`
}

type SlotProperties struct {
	Container   Container   `json:"container"`
	IfCondition interface{} `json:"ifCondition"`
}

type EventHandler struct {
	Type        string `json:"type"`
	ComponentId string `json:"componentId"`
	Method      struct {
		Name       string      `json:"name"`
		Parameters interface{} `json:"parameters"`
	} `json:"method"`
	Wait struct {
		Type string `json:"type"`
		Time int    `json:"time"`
	} `json:"wait"`
	Disabled interface{} `json:"disabled"`
}

type EventProperties struct {
	Handlers []EventHandler `json:"handlers"`
}

type ButtonProperties struct {
	Type     string      `json:"type"`
	Status   string      `json:"status"`
	Long     bool        `json:"long"`
	Size     string      `json:"size"`
	Disabled interface{} `json:"disabled"`
	Loading  interface{} `json:"loading"`
	Shape    string      `json:"shape"`
	Text     string      `json:"text"`
}

type TextProperties struct {
	Text string `json:"text"`
}

type StackProperties struct {
	Spacing   int    `json:"spacing"`
	Direction string `json:"direction"`
	Align     string `json:"align"`
	Wrap      bool   `json:"wrap"`
	Justify   string `json:"justify"`
}

func genHeader(raw CLIJson) []sunmao.ComponentSchema {
	return []sunmao.ComponentSchema{
		{
			Id:   "Header",
			Type: "core/v1/stack",
			Properties: structToMap(StackProperties{
				Spacing:   12,
				Direction: "horizontal",
				Align:     "center",
				Wrap:      false,
				Justify:   "flex-start",
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v1/style",
					Properties: structToMap(struct {
						Styles []Style `json:"styles"`
					}{
						Styles: []Style{
							{
								StyleSlot: "content",
								Style:     "",
								CSSProperties: CSSProperties{
									Width: "100%",
								},
							},
						},
					}),
				},
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   "Layout",
							Slot: "header",
						},
						IfCondition: true,
					}),
				},
			},
		},
		{
			Id:   "HeaderText",
			Type: "core/v2/text",
			Properties: structToMap(TextProperties{
				Text: raw.Name,
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   "Header",
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
				{
					Type: "core/v1/style",
					Properties: structToMap(struct {
						Styles []Style `json:"styles"`
					}{
						Styles: []Style{
							{
								StyleSlot: "content",
								Style:     "",
								CSSProperties: CSSProperties{
									FontSize:   "32px",
									FontWeight: "700",
								},
							},
						},
					}),
				},
			},
		},
		{
			Id:   "HelpBtn",
			Type: "arco/v1/button",
			Properties: structToMap(ButtonProperties{
				Type:     "default",
				Status:   "default",
				Long:     false,
				Size:     "default",
				Disabled: false,
				Loading:  false,
				Shape:    "square",
				Text:     "Help",
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   "Header",
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
				{
					Type: "core/v1/event",
					Properties: structToMap(EventProperties{
						Handlers: []EventHandler{
							{
								Type:        "onClick",
								ComponentId: "HelpModal",
								Method: struct {
									Name       string      `json:"name"`
									Parameters interface{} `json:"parameters"`
								}{
									Name:       "openModal",
									Parameters: make(map[string]interface{}),
								},
								Wait: struct {
									Type string `json:"type"`
									Time int    `json:"time"`
								}{
									Type: "debounce",
									Time: 0,
								},
								Disabled: false,
							},
						},
					}),
				},
			},
		},
	}
}

type ModalProperties struct {
	Title          string `json:"title"`
	Mask           bool   `json:"mask"`
	Simple         bool   `json:"simple"`
	OkText         string `json:"okText"`
	CancelText     string `json:"cancelText"`
	Closable       bool   `json:"closable"`
	MaskClosable   bool   `json:"maskClosable"`
	ConfirmLoading bool   `json:"confirmLoading"`
	DefaultOpen    bool   `json:"defaultOpen"`
	UnmountOnExit  bool   `json:"unmountOnExit"`
}

type TextDisplayProperties struct {
	Text   string `json:"text"`
	Format string `json:"format"`
}

type Method struct {
	Name       string      `json:"name"`
	Parameters interface{} `json:"parameters"`
}

type Wait struct {
	Type string `json:"type"`
	Time int    `json:"time"`
}

func genHelpModal(raw CLIJson) []sunmao.ComponentSchema {

	return []sunmao.ComponentSchema{
		{
			Id:   "HelpModal",
			Type: "arco/v1/modal",
			Properties: structToMap(ModalProperties{
				Title:          "Help 信息",
				Mask:           true,
				Simple:         false,
				OkText:         "confirm",
				CancelText:     "cancel",
				Closable:       true,
				MaskClosable:   true,
				ConfirmLoading: false,
				DefaultOpen:    false,
				UnmountOnExit:  true,
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v1/style",
					Properties: structToMap(struct {
						Styles []Style `json:"styles"`
					}{
						Styles: []Style{
							{
								StyleSlot: "content",
								CSSProperties: CSSProperties{
									Width: "800px",
								},
							},
						},
					}),
				},
			},
		}, {
			Id:   "HelpModalCancelBtn",
			Type: "arco/v1/button",
			Properties: structToMap(ButtonProperties{
				Type:     "default",
				Status:   "default",
				Long:     false,
				Size:     "default",
				Disabled: false,
				Loading:  false,
				Shape:    "square",
				Text:     "关闭",
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   "HelpModal",
							Slot: "footer",
						},
						IfCondition: true,
					}),
				},
				{
					Type: "core/v1/event",
					Properties: structToMap(EventProperties{
						Handlers: []EventHandler{
							{
								Type:        "onClick",
								ComponentId: "HelpModal",
								Method: Method{
									Name:       "closeModal",
									Parameters: map[string]interface{}{},
								},
								Wait: Wait{
									Type: "debounce",
									Time: 0,
								},
								Disabled: false,
							},
						},
					}),
				},
			},
		}, {
			Id:   "HelpModalContent",
			Type: "core/v1/stack",
			Properties: structToMap(StackProperties{
				Spacing:   12,
				Direction: "horizontal",
				Align:     "auto",
				Wrap:      false,
				Justify:   "flex-start",
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   "HelpModal",
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
				{
					Type: "core/v1/style",
					Properties: structToMap(struct {
						Styles []Style `json:"styles"`
					}{
						Styles: []Style{
							{
								StyleSlot: "content",
								Style:     "overflow:auto; box-sizing: border-box;",
								CSSProperties: CSSProperties{
									BackgroundColor: "#333",
									Color:           "white",
									Width:           "100%",
									PaddingTop:      "16px",
									PaddingLeft:     "16px",
									PaddingBottom:   "16px",
									PaddingRight:    "16px",
								},
							},
						},
					}),
				},
			},
		}, {
			Id:   "HelpInfo",
			Type: "cli2ui/v1/TextDisplay",
			Properties: structToMap(TextDisplayProperties{
				Text:   raw.Help,
				Format: "code",
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   "HelpModalContent",
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
				{
					Type: "core/v1/style",
					Properties: structToMap(struct {
						Styles []Style `json:"styles"`
					}{
						Styles: []Style{
							{
								StyleSlot: "content",
								Style:     "height:300px;",
							},
						},
					}),
				},
			},
		}}
}

type TabPropertiesTab struct {
	Title         string `json:"title"`
	Hidden        bool   `json:"hidden"`
	DestroyOnHide bool   `json:"destroyOnHide"`
}

type TabProperties struct {
	Type                          string             `json:"type"`
	DefaultActiveTab              int                `json:"defaultActiveTab"`
	TabPosition                   string             `json:"tabPosition"`
	Size                          string             `json:"size"`
	UpdateWhenDefaultValueChanges bool               `json:"updateWhenDefaultValueChanges"`
	Tabs                          []TabPropertiesTab `json:"tabs"`
}

func genCmdTab(raw CLIJson) []sunmao.ComponentSchema {
	tabs := []TabPropertiesTab{}
	for _, item := range raw.Commands {
		tab := TabPropertiesTab{
			Title:         item.Name,
			Hidden:        false,
			DestroyOnHide: true,
		}
		tabs = append(tabs, tab)
	}

	return []sunmao.ComponentSchema{
		{
			Id:   "CmdTabs",
			Type: "arco/v1/tabs",
			Properties: structToMap(TabProperties{
				Type:                          "line",
				DefaultActiveTab:              0,
				TabPosition:                   "top",
				Size:                          "default",
				UpdateWhenDefaultValueChanges: false,
				Tabs:                          tabs,
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   "Layout",
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
			},
		},
	}
}

type SwitchProperty struct {
	DefaultChecked                bool   `json:"defaultChecked"`
	Disabled                      string `json:"disabled"`
	Loading                       bool   `json:"loading"`
	Type                          string `json:"type"`
	Size                          string `json:"size"`
	UpdateWhenDefaultValueChanges bool   `json:"updateWhenDefaultValueChanges"`
}

type SelectProperty struct {
	AllowClear                    bool           `json:"allowClear"`
	Multiple                      bool           `json:"multiple"`
	AllowCreate                   bool           `json:"allowCreate"`
	Bordered                      bool           `json:"bordered"`
	DefaultValue                  interface{}    `json:"defaultValue"`
	Disabled                      string         `json:"disabled"`
	LabelInValue                  bool           `json:"labelInValue"`
	Loading                       bool           `json:"loading"`
	ShowSearch                    bool           `json:"showSearch"`
	UnmountOnExit                 bool           `json:"unmountOnExit"`
	ShowTitle                     bool           `json:"showTitle"`
	Options                       []SelectOption `json:"options"`
	Placeholder                   string         `json:"placeholder"`
	Size                          string         `json:"size"`
	Error                         bool           `json:"error"`
	UpdateWhenDefaultValueChanges bool           `json:"updateWhenDefaultValueChanges"`
	AutoFixPosition               bool           `json:"autoFixPosition"`
	AutoAlignPopupMinWidth        bool           `json:"autoAlignPopupMinWidth"`
	AutoAlignPopupWidth           bool           `json:"autoAlignPopupWidth"`
	AutoFitPosition               bool           `json:"autoFitPosition"`
	Position                      string         `json:"position"`
	MountToBody                   bool           `json:"mountToBody"`
}

type SelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type TextInputProperty struct {
	AllowClear                    bool   `json:"allowClear"`
	Disabled                      string `json:"disabled"`
	ReadOnly                      bool   `json:"readOnly"`
	DefaultValue                  string `json:"defaultValue"`
	UpdateWhenDefaultValueChanges bool   `json:"updateWhenDefaultValueChanges"`
	Placeholder                   string `json:"placeholder"`
	Error                         string `json:"error"`
	Size                          string `json:"size"`
}

type InputSlotStyle struct {
	StyleSlot     string `json:"styleSlot"`
	Style         string `json:"style"`
	CSSProperties struct {
		MarginBottom string `json:"marginBottom"`
	} `json:"cssProperties"`
}

type TextInputTrait struct {
	Type       string                `json:"type"`
	Properties TextInputTraitContent `json:"properties"`
}

type TextInputTraitContent struct {
	Container   Container `json:"container"`
	IfCondition bool      `json:"ifCondition"`
}

func genCmdFormFiledInput(cmdName string, flag FlagOrArg) sunmao.ComponentSchema {
	fieldId := getFlagFieldId(cmdName, flag.Name)
	inputId := getFlagInputId(cmdName, flag.Name)
	validationId := getFlagValidationId(cmdName, flag.Name)
	options := flag.Options

	switch flag.Type {
	case FlagArgTypeNumber:
		return sunmao.ComponentSchema{
			Id:   inputId,
			Type: "arco/v1/numberInput",
			Properties: map[string]interface{}{
				"defaultValue":                  1,
				"disabled":                      "{{" + cmdName + "FormState.data.isExecuting}}",
				"placeholder":                   "please input",
				"error":                         false,
				"size":                          "default",
				"buttonMode":                    false,
				"min":                           0,
				"max":                           99,
				"readOnly":                      false,
				"step":                          1,
				"precision":                     0,
				"updateWhenDefaultValueChanges": false,
			},
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   fieldId,
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
			},
		}
	case FlagArgTypeArray:
		return sunmao.ComponentSchema{
			Id:   inputId,
			Type: "cli2ui/v1/arrayInput",
			Properties: map[string]interface{}{
				"value":       []string{""},
				"type":        "string",
				"placeholder": "please input",
				"disabled":    "{{" + cmdName + "FormState.data.isExecuting}}",
			},
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   fieldId,
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
			},
		}
	case FlagArgTypeBoolean:
		return sunmao.ComponentSchema{
			Id:   inputId,
			Type: "arco/v1/switch",
			Properties: structToMap(SwitchProperty{
				DefaultChecked:                false,
				Disabled:                      fmt.Sprintf("{{%sFormState.data.isExecuting}}", cmdName),
				Loading:                       false,
				Type:                          "circle",
				Size:                          "default",
				UpdateWhenDefaultValueChanges: false,
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   fieldId,
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
			},
		}
	case FlagArgTypeEnum:
		selectOptions := make([]SelectOption, 0)
		for _, op := range options {
			selectOptions = append(selectOptions, SelectOption{
				Label: op,
				Value: op,
			})
		}
		return sunmao.ComponentSchema{
			Id:   inputId,
			Type: "arco/v1/select",
			Properties: structToMap(SelectProperty{
				AllowClear:                    false,
				Multiple:                      false,
				AllowCreate:                   false,
				Bordered:                      true,
				DefaultValue:                  nil,
				Disabled:                      fmt.Sprintf("{{%sFormState.data.isExecuting}}", cmdName),
				LabelInValue:                  false,
				Loading:                       false,
				ShowSearch:                    false,
				UnmountOnExit:                 true,
				ShowTitle:                     false,
				Options:                       selectOptions,
				Placeholder:                   "Select an option",
				Size:                          "default",
				Error:                         false,
				UpdateWhenDefaultValueChanges: false,
				AutoFixPosition:               false,
				AutoAlignPopupMinWidth:        false,
				AutoAlignPopupWidth:           true,
				AutoFitPosition:               false,
				Position:                      "bottom",
				MountToBody:                   true,
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   fieldId,
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
			},
		}
	default:
		return sunmao.ComponentSchema{
			Id:   inputId,
			Type: "arco/v1/input",
			Properties: structToMap(TextInputProperty{
				AllowClear:                    false,
				Disabled:                      fmt.Sprintf("{{%sFormState.data.isExecuting}}", cmdName),
				ReadOnly:                      false,
				DefaultValue:                  "",
				UpdateWhenDefaultValueChanges: false,
				Placeholder:                   "please input",
				Error:                         fmt.Sprintf("{{%sForm.validatedResult.%s.isInvalid}}", cmdName, validationId),
				Size:                          "default",
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container:   Container{ID: fieldId, Slot: "content"},
						IfCondition: true,
					}),
				},
				{
					Type: "core/v1/event",
					Properties: structToMap(EventProperties{
						Handlers: []EventHandler{
							{
								Type:        "onBlur",
								ComponentId: fmt.Sprintf("%sForm", cmdName),
								Method: Method{
									Name:       "validateAllFields",
									Parameters: map[string]interface{}{},
								},
								Wait: Wait{
									Type: "debounce",
									Time: 0,
								},
								Disabled: false,
							},
						},
					}),
				},
				{
					Type: "core/v1/style",
					Properties: structToMap(struct {
						Styles []Style `json:"styles"`
					}{
						Styles: []Style{
							{
								StyleSlot: "input",
								CSSProperties: CSSProperties{
									MarginBottom: "6px",
								},
							},
						},
					}),
				},
			},
		}
	}
}

type FormControlProperty struct {
	Label      Label  `json:"label"`
	Layout     string `json:"layout"`
	Required   bool   `json:"required"`
	Hidden     bool   `json:"hidden"`
	Extra      string `json:"extra"`
	ErrorMsg   string `json:"errorMsg"`
	LabelAlign string `json:"labelAlign"`
	Colon      bool   `json:"colon"`
	Help       string `json:"help"`
	LabelCol   Column `json:"labelCol"`
	WrapperCol Column `json:"wrapperCol"`
}

type Label struct {
	Format string `json:"format"`
	Raw    string `json:"raw"`
}

type Column struct {
	Span   int `json:"span"`
	Offset int `json:"offset"`
}

func genCmdFormFields(cmdName string, flags []FlagOrArg) []sunmao.ComponentSchema {
	formFieldComponents := make([]sunmao.ComponentSchema, 0)
	flagSelectorId := getFlagSelectorId(cmdName)

	for _, flag := range flags {
		argFieldId := getFlagFieldId(cmdName, flag.Name)

		formFieldComponents = append(formFieldComponents, sunmao.ComponentSchema{
			Id:   argFieldId,
			Type: "arco/v1/formControl",
			Properties: structToMap(FormControlProperty{
				Label: Label{
					Format: "plain",
					Raw:    flag.Name,
				},
				Layout:     "horizontal",
				Required:   flag.Required,
				Hidden:     false,
				Extra:      "",
				ErrorMsg:   "",
				LabelAlign: "left",
				Colon:      false,
				Help:       "",
				LabelCol: Column{
					Span:   6,
					Offset: 0,
				},
				WrapperCol: Column{
					Span:   18,
					Offset: 0,
				},
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   fmt.Sprintf("%sForm", cmdName),
							Slot: "content",
						},
						IfCondition: fmt.Sprintf("{{ %s.value.some(item => item === \"%s\") }}", flagSelectorId, flag.Name),
					}),
				},
			},
		})

		formFieldComponents = append(formFieldComponents, genCmdFormFiledInput(cmdName, flag))
	}

	return formFieldComponents
}

type StateTrait struct {
	Type       string            `json:"type"`
	Properties StateTraitContent `json:"properties"`
}

type StateTraitContent struct {
	Key          string `json:"key"`
	InitialValue string `json:"initialValue"`
}

type TransformerTrait struct {
	Type       string                `json:"type"`
	Properties TransformerProperties `json:"properties"`
}

type TransformerProperties struct {
	Value string `json:"value"`
}

type CheckboxMenuProperties struct {
	Value   []string         `json:"value"`
	Text    string           `json:"text"`
	Options []CheckboxOption `json:"options"`
}

type CheckboxOption struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled"`
}

type ValidationTrait struct {
	Type       string               `json:"type"`
	Properties ValidationProperties `json:"properties"`
}

type ValidationProperties struct {
	Validators []Validator `json:"validators"`
}

type Validator struct {
	Name  string        `json:"name"`
	Value string        `json:"value"`
	Rules []interface{} `json:"rules"`
}

type DividerProperties struct {
	Type        string `json:"type"`
	Orientation string `json:"orientation"`
}

type CollapseProperties struct {
	DefaultActiveKey              []string         `json:"defaultActiveKey"`
	Options                       []CollapseOption `json:"options"`
	UpdateWhenDefaultValueChanges bool             `json:"updateWhenDefaultValueChanges"`
	Accordion                     bool             `json:"accordion"`
	ExpandIconPosition            string           `json:"expandIconPosition"`
	Bordered                      bool             `json:"bordered"`
	DestroyOnHide                 bool             `json:"destroyOnHide"`
	LazyLoad                      bool             `json:"lazyLoad"`
}

type CollapseOption struct {
	Key            string `json:"key"`
	Header         string `json:"header"`
	Disabled       bool   `json:"disabled"`
	ShowExpandIcon bool   `json:"showExpandIcon"`
}

type ResultProperties struct {
	Data string `json:"data"`
}

type TerminalProperties struct {
	Text string `json:"text"`
}

func genCmdInnerTabs(raw CLIJson) []sunmao.ComponentSchema {
	tabs := make([]sunmao.ComponentSchema, 0)

	rawCmd := raw.Commands

	for index, item := range rawCmd {
		flagsAndArgs := make([]struct {
			flagOrArg FlagOrArg
			isArg     bool
		}, 0)
		for idx, foa := range append(item.Flags, item.Args...) {
			isArg := idx < len(item.Flags)
			flagsAndArgs = append(flagsAndArgs, struct {
				flagOrArg FlagOrArg
				isArg     bool
			}{
				flagOrArg: foa,
				isArg:     isArg,
			})
		}

		flagsLen := len(flagsAndArgs)
		transformFlags := make([]string, 0, flagsLen)
		requiredFlags := make([]string, 0)
		checkboxOptions := make([]CheckboxOption, 0, flagsLen)
		validators := make([]Validator, 0, flagsLen)

		for _, flag := range flagsAndArgs {
			flagB, _ := json.Marshal(flag.flagOrArg)
			transformFlags = append(
				transformFlags,
				fmt.Sprintf(`"%s": { "value": %s.value, "isArg": %s, ...%s }`,
					flag.flagOrArg.Name,
					getFlagInputId(item.Name, flag.flagOrArg.Name),
					strconv.FormatBool(flag.isArg),
					string(flagB),
				),
			)

			if flag.flagOrArg.Required {
				requiredFlags = append(requiredFlags, flag.flagOrArg.Name)
			}

			option := CheckboxOption{
				Label:    flag.flagOrArg.Name,
				Value:    flag.flagOrArg.Name,
				Disabled: flag.flagOrArg.Required,
			}
			checkboxOptions = append(checkboxOptions, option)

			inputID := getFlagInputId(item.Name, flag.flagOrArg.Name)
			validationID := getFlagValidationId(item.Name, flag.flagOrArg.Name)
			validation := Validator{
				Name:  validationID,
				Value: fmt.Sprintf("{{%s.value}}", inputID),
				Rules: []interface{}{},
			}

			if flag.flagOrArg.Required {
				validation.Rules = append(validation.Rules, map[string]interface{}{
					"type":     "required",
					"validate": "",
					"error": map[string]interface{}{
						"message": "该字段不可为空。",
					},
					"minLength":     0,
					"maxLength":     0,
					"includeList":   []interface{}{},
					"excludeList":   []interface{}{},
					"min":           0,
					"max":           0,
					"regex":         "",
					"flags":         "",
					"customOptions": make(map[string]interface{}),
				})
			}

			validators = append(validators, validation)
		}

		flagSelectorId := getFlagSelectorId(item.Name)

		formComponents := []sunmao.ComponentSchema{
			{
				Id:         fmt.Sprintf("%sFormState", item.Name),
				Type:       "core/v1/dummy",
				Properties: map[string]interface{}{},
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v1/state",
						Properties: structToMap(StateTraitContent{
							Key:          "data",
							InitialValue: "{{{\n{\n  isExecuting: false,\n}\n}}}",
						}),
					},
				},
			},
			{
				Id:         fmt.Sprintf("%sFormTransformer", item.Name),
				Type:       "core/v1/dummy",
				Properties: map[string]interface{}{},
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v1/transformer",
						Properties: structToMap(TransformerProperties{
							Value: fmt.Sprintf("{{ formatCommand('%s', { %s }, %s.value) }}",
								item.Name,
								strings.Join(
									transformFlags,
									", ",
								),
								flagSelectorId),
						}),
					},
				},
			},
			{
				Id:   fmt.Sprintf("%sTab", item.Name),
				Type: "core/v1/stack",
				Properties: structToMap(StackProperties{
					Spacing:   12,
					Direction: "vertical",
					Align:     "auto",
					Wrap:      false,
					Justify:   "flex-start",
				}),
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v2/slot",
						Properties: structToMap(SlotProperties{
							Container: Container{
								ID:   "CmdTabs",
								Slot: "content",
							},
							IfCondition: fmt.Sprintf("{{CmdTabs.activeTab === %d}}", index),
						}),
					},
					{
						Type: "core/v1/style",
						Properties: structToMap(struct {
							Styles []Style `json:"styles"`
						}{
							Styles: []Style{
								{
									StyleSlot: "content",
									Style:     "",
									CSSProperties: CSSProperties{
										Width: "100%",
									},
								},
							},
						}),
					},
				},
			},
			{
				Id:   fmt.Sprintf("%sFormActionHeader", item.Name),
				Type: "core/v1/stack",
				Properties: structToMap(StackProperties{
					Spacing:   12,
					Direction: "horizontal",
					Align:     "auto",
					Wrap:      false,
					Justify:   "flex-end",
				}),
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v2/slot",
						Properties: structToMap(SlotProperties{
							Container: Container{
								ID:   fmt.Sprintf("%sTab", item.Name),
								Slot: "content",
							},
							IfCondition: true,
						}),
					},
				},
			},
			{
				Id:   flagSelectorId,
				Type: "cli2ui/v1/checkboxMenu",
				Properties: structToMap(CheckboxMenuProperties{
					Value:   requiredFlags,
					Text:    "flags",
					Options: checkboxOptions,
				}),
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v2/slot",
						Properties: structToMap(SlotProperties{
							Container: Container{
								ID:   fmt.Sprintf("%sFormActionHeader", item.Name),
								Slot: "content",
							},
							IfCondition: true,
						}),
					},
				},
			},
			{
				Id:   fmt.Sprintf("%sForm", item.Name),
				Type: "core/v1/stack",
				Properties: structToMap(StackProperties{
					Spacing:   12,
					Direction: "vertical",
					Align:     "auto",
					Wrap:      false,
					Justify:   "flex-start",
				}),
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v2/slot",
						Properties: structToMap(SlotProperties{
							Container: Container{
								ID:   fmt.Sprintf("%sTab", item.Name),
								Slot: "content",
							},
							IfCondition: true,
						}),
					},
					{
						Type: "core/v1/validation",
						Properties: structToMap(ValidationProperties{
							Validators: validators,
						}),
					},
				},
			},
		}
		resultComponents := []sunmao.ComponentSchema{{
			Id:   fmt.Sprintf("%sTabDivider", item.Name),
			Type: "arco/v1/divider",
			Properties: structToMap(DividerProperties{
				Type:        "horizontal",
				Orientation: "center",
			}),
			Traits: []sunmao.TraitSchema{
				{
					Type: "core/v2/slot",
					Properties: structToMap(SlotProperties{
						Container: Container{
							ID:   fmt.Sprintf("%sTab", item.Name),
							Slot: "content",
						},
						IfCondition: true,
					}),
				},
			},
		},
			{
				Id:   fmt.Sprintf("%sTabResult", item.Name),
				Type: "arco/v1/collapse",
				Properties: structToMap(CollapseProperties{
					DefaultActiveKey: []string{"0"},
					Options: []CollapseOption{
						{
							Key:            "0",
							Header:         "执行结果\n",
							Disabled:       false,
							ShowExpandIcon: true,
						},
					},
					UpdateWhenDefaultValueChanges: false,
					Accordion:                     false,
					ExpandIconPosition:            "left",
					Bordered:                      false,
					DestroyOnHide:                 false,
					LazyLoad:                      true,
				}),
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v2/slot",
						Properties: structToMap(SlotProperties{
							Container: Container{
								ID:   fmt.Sprintf("%sTab", item.Name),
								Slot: "content",
							},
							IfCondition: true,
						}),
					},
					{
						Type: "core/v1/style",
						Properties: structToMap(struct {
							Styles []Style `json:"styles"`
						}{
							Styles: []Style{
								{
									StyleSlot: "content",
									Style:     ".arco-collapse-item-content-box {\nbackground: white;\n}",
									CSSProperties: CSSProperties{
										BackgroundColor: "",
									},
								},
							},
						}),
					},
				},
			},
			{
				Id:   fmt.Sprintf("%sTabResultContent", item.Name),
				Type: "cli2ui/v1/result",
				Properties: structToMap(ResultProperties{
					Data: "",
				}),
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v2/slot",
						Properties: structToMap(SlotProperties{
							Container: Container{
								ID:   fmt.Sprintf("%sTabResult", item.Name),
								Slot: "content",
							},
							IfCondition: true,
						}),
					},
				},
			},
			{
				Id:   fmt.Sprintf("%sTabResultTerminal", item.Name),
				Type: "cli2ui/v1/terminal",
				Properties: structToMap(TerminalProperties{
					Text: "{{ exec.state.stdout }}",
				}),
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v2/slot",
						Properties: structToMap(SlotProperties{
							Container: Container{
								ID:   fmt.Sprintf("%sTabResult", item.Name),
								Slot: "content",
							},
							IfCondition: true,
						}),
					},
				},
			},
			{
				Id:   fmt.Sprintf("%sStopBtn", item.Name),
				Type: "arco/v1/button",
				Properties: structToMap(ButtonProperties{
					Type:     "primary",
					Status:   "default",
					Long:     false,
					Size:     "default",
					Disabled: fmt.Sprintf("{{!%sFormState.data.isExecuting}}", item.Name),
					Loading:  false,
					Shape:    "square",
					Text:     "Stop",
				}),
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v2/slot",
						Properties: structToMap(SlotProperties{
							Container: Container{
								ID:   fmt.Sprintf("%sTab", item.Name),
								Slot: "content",
							},
							IfCondition: true,
						}),
					},
					{
						Type: "core/v1/event",
						Properties: structToMap(EventProperties{
							Handlers: []EventHandler{
								{
									Type:        "onClick",
									ComponentId: fmt.Sprintf("%sFormState", item.Name),
									Method: Method{
										Name: fmt.Sprintf("setValue"),
										Parameters: map[string]interface{}{
											"key":   "data",
											"value": fmt.Sprintf("{{{...%sFormState.data, isExecuting: false}}}", item.Name),
										},
									},
									Wait: Wait{
										Type: "debounce",
										Time: 0,
									},
									Disabled: false,
								},
								{
									Type:        "onClick",
									ComponentId: "$utils",
									Method: Method{
										Name:       "binding/v1/stop",
										Parameters: map[string]interface{}{},
									},
								},
							},
						}),
					},
					{
						Type: "core/v1/style",
						Properties: structToMap(struct {
							Styles []Style `json:"styles"`
						}{
							Styles: []Style{
								{
									StyleSlot: "content",
									Style:     "",
									CSSProperties: CSSProperties{
										Width: "",
									},
								},
							},
						}),
					},
				},
			},
		}

		tabs = append(tabs, formComponents...)
		tabs = append(tabs, genCmdFormFields(item.Name, append(item.Flags, item.Args...))...)
		tabs = append(tabs,
			sunmao.ComponentSchema{
				Id:   fmt.Sprintf("%sButton", item.Name),
				Type: "arco/v1/button",
				Properties: structToMap(ButtonProperties{
					Type:     "primary",
					Status:   "default",
					Long:     false,
					Size:     "default",
					Disabled: false,
					Loading:  fmt.Sprintf("{{%sFormState.data.isExecuting}}", item.Name),
					Shape:    "square",
					Text:     "Run",
				}),
				Traits: []sunmao.TraitSchema{
					{
						Type: "core/v2/slot",
						Properties: structToMap(SlotProperties{
							Container: Container{
								ID:   fmt.Sprintf("%sForm", item.Name),
								Slot: "content",
							},
							IfCondition: true,
						}),
					},
					{
						Type: "core/v1/event",
						Properties: structToMap(EventProperties{
							Handlers: []EventHandler{
								{
									Type:        "onClick",
									ComponentId: fmt.Sprintf("%sForm", item.Name),
									Method: Method{
										Name:       "validateAllFields",
										Parameters: map[string]interface{}{},
									},
									Wait: Wait{
										Type: "debounce",
										Time: 0,
									},
									Disabled: false,
								},
								{
									Type:        "onClick",
									ComponentId: fmt.Sprintf("%sFormState", item.Name),
									Method: Method{
										Name: fmt.Sprintf("setValue"),
										Parameters: map[string]interface{}{
											"key":   "data",
											"value": fmt.Sprintf("{{{...%sFormState.data, isExecuting: true}}}", item.Name),
										},
									},
									Wait: Wait{
										Type: "debounce",
										Time: 0,
									},
									Disabled: fmt.Sprintf("{{%sForm.isInvalid}}", item.Name),
								},
								{
									Type:        "onClick",
									ComponentId: "$utils",
									Method: Method{
										Name:       "binding/v1/run",
										Parameters: fmt.Sprintf("{{ %sFormTransformer.value }}", item.Name),
									},
								},
							},
						}),
					},
				},
			})
		tabs = append(tabs, resultComponents...)
	}

	return tabs
}

func genSchemaComponents(raw CLIJson) []sunmao.ComponentSchema {
	components := make([]sunmao.ComponentSchema, 0)

	components = append(components, genLayout()...)
	components = append(components, genHeader(raw)...)
	components = append(components, genCmdTab(raw)...)
	components = append(components, genCmdInnerTabs(raw)...)
	components = append(components, genHelpModal(raw)...)

	return components
}
