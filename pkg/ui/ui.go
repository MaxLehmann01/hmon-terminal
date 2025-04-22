package ui

import (
	"errors"
	"fmt"

	"github.com/maxlehmann01/hmon-terminal/pkg/config"
	"github.com/maxlehmann01/hmon-terminal/pkg/plug"
)

type UserInterface interface {
	StartControlListener(pm *plug.PlugManager) error
	OutputSelectedPlug(p *plug.Plug) error
}

var UI UserInterface

func SetUserInterface(ui UserInterface) {
	UI = ui
}

func StartControlListener(pm *plug.PlugManager) error {
	if UI == nil {
		return errors.New("no user interface set")
	}

	UI.StartControlListener(pm)

	return nil
}

func OutputSelectedPlug(pm *plug.PlugManager) error {
	if UI == nil {
		return errors.New("no user interface set")
	}

	selectedPlug := pm.GetSelected()
	UI.OutputSelectedPlug(selectedPlug)

	return nil
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
