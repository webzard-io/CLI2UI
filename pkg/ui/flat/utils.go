package flat

import (
	"CLI2UI/pkg/ui"
	"fmt"
)

type Path struct {
	ui.Path
}

func (p Path) menuItemKey() string {
	return fmt.Sprintf("MenuItem%s", p)
}
