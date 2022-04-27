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

	// ctrl, err := f1.Open()
	// common.Must(err)
	// log.Debug(ctrl.String())

	if !hid.Supported() {
		log.Fatalln("hid not supported")
	}

	dev, err := f1.Open()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer dev.Close()

	events := make(chan event.Event, 100)

	dev.Subscribe(events)

	go func() {
		for evt := range events {
			log.Infof("%v", evt)

			switch {
			case evt.Btn.IsMute():
				dev.SetBrightness(evt.Btn, 0)
			}
		}
	}()

	padcount := len(button.Pads())
	for idx, col := range color.Colors() {
		if idx < padcount {
			err = dev.SetPadColor(button.PadA1+button.Button(idx), col)
			log.ErrorIf(err, "failed to ")
		}
	}

	// for _, btn := range button.Pads() {
	// 	dev.SetPadColor(btn, color.Blue)
	// }

	for _, btn := range button.Mutes() {
		dev.SetBrightness(btn, 127)
	}

	for _, btn := range button.Functions() {
		dev.SetBrightness(btn, 255)
	}

	common.WaitExit()
}
