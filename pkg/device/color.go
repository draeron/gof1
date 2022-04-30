package device

import (
	"github.com/pkg/errors"

	"github.com/draeron/gof1/pkg/f1/button"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/color/7bits"
)

func (d *Device) SetDial(val int8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.state.out.SevenSegment = val

	err := d.state.out.Write(d.device)
	return errors.WithMessage(err, "failed to write to HID device")
}

func (d *Device) SetBrightness(btn button.Button, val uint8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if btn.IsPad() || btn.IsKnob() || btn.IsFader() {
		return errors.Errorf("button %v is cannot set on")
	}

	bright := LEDIntensity(val)

	switch {
	case btn.IsMute():
		idx := btn - button.Mute1
		d.state.out.Mute[idx] = bright
	case btn.IsFunctions():
		d.state.out.Functions[btn] = bright
	}

	err := d.state.out.Write(d.device)
	return errors.WithMessage(err, "failed to write to HID device")
}

func (d *Device) SetPadColorAll(col color.Color) error {
	return d.SetPadColorMany(button.Values(), col)
}

func (d *Device) SetPadColorMany(btns []button.Button, col color.Color) error {
	mapp := button.ColorMap{}
	t := seven_bits.FromColor(col)
	for _, b := range btns {
		mapp[b] = t
	}
	return d.SetPadColors(mapp)
}

func (d *Device) SetPadColor(btn button.Button, col color.Color) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if btn.IsPad() {
		d.state.out.Pads[btn-button.PadA1] = col
	} else {
		return errors.Errorf("button %v is not a pad", btn)
	}

	err := d.state.out.Write(d.device)
	return errors.WithMessage(err, "failed to write to HID device")
}

func (d *Device) SetPadColors(mapp button.ColorMap) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for btn, col := range mapp {
		if btn.IsPad() {
			d.state.out.Pads[btn-button.PadA1] = col
		} else {
			return errors.Errorf("button %v is not a pad", btn)
		}
	}

	err := d.state.out.Write(d.device)
	return errors.WithMessage(err, "failed to write to HID device")
}

// func (c *Device) SetPadColor(btn button.Button, color color.Color) error {
// 	c.mutex.Lock()
// 	defer c.mutex.Unlock()
//
// 	if !btn.IsPad() {
// 		return errors.Errorf("button %v is not a pad button", btn)
// 	}
//
// 	idx := btn - button.PadA1
// 	c.state.out.Pads[idx] = color
//
// 	err := c.state.out.Write(c.device)
// 	return errors.WithMessage(err, "failed to write to HID device")
// }
