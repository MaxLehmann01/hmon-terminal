package main

import (
	"log/slog"

	"github.com/maxlehmann01/hmon-terminal/pkg/config"
)

func main() {
	flags := config.GetParsedFlags()

	slog.Info("Starting hmon-terminal", "devMode", flags.DevMode)
}
