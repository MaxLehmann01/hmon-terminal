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
	if p == nil {
		return "No plug selected", ""
	}

	maxWidth := config.OUTPUT_WIDTH

	state := "[OFF]"
	if p.IsOn {
		state = "[ON]"
	}

	nameMaxLength := maxWidth - len(state)
	name := p.Name
	if len(name) > nameMaxLength {
		name = name[:nameMaxLength]
	}

	power := fmt.Sprintf("%.1fW", p.PowerUsage)

	protected := "[U]"
	if p.IsProtected {
		protected = "[P]"
	}

	protectedMaxLength := maxWidth - len(power)

	line1 = fmt.Sprintf("%-*s%s", nameMaxLength, name, state)
	line2 = fmt.Sprintf("%-*s%s", protectedMaxLength, protected, power)
	return line1, line2
}
