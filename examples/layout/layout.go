package main

import (
	"github.com/draeron/gof1/examples/common"
	"github.com/draeron/gof1/pkg/device"
	"github.com/draeron/gof1/pkg/f1"
	"github.com/draeron/gof1/pkg/layout"
	"github.com/draeron/gopkgs/logger"
)

func main() {
	log := logger.NewLogrus("main")
	common.Setup()

	f1dev, err := device.Open()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	layout := layout.NewLayoutPreset(f1.MaskMutes)
	layout.Connect(f1dev)

	layout.SetHandler()

	common.WaitExit()
}
