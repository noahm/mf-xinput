package main

import (
	"fmt"
	"strings"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	driver "gitlab.com/gomidi/rtmididrv"

	"github.com/WesleiRamos/goxinput"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	controller := goxinput.NewController()
	fmt.Println("Created controller")

	if !controller.IsVBusExists() {
		panic("VBus driver is not installed")
	}
	fmt.Println("controller bus exists")

	// Plugin controller
	if error := controller.PlugIn(); error != nil {
		panic(error)
	}
	defer controller.Unplug()
	fmt.Println("Controller plugged in")

	drv, err := driver.New()
	must(err)
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	var midiFighter midi.In
	for i, input := range ins {
		if strings.Contains(input.String(), "Midi Fighter") {
			fmt.Println("Found midi fighter!")
			midiFighter = input
		} else if i+1 == len(ins) {
			printInPorts(ins)
			panic("Couldn't find midi fighter!")
		}
	}

	must(midiFighter.Open())
	defer midiFighter.Close()

	setKey := func(key uint8, isOn bool) {
		switch key {
		case 47: //p1 up
			controller.SetBtn(goxinput.BUTTON_LS, isOn)
			break
		case 50: //p1 down
			controller.SetBtn(goxinput.BUTTON_RS, isOn)
			break
		case 51: //p1 left
			if isOn {
				controller.SetTrigger(goxinput.BUTTON_LT, 1)
			} else {
				controller.SetTrigger(goxinput.BUTTON_LT, 0)
			}
			break
		case 46: //p1 right
			if isOn {
				controller.SetAxis(goxinput.AXIS_LX, 1)
			} else {
				controller.SetAxis(goxinput.AXIS_LX, 0)
			}
			break
		case 49: //p1 start
			controller.SetBtn(goxinput.BUTTON_START, isOn)
			break
		case 43: //p1 back
			controller.SetBtn(goxinput.BUTTON_BACK, isOn)
			break
		case 37: //p2 up
			controller.SetBtn(goxinput.BUTTON_Y, isOn)
			break
		case 40: //p2 down
			controller.SetBtn(goxinput.BUTTON_A, isOn)
			break
		case 41: //p2 left
			controller.SetBtn(goxinput.BUTTON_X, isOn)
			break
		case 36: //p2 right
			controller.SetBtn(goxinput.BUTTON_B, isOn)
			break
		case 44: //p2 start
			controller.SetBtn(goxinput.BUTTON_RB, isOn)
			break
		case 38: //p2 back
			controller.SetBtn(goxinput.BUTTON_LB, isOn)
			break
		default:
			fmt.Printf("Key %v is on: %t\n", key, isOn)
		}
	}

	rd := reader.New(
		reader.NoLogger(),
		reader.NoteOn(func(p *reader.Position, channel, key, vel uint8) {
			if channel == 2 {
				setKey(key, true)
			}
		}),
		reader.NoteOff(func(p *reader.Position, channel, key, vel uint8) {
			if channel == 2 {
				setKey(key, false)
			}
		}),
	)

	// listen for MIDI
	err = rd.ListenTo(midiFighter)
	must(err)

	// suspend this func without returning
	select {}
}

func printInPorts(ports []midi.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		fmt.Printf("[%v] %s\n", port.Number(), port.String())
	}
	fmt.Printf("\n\n")
}
