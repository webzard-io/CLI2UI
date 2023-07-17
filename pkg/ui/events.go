package ui

import (
	"CLI2UI/pkg/config"
	"fmt"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

type UpdateOptionValueParams struct {
	Path       Path
	OptionName string
	Value      any
}

func UpdateValueEvent(field string, p Path, o config.Option) []sunmao.EventHandler {
	return []sunmao.EventHandler{
		{
			Type:        "onChange",
			ComponentId: "$utils",
			Method: sunmao.EventMethod{
				Name: "binding/v1/UpdateOptionValue",
				Parameters: UpdateOptionValueParams{
					OptionName: o.Name,
					Path:       p,
					Value:      fmt.Sprintf("{{ %s.%s }}", p.OptionValueInputId(o.Name), field),
				},
			},
		},
	}
}

func UpdateCheckedOptions(v *map[string]*config.OptionValue, checked []string) {
	for k, v := range *v {
		found := false
		for _, cv := range checked {
			if k == cv {
				v.Enabled = true
				found = true
				break
			}
		}
		if !found {
			v.Enabled = false
			v.ResetValue()
		}
	}
}
