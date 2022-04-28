package f1

import (
	"github.com/pkg/errors"

	"github.com/draeron/gof1/pkg/f1/button"
	"github.com/draeron/gopkgs/color"
)

func (c *Controller) SetPadColor(btn button.Button, color color.Color) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !btn.IsPad() {
		return errors.Errorf("button %v is not a pad button", btn)
	}

	idx := btn - button.PadA1
	c.state.out.Pads[idx] = color

	err := c.state.out.Write(c.device)
	return errors.WithMessage(err, "failed to write to HID device")
}

func (c *Controller) SetDial(val int8) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.state.out.SevenSegment = val

	err := c.state.out.Write(c.device)
	return errors.WithMessage(err, "failed to write to HID device")
}

func (c *Controller) SetBrightness(btn button.Button, val uint8) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if btn.IsPad() || btn.IsKnob() || btn.IsFader() {
		return errors.Errorf("button %v is cannot set on")
	}

	bright := LEDIntensity(val)

	switch {
	case btn.IsMute():
		idx := btn - button.Mute1
		c.state.out.Mute[idx] = bright
	case btn.IsFunctions():
		c.state.out.Functions[btn] = bright
	}

	err := c.state.out.Write(c.device)
	return errors.WithMessage(err, "failed to write to HID device")
}
