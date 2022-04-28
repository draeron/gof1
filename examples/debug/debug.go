package main

import (
	"github.com/bearsh/hid"

	"github.com/draeron/gof1/examples/common"
	"github.com/draeron/gof1/pkg/f1"
	"github.com/draeron/gof1/pkg/f1/button"
	"github.com/draeron/gof1/pkg/f1/event"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
)

func main() {
	log := logger.NewLogrus("main")
	common.Setup()

	if !hid.Supported() {
		log.Fatalln("hid not supported")
	}

	dev, err := f1.Open()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer dev.Close()

	// Set startup state
	dial := int8(0)
	colors := map[button.Button]color.PaletteColor{}
	padcount := len(button.Pads())
	for idx, col := range color.Colors() {
		col = col % color.White
		if idx < padcount {
			colors[button.PadA1+button.Button(idx)] = col
		}
	}

	for btn, col := range colors {
		err = dev.SetPadColor(btn, col)
		log.ErrorIf(err, "failed to ")
	}
	dev.SetDial(0)
	for _, btn := range button.Mutes() {
		dev.SetBrightness(btn, 127)
	}
	for _, btn := range button.Functions() {
		dev.SetBrightness(btn, 255)
	}

	// Start event listening
	events := make(chan event.Event, 100)
	dev.Subscribe(events)

	dev.AddCallback(event.IsButtonOfType(button.Pads()...), func(evt event.Event) {
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

			case evt.Btn == button.Dial && evt.Type == event.Pressed:
				dial = 0
				dev.SetDial(0)

			case evt.Type == event.Decrement:
				if dial > -99 {
					dial--
				}
				dev.SetDial(dial)

			case evt.Type == event.Increment:
				if dial < 99 {
					dial++
				}
				dev.SetDial(dial)
			}
		}
	}()

	common.WaitExit()
}
