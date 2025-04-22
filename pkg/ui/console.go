package ui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eiannone/keyboard"
	"github.com/maxlehmann01/hmon-terminal/pkg/plug"
)

type ConsoleUserInterface struct{}

func (cui *ConsoleUserInterface) StartControlListener(pm *plug.PlugManager) error {
	if err := keyboard.Open(); err != nil {
		return err
	}
	defer keyboard.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		keyboard.Close()
		os.Exit(0)
	}()

	OutputSelectedPlug(pm)

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		switch key {
		case keyboard.KeyArrowRight:
			pm.SelectNext()
			OutputSelectedPlug(pm)

		case keyboard.KeyEnter:
			pm.ToggleSelected()
			OutputSelectedPlug(pm)

		case keyboard.KeyEsc:
			return nil
		}
	}
}

func (cui *ConsoleUserInterface) OutputSelectedPlug(p *plug.Plug) error {
	line1, line2 := formatPlugOutput(p)

	fmt.Print("\033[H\033[2J")
	fmt.Println(line1)
	fmt.Println(line2)

	return nil
}
