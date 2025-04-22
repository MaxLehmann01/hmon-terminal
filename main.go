package main

import (
	"log/slog"

	"github.com/maxlehmann01/hmon-terminal/pkg/config"
	"github.com/maxlehmann01/hmon-terminal/pkg/plug"
	"github.com/maxlehmann01/hmon-terminal/pkg/ui"
)

func main() {
	flags := config.GetParsedFlags()

	slog.Info("Starting hmon-terminal", "devMode", flags.DevMode)

	plugManager := plug.NewPlugManager()

	plugManager.AddPlug(&plug.Plug{ID: 1, Name: "Plug 1", PowerUsage: 100})
	plugManager.AddPlug(&plug.Plug{ID: 2, Name: "Plug 2", PowerUsage: 201.2})
	plugManager.AddPlug(&plug.Plug{ID: 3, Name: "Plug 3", PowerUsage: 35.1})

	if flags.DevMode {
		ui.SetUserInterface(&ui.ConsoleUserInterface{})
	}

	ui.StartControlListener(plugManager)
}
