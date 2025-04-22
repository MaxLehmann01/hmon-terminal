package config

import "flag"

type Flags struct {
	DevMode bool
}

func GetParsedFlags() *Flags {
	devModeFlag := flag.Bool("dev", false, "Run in development mode (output in console instead of i2c display)")
	flag.Parse()

	return &Flags{
		DevMode: *devModeFlag,
	}
}
