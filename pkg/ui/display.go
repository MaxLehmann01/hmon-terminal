package ui

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/maxlehmann01/hmon-terminal/pkg/config"
	"github.com/maxlehmann01/hmon-terminal/pkg/i2c"
	"github.com/maxlehmann01/hmon-terminal/pkg/plug"
	"github.com/stianeikeland/go-rpio"
)

type DisplayUserInterface struct {
	device        *i2c.I2C
	displayConfig config.I2CDisplayConfig
	displayMutex  sync.Mutex
}

func (dui *DisplayUserInterface) StartControlListener(pm *plug.PlugManager) error {
	displayConfig := config.GetI2CDisplayConfig()

	i2cDevice, err := i2c.NewI2C(displayConfig.LCDAddress, displayConfig.LCDBus)
	if err != nil {
		return err
	}

	dui.device = i2cDevice
	dui.displayConfig = displayConfig

	if err := displayInit(i2cDevice, displayConfig); err != nil {
		return err
	}

	if err := rpio.Open(); err != nil {
		return errors.New("failed to open rpio")
	}

	buttonSelect := rpio.Pin(config.BUTTON_SELECT_PIN)
	buttonToggle := rpio.Pin(config.BUTTON_TOGGLE_PIN)

	buttonSelect.Input()
	buttonSelect.PullUp()

	buttonToggle.Input()
	buttonToggle.PullUp()

	go func() {
		for {
			if buttonSelect.Read() == rpio.Low {
				pm.SelectNext()
				OutputSelectedPlug(pm)

				time.Sleep(300 * time.Millisecond)

				for buttonSelect.Read() == rpio.Low {
					time.Sleep(10 * time.Millisecond)
				}
			}

			if buttonToggle.Read() == rpio.Low {
				pm.ToggleSelected()
				OutputSelectedPlug(pm)

				time.Sleep(300 * time.Millisecond)

				for buttonToggle.Read() == rpio.Low {
					time.Sleep(10 * time.Millisecond)
				}
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()

	select {}
}

func (dui *DisplayUserInterface) OutputSelectedPlug(p *plug.Plug) error {
	if dui.device == nil {
		return errors.New("display device not initialized")
	}

	dui.displayMutex.Lock()
	defer dui.displayMutex.Unlock()

	line1, line2 := formatPlugOutput(p)

	if err := displaySendString(dui.device, padToWidth(line1, config.OUTPUT_WIDTH), dui.displayConfig.LCDLine1, dui.displayConfig); err != nil {
		log.Printf("Failed to write line 1: %v", err)
	}

	if err := displaySendString(dui.device, padToWidth(line2, config.OUTPUT_WIDTH), dui.displayConfig.LCDLine2, dui.displayConfig); err != nil {
		log.Printf("Failed to write line 2: %v", err)
	}

	return nil
}

func displayInit(device *i2c.I2C, cfg config.I2CDisplayConfig) error {
	time.Sleep(50 * time.Millisecond)

	initSequence := []byte{
		0x33,
		0x32,
		0x28,
		0x0C,
		0x06,
		0x01,
	}

	for _, cmd := range initSequence {
		if err := displaySendByte(device, cmd, 0, cfg); err != nil {
			return errors.New("failed to send init command")
		}
		time.Sleep(2 * time.Millisecond)
	}

	return nil
}

func displaySendByte(device *i2c.I2C, data byte, mode byte, cfg config.I2CDisplayConfig) error {
	if err := displayWriteNibble(device, data>>4, mode, cfg); err != nil {
		return err
	}

	if err := displayWriteNibble(device, data&0x0F, mode, cfg); err != nil {
		return err
	}

	return nil
}

func displayWriteNibble(device *i2c.I2C, nibble, mode byte, cfg config.I2CDisplayConfig) error {
	data := (nibble << 4) | mode | cfg.LCDBacklight

	if _, err := device.WriteBytes([]byte{data | cfg.Enable}); err != nil {
		return err
	}

	time.Sleep(cfg.EnableDelay)

	if _, err := device.WriteBytes([]byte{data & ^cfg.Enable}); err != nil {
		return err
	}

	time.Sleep(cfg.EnableDelay)

	return nil
}

func displaySendString(device *i2c.I2C, str string, line byte, cfg config.I2CDisplayConfig) error {
	if err := displaySendByte(device, line, 0, cfg); err != nil {
		return err
	}

	if len(str) > config.OUTPUT_WIDTH {
		str = str + fmt.Sprintf("%*s", config.OUTPUT_WIDTH-len(str), "")
	}

	for i := 0; i < len(str); i++ {
		if err := displaySendByte(device, str[i], cfg.RegisterSelect, cfg); err != nil {
			return fmt.Errorf("lcdString: error sending character '%c': %v", str[i], err)
		}
	}

	return nil
}

func padToWidth(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + fmt.Sprintf("%-*s", width-len(s), "")
}
