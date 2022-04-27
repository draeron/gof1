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
	c.lastOut.Pads[idx] = color

	err := c.lastOut.Write(c.device)
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
		c.lastOut.Mute[idx] = bright
	case btn.IsFunctions():
		c.lastOut.Functions[btn] = bright
	}

	err := c.lastOut.Write(c.device)
	return errors.WithMessage(err, "failed to write to HID device")
}
