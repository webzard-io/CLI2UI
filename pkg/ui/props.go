package ui

// this file should be later integrated into https://github.com/webzard-io/sunmao-ui-go-binding/
// tags generated using: gomodifytags -file pkg/ui/props.go -all -add-tags json

type ButtonProperties[T string | bool] struct {
	Text     string `json:"text"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Size     string `json:"size"`
	Shape    string `json:"shape"`
	Disabled T      `json:"disabled"`
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

type NumberInputProperties[T bool | string] struct {
	DefaultValue                  int    `json:"defaultValue"`
	UpdateWhenDefaultValueChanges bool   `json:"updateWhenDefaultValueChanges"`
	Min                           int    `json:"min"`
	Max                           int    `json:"max"`
	Placeholder                   string `json:"placeholder"`
	Disabled                      T      `json:"disabled"`
	ButtonMode                    bool   `json:"buttonMode"`
	Precision                     int    `json:"precision"`
	Step                          int    `json:"step"`
	Size                          string `json:"size"`
	ReadOnly                      bool   `json:"readOnly"`
	Error                         bool   `json:"error"`
}

type SwitchProperties[T bool | string] struct {
	DefaultChecked                bool   `json:"defaultChecked"`
	UpdateWhenDefaultValueChanges bool   `json:"updateWhenDefaultValueChanges"`
	Disabled                      T      `json:"disabled"`
	Type                          string `json:"type"`
	Size                          string `json:"size"`
	Loading                       bool   `json:"loading"`
}

type SelectOptionProperties struct {
	Value    string `json:"value"`
	Text     string `json:"text"`
	Disabled bool   `json:"disabled"`
}

type SelectProperties[T string | bool] struct {
	DefaultValue                  string                   `json:"defaultValue"`
	Options                       []SelectOptionProperties `json:"options"`
	UpdateWhenDefaultValueChanges bool                     `json:"updateWhenDefaultValueChanges"`
	Multiple                      bool                     `json:"multiple"`
	LabelInValue                  bool                     `json:"labelInValue"`
	Placeholder                   string                   `json:"placeholder"`
	Bordered                      bool                     `json:"bordered"`
	Size                          string                   `json:"size"`
	Disabled                      T                        `json:"disabled"`
	Loading                       bool                     `json:"loading"`
	ShowSearch                    bool                     `json:"showSearch"`
	RetainInputValue              bool                     `json:"retainInputValue"`
	AllowClear                    bool                     `json:"allowClear"`
	AllowCreate                   bool                     `json:"allowCreate"`
	ShowTitle                     bool                     `json:"showTitle"`
	Error                         bool                     `json:"error"`
	UnmountOnExit                 bool                     `json:"unmountOnExit"`
	MountToBody                   bool                     `json:"mountToBody"`
	AutoFixPosition               bool                     `json:"autoFixPosition"`
	AutoAlignPopupMinWidth        bool                     `json:"autoAlignPopupMinWidth"`
	AutoAlignPopupWidth           bool                     `json:"autoAlignPopupWidth"`
	AutoFitPosition               bool                     `json:"autoFitPosition"`
	Position                      string                   `json:"position"`
}

type InputProperties[T string | bool] struct {
	DefaultValue                  string `json:"defaultValue"`
	Placeholder                   string `json:"placeholder"`
	UpdateWhenDefaultValueChanges bool   `json:"updateWhenDefaultValueChanges"`
	AllowClear                    bool   `json:"allowClear"`
	Disabled                      T      `json:"disabled"`
	ReadOnly                      bool   `json:"readOnly"`
	Error                         bool   `json:"error"`
	Size                          string `json:"size"`
}

type TabProperties struct {
	Title         string `json:"title"`
	Hidden        bool   `json:"hidden"`
	DestroyOnHide bool   `json:"destroyOnHide"`
}

type TabsProperties struct {
	DefaultActiveTab              int             `json:"defaultActiveTab"`
	Tabs                          []TabProperties `json:"tabs"`
	UpdateWhenDefaultValueChanges bool            `json:"updateWhenDefaultValueChanges"`
	Type                          string          `json:"type"`
	TabPosition                   string          `json:"tabPosition"`
	Size                          string          `json:"size"`
}

type CheckboxOptionProperties struct {
	Label        string `json:"label"`
	Value        string `json:"value"`
	Disabled     bool   `json:"disabled"`
	Intermediate bool   `json:"intermediate"`
}

type CheckboxProperties[T string | bool] struct {
	Options                       []CheckboxOptionProperties `json:"options"`
	DefaultCheckedValues          []string                   `json:"defaultCheckedValues"`
	Direction                     string                     `json:"direction"`
	ShowCheckAll                  bool                       `json:"showCheckAll"`
	CheckAllText                  string                     `json:"checkAllText"`
	Disabled                      T                          `json:"disabled"`
	UpdateWhenDefaultValueChanges bool                       `json:"updateWhenDefaultValueChanges"`
}

type CollapseItemProperties struct {
	Key            string `json:"key"`
	Header         string `json:"header"`
	Disabled       bool   `json:"disabled"`
	ShowExpandIcon bool   `json:"showExpandIcon"`
	DestroyOnHide  bool   `json:"destroyOnHide"`
}

type CollapseProperties struct {
	DefaultActiveKey              []string                 `json:"defaultActiveKey"`
	Accordion                     bool                     `json:"accordion"`
	ExpandIconPosition            string                   `json:"expandIconPosition"`
	Bordered                      bool                     `json:"bordered"`
	Options                       []CollapseItemProperties `json:"options"`
	UpdateWhenDefaultValueChanges bool                     `json:"updateWhenDefaultValueChanges"`
	LazyLoad                      bool                     `json:"lazyLoad"`
	DestroyOnHide                 bool                     `json:"destroyOnHide"`
}
