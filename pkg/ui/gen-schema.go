package ui

import (
	"encoding/json"
	"fmt"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/arco"
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

func (u *UI) genHelpModal(raw CLIJson) []sunmao.BaseComponentBuilder {
	return []sunmao.BaseComponentBuilder{
		u.b.NewModal().Id("HelpModal").Properties(structToMap(ModalProperties{
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
		})).Style("content", "width: 800px").Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {
				u.b.NewStack().Id("HelpModalContent").Properties(structToMap(StackProperties{
					Spacing:   12,
					Direction: "horizontal",
					Align:     "auto",
					Wrap:      false,
					Justify:   "flex-start",
				})).Style("content", `
					overflow:auto; box-sizing: border-box; background-color: #333;
					color: white; width: 100%; padding: 16px;
				`),
				u.bb.NewTextDisplay().Id("HelpInfo").Content(TextDisplayProperties{
					Text:   raw.Help,
					Format: "code",
				}).Style("content", "height: 300px;"),
			},
			"footer": {
				u.b.NewButton().Id("HelpModalCancelBtn").Properties(structToMap(ButtonProperties{
					Type:     "default",
					Status:   "default",
					Long:     false,
					Size:     "default",
					Disabled: false,
					Loading:  false,
					Shape:    "square",
					Text:     "关闭",
				})).Event([]sunmao.EventHandler{
					{
						Type:        "onClick",
						ComponentId: "HelpModal",
						Method: sunmao.EventMethod{
							Name:       "closeModal",
							Parameters: map[string]interface{}{},
						},
					},
				}),
			},
		}),
	}

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

