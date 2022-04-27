package f1

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"

	"github.com/bearsh/hid"

	seven_bits "github.com/draeron/gof1/pkg/color/7bits"
	"github.com/draeron/gof1/pkg/f1/button"
	"github.com/draeron/gopkgs/color"
)

type OutState struct {
	SevenSegment int8 // [-99,99] sign means dot is turned on
	Functions    map[button.Button]LEDIntensity
	Pads         [16]color.Color
	Mute         [4]LEDIntensity
}

type LEDIntensity uint8

func (l LEDIntensity) Value() uint8 {
	if l > 127 {
		return 127
	}
	return uint8(l)
}

func NewOutState() *OutState {
	o := OutState{
		Functions: map[button.Button]LEDIntensity{},
	}
	for idx, _ := range o.Pads {
		o.Pads[idx] = color.Black
	}
	for _, btn := range button.Functions() {
		o.Functions[btn] = 0
	}
	return &o
}

func (o OutState) Write(device *hid.Device) error {
	var err error

	writer := bytes.NewBuffer([]byte{})

	/*
		The LEDs of the F1 are set using a single output report with a length of 81 Bytes.
		All animations including blinking states, and the white line animation displayed when changing
		banks are controlled by the software. That is a blinking led is turned on by the software, then
		toggled by the software at a certain interval, rather than being put into a blinking state where
		the strobing is controlled purely by the F1 hardware.

		All leds appear to have some level of brightness control, Traktor allows you to vary the
		“On State Brightness” and “Dim State Percentage” levels using the configuration menu. All the values
		listed in this document are for an On state brightness of 100%, and a Dim state Percentage of 0%.

		Byte 01              ID (= 80)
		Byte 02 .. 17      7-segment displays
		Byte 18 .. 25      Small Function Keys
		Byte 26 .. 73      RGB Pads
		Byte 74 .. 81      Stop Keys

	*/
	// 		The first byte is always 80.
	err = binary.Write(writer, binary.LittleEndian, byte(0x80))
	if err != nil {
		return errors.WithMessage(err, "failed to write HID packet")
	}

	/*
		Bytes 02 thru 17                7 Segment Displays

		The next 16 bytes control the 7 segment displays, each byte represents one of the segments. A value of 64
		is used when that segment is on, and a zero value when it is off. Bytes 2 - 9 control the right hand digit,
		and bytes 10 - 17 control the left digit.

		Note that the DP actually appears top left of the digit.

		Byte 1: DP
		Byte 2: Segment G
		Byte 3: Segment C
		Byte 4: Segment B
		Byte 5: Segment A
		Byte 6: Segment F
		Byte 7: Segment E
		Byte 8: Segment D
	*/
	err = binary.Write(writer, binary.LittleEndian, make([]byte, 16))
	if err != nil {
		return errors.WithMessage(err, "failed to write HID packet")
	}

	/*
		Bytes 18 thru 25        Small Function Keys

		The next 8 bytes control the brightness of the 8 small functions keys located near the middle of the device.
		The full brightness value is 7F, off is 0.
		Byte 1     Browse
		Byte 2     Size
		Byte 3     Type
		Byte 4     Reverse
		Byte 5     Shift
		Byte 6     Capture
		Byte 7     Quant
		Byte 8     Sync
	*/
	err = binary.Write(writer, binary.LittleEndian, []byte{
		o.Functions[button.Browse].Value(),
		o.Functions[button.Size].Value(),
		o.Functions[button.Type].Value(),
		o.Functions[button.Reverse].Value(),
		o.Functions[button.Shift].Value(),
		o.Functions[button.Capture].Value(),
		o.Functions[button.Quant].Value(),
		o.Functions[button.Sync].Value(),
	})
	if err != nil {
		return errors.WithMessage(err, "failed to write HID packet")
	}

	/*
		Bytes 26 thru 73     RGB Pads

		The next 48 bytes are used to set the color of each of the 16 pads. 3 bytes are used for the RGB color
		settings of each pad are arranged in BRG order: Blue, Red, Green.
		Byte 1     Pad #1 Blue
		Byte 2     Pad #1 Red
		Byte 3     Pad #1 Green
		Byte 4     Pad #2 Blue
		Byte 5     Pad #2 Red
		Byte 6     Pad #2 Green

		White is sent as 73 73 73. Each color uses a 7 bit resolution, min (off) = 0, max (100% on) = 0x7D

		I will add further information on the color settings later.

		There is no buffering of the color states of other banks, the color information is fully refreshed
		from the host in response a bank change key press.
	*/
	for _, pad := range o.Pads {
		rgb := seven_bits.FromColor(pad.RGB())
		err = binary.Write(writer, binary.LittleEndian, []byte{
			rgb.B,
			rgb.R,
			rgb.G,
		})
		if err != nil {
			return errors.WithMessage(err, "failed to write HID packet")
		}
	}

	/*
		Bytes 74 thru 81     Stop Keys

		The next 8 bytes are used to control the brightness of the 4 stop keys at the bottom of the device.
		Each of these keys uses two leds to provide sufficient illumination along its length, it appears you can control
		these LEDs separately as each key uses two control bytes.
		Byte 74     Column 4 Stop Key LED 1
		Byte 75     Column 4 Stop Key LED 2
		Byte 76     Column 3 Stop Key LED 1
		Byte 77     Column 3 Stop Key LED 2
		Byte 78     Column 2 Stop Key LED 1
		Byte 79     Column 2 Stop Key LED 2
		Byte 80     Column 1 Stop Key LED 1
		Byte 81     Column 1 Stop Key LED 2
	*/
	for idx := len(button.Mutes()) - 1; idx >= 0; idx-- {
		mute := o.Mute[idx]
		err = binary.Write(writer, binary.LittleEndian, []byte{
			mute.Value(),
			mute.Value(),
		})
		if err != nil {
			return errors.WithMessage(err, "failed to write HID packet")
		}
	}

	packet := writer.Bytes()
	// if len(packet) != 81 {
	// 	return errors.New("wrong computed packet length")
	// }

	wrote, err := device.Write(packet)
	log.Debugf("wrote %d bytes to HID devices", wrote)
	return errors.WithMessage(err, "failed to write HID packet")
}
