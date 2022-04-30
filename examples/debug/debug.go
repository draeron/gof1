package main

import (
	"github.com/bearsh/hid"

	"github.com/draeron/gof1/examples/common"
	"github.com/draeron/gof1/pkg/device"
	button2 "github.com/draeron/gof1/pkg/f1/button"
	event2 "github.com/draeron/gof1/pkg/f1/event"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
)

func main() {
	log := logger.NewLogrus("main")
	common.Setup()

	if !hid.Supported() {
		log.Fatalln("hid not supported")
	}

	dev, err := device.Open()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer dev.Close()

	dev.EnableDebugLogger()

	// Set startup state
	dial := int8(0)
	colors := map[button2.Button]color.PaletteColor{}
	padcount := len(button2.Pads())
	for idx, col := range color.Colors() {
		col = col % color.White
		if idx < padcount {
			colors[button2.PadA1+button2.Button(idx)] = col
		}
	}

	for btn, col := range colors {
		err = dev.SetPadColor(btn, col)
		log.ErrorIf(err, "failed to ")
	}
	dev.SetDial(0)
	for _, btn := range button2.Mutes() {
		dev.SetBrightness(btn, 127)
	}
	for _, btn := range button2.Functions() {
		dev.SetBrightness(btn, 255)
	}

	// Start event listening
	events := make(chan event2.Event, 100)
	dev.Subscribe(events)

	dev.AddCallback(event2.IsButtonOfType(button2.Pads()...), func(evt event2.Event) {
		col, _ := colors[evt.Btn]
		col = (col + 1) % color.White
		if col > color.White {
			col = color.PaletteColor(0)
		}
		dev.SetPadColor(evt.Btn, col)
		colors[evt.Btn] = col
	})

	go func() {
		for evt := range events {
			log.Infof("%v", evt)

			switch {
			case evt.Btn.IsFunctions(), evt.Btn.IsMute():
				on := uint8(0)
				if evt.Value <= 0 {
					on = 255
				}
				dev.SetBrightness(evt.Btn, on)

			case evt.Btn == button2.Dial && evt.Type == event2.Pressed:
				dial = 0
				dev.SetDial(0)

			case evt.Type == event2.Decrement:
				if dial > -99 {
					dial--
				}
				dev.SetDial(dial)

			case evt.Type == event2.Increment:
				if dial < 99 {
					dial++
				}
				dev.SetDial(dial)
			}
		}
	}()

	common.WaitExit()
}
