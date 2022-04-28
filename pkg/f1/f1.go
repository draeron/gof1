package f1

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"

	"github.com/bearsh/hid"

	"github.com/draeron/gof1/pkg/f1/button"
	f1event "github.com/draeron/gof1/pkg/f1/event"
)

type Controller struct {
	device      *hid.Device
	subscribers []chan<- f1event.Event
	lastOut     *OutState
	mutex       sync.RWMutex
}

func Open() (*Controller, error) {

	if !hid.Supported() {
		return nil, errors.New("HID USB operations not supported on this platform")
	}

	var err error
	ctrl := &Controller{
		lastOut: NewOutState(),
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
		break
	}

	if selected == nil {
		return nil, errors.New("no F1 controller were found")
	}

	ctrl.device, err = selected.Open()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open Traktor F1 HID device")
	}

	log.Infof("opened device: %v", ctrl.device)

	err = ctrl.lastOut.Write(ctrl.device)
	if err != nil {
		log.Errorf("failed to init HID state: %+v", err)
	}

	go ctrl.processInput()

	return ctrl, nil
}

func (c *Controller) processInput() {
	log.Infof("starting to read from HID device")
	defer log.Infof("stopped reading from HID device")

	previous := NewInState()

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

			// ignore value on first dial event
			if first {
				log.Infof("first received message is ignored")
				previous.Dial = current.Dial
				first = false
			} else {
			}

			c.compareState(previous, current)

			previous = current

			// log.Infof("%v bytes were read from HID device", length)
		}
	}
}

/*
	Compare state to previous and emit events based on difference
*/
func (c *Controller) compareState(previous, current *InState) {
	evt := f1event.Event{}

	if current.Dial != previous.Dial {
		evt.Btn = button.Dial

		pd := previous.Dial
		cd := current.Dial
		if (pd == 255 && cd == 0) || (pd == 0 && cd == 255) {
			pd = current.Dial
			cd = previous.Dial
		}

		if cd > pd {
			evt.Type = f1event.Increment
			evt.Value = int16(current.Dial)
		} else {
			evt.Type = f1event.Decrement
			evt.Value = int16(current.Dial)
		}

		c.sendToSubscribers(evt)
	}

	for key, state := range current.PressedButtons {
		if current.PressedButtons[key] != previous.PressedButtons[key] {
			evt.Btn, _ = button.ParseButton(key)
			if state {
				evt.Type = f1event.Pressed
				evt.Value = 1
			} else {
				evt.Type = f1event.Released
				evt.Value = 0
			}
			c.sendToSubscribers(evt)
		}
	}

	const maxval = 4090

	for idx, value := range current.Volumes {
		if current.Volumes[idx] != previous.Volumes[idx] {
			evt.Btn = button.Volume1 + button.Button(idx)
			evt.Type = f1event.Changed
			evt.Value = int16(float64(value) / maxval * 256)
			c.sendToSubscribers(evt)
		}
	}

	for idx, value := range current.Filters {
		if current.Filters[idx] != previous.Filters[idx] {
			evt.Btn = button.Filter1 + button.Button(idx)
			evt.Type = f1event.Changed
			evt.Value = int16(float32(value) / maxval * 256)
			c.sendToSubscribers(evt)
		}
	}
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

func (c *Controller) Subscribe(channel chan<- f1event.Event) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.subscribers = append(c.subscribers, channel)
}

func (c *Controller) sendToSubscribers(evt f1event.Event) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for _, channel := range c.subscribers {
		channel <- evt
	}
}
