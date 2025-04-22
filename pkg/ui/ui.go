package ui

import (
	"fmt"

	"github.com/maxlehmann01/hmon-terminal/pkg/config"
	"github.com/maxlehmann01/hmon-terminal/pkg/plug"
)

type UserInterface interface {
	StartControlListener(pm *plug.PlugManager) error
	OutputSelectedPlug(p *plug.Plug)
}

var UI UserInterface

func SetUserInterface(ui UserInterface) {
	UI = ui
}

func StartControlListener(pm *plug.PlugManager) {
	if UI != nil {
		UI.StartControlListener(pm)
	}
}

func OutputSelectedPlug(pm *plug.PlugManager) {
	if UI != nil {
		selectedPlug := pm.GetSelected()
		UI.OutputSelectedPlug(selectedPlug)
	}
}

func formatPlugOutput(p *plug.Plug) (line1 string, line2 string) {
	width := config.OUTPUT_WIDTH

	state := "[OFF]"
	if p.IsOn {
		state = "[ON]"
	}

	nameMaxLen := width - len(state)
	name := p.Name

	if len(name) > nameMaxLen {
		name = name[:nameMaxLen]
	}

	line1 = fmt.Sprintf("%-*s%s", nameMaxLen, name, state)

	powerStr := fmt.Sprintf("%.1fW", p.PowerUsage)
	line2 = fmt.Sprintf("%*s", width, powerStr)

	return line1, line2
}
