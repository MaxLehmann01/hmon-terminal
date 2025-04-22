package ui

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/maxlehmann01/hmon-terminal/pkg/config"
	"github.com/maxlehmann01/hmon-terminal/pkg/i2c"
	"github.com/maxlehmann01/hmon-terminal/pkg/plug"
	"github.com/stianeikeland/go-rpio"
)

const (
	SELECT_BUTTON_PIN = 4
	TOGGLE_BUTTON_PIN = 17
)

type GPIOUserInterface struct {
	display       *i2c.I2C
	displayConfig config.I2CDisplayConfig
	displayMutex  sync.Mutex
}

func (gpioui *GPIOUserInterface) StartControlListener(pm *plug.PlugManager) error {
	displayConfig := config.GetI2CDisplayConfig()

	i2cDevice, err := i2c.NewI2C(displayConfig.LCDAddress, displayConfig.LCDBus)
	if err != nil {
		return err
	}

	gpioui.display = i2cDevice
	gpioui.displayConfig = displayConfig

	if err := displayInit(i2cDevice, displayConfig); err != nil {
		return err
	}

	if err := rpio.Open(); err != nil {
		return errors.New("failed to open rpio")
	}

	selectButton := rpio.Pin(SELECT_BUTTON_PIN)
	selectButton.Input()
	selectButton.PullUp()

	toggleButton := rpio.Pin(TOGGLE_BUTTON_PIN)
	toggleButton.Input()
	toggleButton.PullUp()

	go buttonControlLoop(selectButton, toggleButton, pm)

	select {}
}

func (gpioui *GPIOUserInterface) OutputSelectedPlug(p *plug.Plug) error {
	if gpioui.display == nil {
		return errors.New("display device not initialized")
	}

	gpioui.displayMutex.Lock()
	defer gpioui.displayMutex.Unlock()

	line1, line2 := formatPlugOutput(p)

	if err := displaySendString(gpioui.display, padToWidth(line1, config.OUTPUT_WIDTH), gpioui.displayConfig.LCDLine1, gpioui.displayConfig); err != nil {
		return err
	}

	if err := displaySendString(gpioui.display, padToWidth(line2, config.OUTPUT_WIDTH), gpioui.displayConfig.LCDLine2, gpioui.displayConfig); err != nil {
		return err
	}

	return nil
}

func displayInit(dvc *i2c.I2C, cfg config.I2CDisplayConfig) error {
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
		if err := displaySendByte(dvc, cmd, 0, cfg); err != nil {
			return errors.New("failed to send init command")
		}
		time.Sleep(2 * time.Millisecond)
	}

	return nil
}

func displaySendByte(dvc *i2c.I2C, data byte, mode byte, cfg config.I2CDisplayConfig) error {
	if err := displayWriteNibble(dvc, data>>4, mode, cfg); err != nil {
		return err
	}

	if err := displayWriteNibble(dvc, data&0x0F, mode, cfg); err != nil {
		return err
	}

	return nil
}

func displayWriteNibble(dvc *i2c.I2C, nibble, mode byte, cfg config.I2CDisplayConfig) error {
	data := (nibble << 4) | mode | cfg.LCDBacklight

	if _, err := dvc.WriteBytes([]byte{data | cfg.Enable}); err != nil {
		return err
	}

	time.Sleep(cfg.EnableDelay)

	if _, err := dvc.WriteBytes([]byte{data & ^cfg.Enable}); err != nil {
		return err
	}

	time.Sleep(cfg.EnableDelay)

	return nil
}

func displaySendString(dvc *i2c.I2C, str string, line byte, cfg config.I2CDisplayConfig) error {
	if err := displaySendByte(dvc, line, 0, cfg); err != nil {
		return err
	}

	if len(str) > config.OUTPUT_WIDTH {
		str = str + fmt.Sprintf("%*s", config.OUTPUT_WIDTH-len(str), "")
	}

	for i := 0; i < len(str); i++ {
		if err := displaySendByte(dvc, str[i], cfg.RegisterSelect, cfg); err != nil {
			return fmt.Errorf("lcdString: error sending character '%c': %v", str[i], err)
		}
	}

	return nil
}

func buttonControlLoop(selectButton, toggleButton rpio.Pin, pm *plug.PlugManager) {
	for {
		if selectButton.Read() == rpio.Low {
			pm.SelectNext()
			OutputSelectedPlug(pm)

			time.Sleep(300 * time.Millisecond)

			for selectButton.Read() == rpio.Low {
				time.Sleep(10 * time.Millisecond)
			}
		}

		if toggleButton.Read() == rpio.Low {
			pm.ToggleSelected()
			OutputSelectedPlug(pm)

			time.Sleep(300 * time.Millisecond)

			for toggleButton.Read() == rpio.Low {
				time.Sleep(10 * time.Millisecond)
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func padToWidth(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + fmt.Sprintf("%-*s", width-len(s), "")
}
