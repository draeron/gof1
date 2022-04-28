package f1

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"

	"github.com/bearsh/hid"

	"github.com/draeron/gof1/pkg/f1/event"
)

type Controller struct {
	device      *hid.Device
	subscribers []chan<- event.Event
	state       State
	mutex       sync.RWMutex
}

func Open() (*Controller, error) {

	if !hid.Supported() {
		return nil, errors.New("HID USB operations not supported on this platform")
	}

	var err error
	ctrl := &Controller{
		state: State{
			out: NewOutState(),
			in:  InState{},
		},
	}

	var selected *hid.DeviceInfo
	for _, devinfo := range hid.Enumerate(6092, 4384) {
		jinfo, _ := json.MarshalIndent(devinfo, "", "  ")
		log.Infof("info: \n%v", string(jinfo))

		const f1product = "Traktor Kontrol F1"
		if devinfo.Product != f1product {
			log.Warnf("usb product name '%s' is not equal to '%s'", f1product)
		}

		selected = &devinfo

	}

	if selected == nil {
		return nil, errors.New("no F1 controller were found")
	}

	ctrl.device, err = selected.Open()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open Traktor F1 HID device")
	}

	log.Infof("opened device: %v", ctrl.device)

	err = ctrl.state.out.Write(ctrl.device)
	if err != nil {
		log.Errorf("failed to init HID state: %+v", err)
	}

	go ctrl.processInput()

	return ctrl, nil
}

func (c *Controller) processInput() {
	log.Infof("starting to read from HID device")
	defer log.Infof("stopped reading from HID device")

	first := true

	buffer := make([]byte, 22)
	for {
		length, err := c.device.Read(buffer)
		if err != nil {
			log.Errorf("failed to read buffer from HID device")
			return
		} else if length > 0 {
			current := NewInState()

			err = current.UnpackPacket(bytes.NewReader(buffer))
			if err != nil {
				log.Errorf("failed to parse HID packet")
			}

			previous := c.State()

			// ignore value on first dial event
			if first {
				log.Infof("first received message is ignored")
				previous.in.Dial = current.Dial
				first = false
			}

			events := previous.eventFromDiff(*current)

			for _, evt := range events {
				c.sendToSubscribers(evt)
			}

			// replace input state
			c.mutex.Lock()
			c.state.in = *current
			c.mutex.Unlock()

			// log.Infof("%v bytes were read from HID device", length)
		}
	}
}

/*
	thread safe state retrieval
*/
func (c *Controller) State() State {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.state.Copy()
}

func (c *Controller) Close() {
	if c.device != nil {
		c.device.Close()
	}
}

func (c *Controller) String() string {
	return "papapa"
	// return fmt.Sprintf("F1 Controller, IN: %v, Out: %v", c.dev.port.In, c.dev.port.Out)
}
