package ui

// this file should be later integrated into https://github.com/webzard-io/sunmao-ui-go-binding/
// tags generated using: gomodifytags -file pkg/ui/props.go -all -add-tags json

type ButtonProperties struct {
	Text     string `json:"text"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Size     string `json:"size"`
	Shape    string `json:"shape"`
	Disabled bool   `json:"disabled"`
	Loading  bool   `json:"loading"`
	Long     bool   `json:"long"`
}

type StackProperties struct {
	Align     string `json:"align"`
	Direction string `json:"direction"`
	Justify   string `json:"justify"`
	Spacing   int    `json:"spacing"`
	Wrap      bool   `json:"wrap"`
}

type LayoutProperties struct {
	ShowHeader              bool `json:"showHeader"`
	ShowFooter              bool `json:"showFooter"`
	ShowSideBar             bool `json:"showSideBar"`
	SidebarCollapsible      bool `json:"sidebarCollapsible"`
	SidebarDefaultCollapsed bool `json:"sidebarDefaultCollapsed"`
}

type ModalProperties struct {
	Title          string `json:"title"`
	DefaultOpen    bool   `json:"defaultOpen"`
	Mask           bool   `json:"mask"`
	Simple         bool   `json:"simple"`
	OkText         string `json:"okText"`
	CancelText     string `json:"cancelText"`
	Closable       bool   `json:"closable"`
	MaskClosable   bool   `json:"maskClosable"`
	ConfirmLoading bool   `json:"confirmLoading"`
	UnmountOnExit  bool   `json:"unmountOnExit"`
}

type TextProperties struct {
	Raw    string `json:"raw"`
	Format string `json:"format"`
}

type ColumnProperties struct {
	Span   int `json:"span"`
	Offset int `json:"offset"`
}

type FormControlProperties struct {
	Label      TextProperties   `json:"label"`
	Required   bool             `json:"required"`
	Hidden     bool             `json:"hidden"`
	Layout     string           `json:"layout"`
	Extra      string           `json:"extra"`
	ErrorMsg   string           `json:"errorMsg"`
	Help       string           `json:"help"`
	LabelAlign string           `json:"labelAlign"`
	Colon      bool             `json:"colon"`
	LabelCol   ColumnProperties `json:"labelCol"`
	WrapperCol ColumnProperties `json:"wrapperCol"`
}
