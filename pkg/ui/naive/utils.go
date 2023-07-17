package naive

import (
	"CLI2UI/pkg/ui"
	"fmt"
)

type Path struct {
	ui.Path
}

func (p Path) subcommandTabsId() string {
	return fmt.Sprintf("%sSubcommandTabs", p)
}
