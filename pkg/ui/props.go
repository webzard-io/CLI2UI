package ui

// this file should be later integrated into https://github.com/webzard-io/sunmao-ui-go-binding/
// tags generated using: gomodifytags -file pkg/ui/props.go -all -add-tags json

type TextDisplayProperties struct {
	Text   string `json:"text"`
	Format string `json:"format"`
}

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