func (u *UI) genCmdTab(raw CLIJson) []sunmao.BaseComponentBuilder {
	c := u.b.NewTabs().Id("CmdTabs").Properties(structToMap(TabProperties{
		Type:                          "line",
		DefaultActiveTab:              0,
		TabPosition:                   "top",
		Size:                          "default",
		UpdateWhenDefaultValueChanges: false,
		Tabs:                          []TabPropertiesTab{},
	}))

	for _, item := range raw.Commands {
		c.Tab(&arco.ArcoTabsTab{
			Title:         item.Name,
			Hidden:        false,
			DestroyOnHide: true,
		})
	}

	return []sunmao.BaseComponentBuilder{c}
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

func (u *UI) genCmdFormFiledInput(cmdName string, flag FlagOrArg) []sunmao.BaseComponentBuilder {
	inputId := getFlagInputId(cmdName, flag.Name)
	validationId := getFlagValidationId(cmdName, flag.Name)
	options := flag.Options

	switch flag.Type {
	case FlagArgTypeNumber:
		return []sunmao.BaseComponentBuilder{
			u.b.NewNumberInput().Id(inputId).Properties(map[string]interface{}{
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
			}),
		}
	case FlagArgTypeArray:
		return []sunmao.BaseComponentBuilder{
			u.bb.NewArrayInput().Id(inputId).Properties(map[string]interface{}{
				"value":       []string{""},
				"type":        "string",
				"placeholder": "please input",
				"disabled":    "{{" + cmdName + "FormState.data.isExecuting}}",
			}),
		}
	case FlagArgTypeBoolean:
		return []sunmao.BaseComponentBuilder{
			u.b.NewSwitch().Id(inputId).Properties(structToMap(SwitchProperty{
				DefaultChecked:                false,
				Disabled:                      fmt.Sprintf("{{%sFormState.data.isExecuting}}", cmdName),
				Loading:                       false,
				Type:                          "circle",
				Size:                          "default",
				UpdateWhenDefaultValueChanges: false,
			})),
		}
	case FlagArgTypeEnum:
		selectOptions := make([]SelectOption, 0)
		for _, op := range options {
			selectOptions = append(selectOptions, SelectOption{
				Label: op,
				Value: op,
			})
		}
		return []sunmao.BaseComponentBuilder{
			u.b.NewSelect().Id(inputId).Properties(structToMap(SelectProperty{
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
			})),
		}
	default:
		return []sunmao.BaseComponentBuilder{
			u.b.NewInput().Id(inputId).Properties(structToMap(TextInputProperty{
				AllowClear:                    false,
				Disabled:                      fmt.Sprintf("{{%sFormState.data.isExecuting}}", cmdName),
				ReadOnly:                      false,
				DefaultValue:                  "",
				UpdateWhenDefaultValueChanges: false,
				Placeholder:                   "please input",
				Error:                         fmt.Sprintf("{{%sForm.validatedResult.%s.isInvalid}}", cmdName, validationId),
				Size:                          "default",
			})).Event([]sunmao.EventHandler{
				{
					Type:        "onBlur",
					ComponentId: fmt.Sprintf("%sForm", cmdName),
					Method: sunmao.EventMethod{
						Name:       "validateAllFields",
						Parameters: map[string]interface{}{},
					},
				},
			}).Style("input", "margin-bottom: 6px;"),
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

func (u *UI) genCmdFormFields(cmdName string, flags []FlagOrArg) []sunmao.BaseComponentBuilder {
	formFieldComponents := make([]sunmao.BaseComponentBuilder, 0)
	flagSelectorId := getFlagSelectorId(cmdName)

	for _, flag := range flags {
		argFieldId := getFlagFieldId(cmdName, flag.Name)

		u.b.Component(u.b.NewFormControl().Id(argFieldId).
			Properties(structToMap(FormControlProperty{
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
			})).
			Slot(sunmao.Container{
				ID:   fmt.Sprintf("%sForm", cmdName),
				Slot: "content",
			}, fmt.Sprintf("{{ %s.value.some(item => item === \"%s\") }}", flagSelectorId, flag.Name)).
			Children(map[string][]sunmao.BaseComponentBuilder{
				"content": u.genCmdFormFiledInput(cmdName, flag),
			}))
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

func (u *UI) genCmdInnerTabs(raw CLIJson) []sunmao.BaseComponentBuilder {
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

		u.b.Component(u.b.NewState("data", map[string]interface{}{
			"isExecuting": false,
		}).Id(fmt.Sprintf("%sFormState", item.Name)))

		u.b.Component(u.b.NewTransformer(fmt.Sprintf("{{ formatCommand('%s', { %s }, %s.value) }}",
			item.Name,
			strings.Join(
				transformFlags,
				", ",
			),
			flagSelectorId)).Id(fmt.Sprintf("%sFormTransformer", item.Name)))

		u.b.Component(u.b.NewStack().Id(fmt.Sprintf("%sTab", item.Name)).
			Properties(structToMap(StackProperties{
				Spacing:   12,
				Direction: "vertical",
				Align:     "auto",
				Wrap:      false,
				Justify:   "flex-start",
			})).
			Slot(sunmao.Container{
				ID:   "CmdTabs",
				Slot: "content",
			}, fmt.Sprintf("{{CmdTabs.activeTab === %d}}", index)).
			Style("content", "width: 100%").
			Children(map[string][]sunmao.BaseComponentBuilder{
				"content": {
					u.b.NewStack().Id(fmt.Sprintf("%sFormActionHeader", item.Name)).
						Properties(structToMap(StackProperties{
							Spacing:   12,
							Direction: "horizontal",
							Align:     "auto",
							Wrap:      false,
							Justify:   "flex-end",
						})).
						Children(map[string][]sunmao.BaseComponentBuilder{
							"content": {
								u.bb.NewCheckboxMenu().Id(flagSelectorId).Properties(structToMap(CheckboxMenuProperties{
									Value:   requiredFlags,
									Text:    "flags",
									Options: checkboxOptions,
								})),
							},
						}),

					u.b.NewStack().Id(fmt.Sprintf("%sForm", item.Name)).Properties(structToMap(StackProperties{
						Spacing:   12,
						Direction: "vertical",
						Align:     "auto",
						Wrap:      false,
						Justify:   "flex-start",
					})).Trait(u.b.NewTrait().Type("core/v1/validation").Properties(structToMap(ValidationProperties{
						Validators: validators,
					}))).Children(map[string][]sunmao.BaseComponentBuilder{
						"content": {
							u.b.NewButton().Id(fmt.Sprintf("%sButton", item.Name)).Properties(structToMap(ButtonProperties{
								Type:     "primary",
								Status:   "default",
								Long:     false,
								Size:     "default",
								Disabled: false,
								Loading:  fmt.Sprintf("{{%sFormState.data.isExecuting}}", item.Name),
								Shape:    "square",
								Text:     "Run",
							})).Event([]sunmao.EventHandler{
								{
									Type:        "onClick",
									ComponentId: fmt.Sprintf("%sForm", item.Name),
									Method: sunmao.EventMethod{
										Name:       "validateAllFields",
										Parameters: map[string]interface{}{},
									},
								},
								{
									Type:        "onClick",
									ComponentId: fmt.Sprintf("%sFormState", item.Name),
									Method: sunmao.EventMethod{
										Name: fmt.Sprintf("setValue"),
										Parameters: map[string]interface{}{
											"key":   "data",
											"value": fmt.Sprintf("{{{...%sFormState.data, isExecuting: true}}}", item.Name),
										},
									},
									Disabled: fmt.Sprintf("{{%sForm.isInvalid}}", item.Name),
								},
								{
									Type:        "onClick",
									ComponentId: "$utils",
									Method: sunmao.EventMethod{
										Name:       "binding/v1/run",
										Parameters: fmt.Sprintf("{{ %sFormTransformer.value }}", item.Name),
									},
								},
							}),
						},
					}),

					u.b.NewDivider().Id(fmt.Sprintf("%sTabDivider", item.Name)).Properties(structToMap(DividerProperties{
						Type:        "horizontal",
						Orientation: "center",
					})),

					u.b.NewCollapse().Id(fmt.Sprintf("%sTabResult", item.Name)).
						Properties(structToMap(CollapseProperties{
							DefaultActiveKey: []string{"0"},
							Options: []CollapseOption{
								{
									Key:            "0",
									Header:         "执行结果",
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
						})).
						Style("content", ".arco-collapse-item-content-box { background: white; }").
						Children(map[string][]sunmao.BaseComponentBuilder{
							"content": {
								u.bb.NewResult().Id(fmt.Sprintf("%sTabResultContent", item.Name)).Data(""),

								u.bb.NewTerminal().Id(fmt.Sprintf("%sTabResultTerminal", item.Name)).Text("{{ exec.state.stdout }}"),
							},
						}),

					u.b.NewButton().Id(fmt.Sprintf("%sStopBtn", item.Name)).Properties(structToMap(ButtonProperties{
						Type:     "primary",
						Status:   "default",
						Long:     false,
						Size:     "default",
						Disabled: fmt.Sprintf("{{!%sFormState.data.isExecuting}}", item.Name),
						Loading:  false,
						Shape:    "square",
						Text:     "Stop",
					})).Event([]sunmao.EventHandler{
						{
							Type:        "onClick",
							ComponentId: fmt.Sprintf("%sFormState", item.Name),
							Method: sunmao.EventMethod{
								Name: fmt.Sprintf("setValue"),
								Parameters: map[string]interface{}{
									"key":   "data",
									"value": fmt.Sprintf("{{{...%sFormState.data, isExecuting: false}}}", item.Name),
								},
							},
						},
						{
							Type:        "onClick",
							ComponentId: "$utils",
							Method: sunmao.EventMethod{
								Name:       "binding/v1/stop",
								Parameters: map[string]interface{}{},
							},
						},
					}),
				},
			}))

		u.genCmdFormFields(item.Name, append(item.Flags, item.Args...))
	}

	return []sunmao.BaseComponentBuilder{}
}

func (u *UI) genLayout(raw CLIJson) []sunmao.BaseComponentBuilder {
	return []sunmao.BaseComponentBuilder{
		u.b.NewLayout().
			Id("Layout").
			Properties(structToMap(LayoutProperty{
				ShowHeader:              true,
				ShowSideBar:             false,
				SidebarCollapsible:      false,
				SidebarDefaultCollapsed: false,
				ShowFooter:              false,
			})).
			Style("layout", "padding: 16px 32px;").
			Children(map[string][]sunmao.BaseComponentBuilder{
				"header": {
					u.b.NewStack().Id("Header").Properties(structToMap(StackProperties{
						Spacing:   12,
						Direction: "horizontal",
						Align:     "center",
						Wrap:      false,
						Justify:   "flex-start",
					})).
						Style("content", "width: 100%").
						Children(map[string][]sunmao.BaseComponentBuilder{
							"content": u.genHeader(raw),
						}),
				},
				"content": u.genCmdTab(raw),
			})}
}

func (u *UI) genHeader(raw CLIJson) []sunmao.BaseComponentBuilder {
	return []sunmao.BaseComponentBuilder{
		u.b.NewText().Id("HeaderText").Content(raw.Name).
			Style("content", "font-size: 32px; font-weight: 700;"),
		u.b.NewButton().Id("HelpButton").Properties(structToMap(ButtonProperties{
			Shape: "square",
			Text:  "Help",
			Type:  "default",
		})).Event([]sunmao.EventHandler{
			{
				Type:        "onClick",
				ComponentId: "HelpModal",
				Method: sunmao.EventMethod{
					Name:       "openModal",
					Parameters: make(map[string]interface{}),
				},
			},
		}),
	}
}

func (u *UI) genSchemaComponents(raw CLIJson) {
	components := u.genLayout(raw)
	components = append(components, u.genHelpModal(raw)...)
	components = append(components, u.genCmdInnerTabs(raw)...)
	for _, c := range components {
		u.b.Component(c)
	}
}
