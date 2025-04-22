package config

import (
	"flag"
	"time"
)

const OUTPUT_WIDTH = 16

type Flags struct {
	DevMode bool
}

type I2CDisplayConfig struct {
	LCDAddress     uint8
	LCDBus         int
	LCDLine1       uint8
	LCDLine2       uint8
	LCDBacklight   uint8
	Enable         uint8
	ReadWrite      uint8
	RegisterSelect uint8
	EnableDelay    time.Duration
}

func GetParsedFlags() *Flags {
	devModeFlag := flag.Bool("dev", false, "Run in development mode (output in console instead of i2c display)")
	flag.Parse()

	return &Flags{
		DevMode: *devModeFlag,
	}
}

func GetI2CDisplayConfig() I2CDisplayConfig {
	return I2CDisplayConfig{
		LCDAddress:     0x27,
		LCDBus:         1,
		LCDLine1:       0x80,
		LCDLine2:       0xC0,
		LCDBacklight:   0x08,
		Enable:         0x04,
		ReadWrite:      0x02,
		RegisterSelect: 0x01,
		EnableDelay:    1 * time.Millisecond,
	}
}
