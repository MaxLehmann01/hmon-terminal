package main

import (
	"log/slog"

	"github.com/maxlehmann01/hmon-terminal/pkg/config"
	"github.com/maxlehmann01/hmon-terminal/pkg/plug"
	"github.com/maxlehmann01/hmon-terminal/pkg/server"
	"github.com/maxlehmann01/hmon-terminal/pkg/ui"
)

func main() {
	flags := config.GetParsedFlags()

	slog.Info("Starting hmon-terminal", "devMode", flags.DevMode)

	plugManager := plug.NewPlugManager(flags.BackendUrl)

	if flags.DevMode {
		ui.SetUserInterface(&ui.ConsoleUserInterface{})
	} else {
		ui.SetUserInterface(&ui.GPIOUserInterface{})
	}

	server.Start(plugManager, flags.Port)

	ui.StartControlListener(plugManager)
}
