package device

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"

	button2 "github.com/draeron/gof1/pkg/f1/button"
)

type InState struct {
	Version        byte
	PressedButtons map[button2.Button]button2.PushState
	Dial           uint8
	Filters        [4]uint16
	Volumes        [4]uint16
}

func NewInState() *InState {
	in := &InState{
		PressedButtons: map[button2.Button]button2.PushState{},
	}
	for _, btn := range button2.Push.Buttons() {
		in.PressedButtons[btn] = button2.Released
	}
	return in
}

/*
	The state of all input controls is communicated via a single input report of 22 Bytes
	The first byte is the version number, currently 0x01
*/
func (packet *InState) UnpackPacket(rdr io.Reader) error {
	var err error

	err = binary.Read(rdr, binary.LittleEndian, &packet.Version)
	if err != nil {
		return errors.WithMessage(err, "failed to read HID packet")
	}
	if packet.Version != 0x1 {
		log.Warnf("received a HID packet with invalid version")
	}

	sbyte := byte(0)

	// The next two byte contain the bit encoded boolean state of the pads, true = pressed.
	/*
		Byte 2 Bit 7 (MSB) = Pad 1
		Byte 2 Bit 6       = Pad 2
		Byte 2 Bit 5       = Pad 3
		Byte 2 Bit 4       = Pad 4
		Byte 2 Bit 3       = Pad 5
		Byte 2 Bit 2       = Pad 6
		Byte 2 Bit 1       = Pad 7
		Byte 2 Bit 0       = Pad 8
	*/
	err = binary.Read(rdr, binary.BigEndian, &sbyte)
	if err != nil {
		return errors.WithMessage(err, "failed to read HID packet")
	}
	packet.unpackbools(sbyte, []button2.Button{
		button2.PadA1,
		button2.PadA2,
		button2.PadA3,
		button2.PadA4,
		button2.PadB1,
		button2.PadB2,
		button2.PadB3,
		button2.PadB4,
	})

	/*
		Byte 3 Bit 7 (MSB) = Pad 9
		Byte 3 Bit 6       = Pad 10
		Byte 3 Bit 5       = Pad 11
		Byte 3 Bit 4       = Pad 12
		Byte 3 Bit 3       = Pad 13
		Byte 3 Bit 2       = Pad 14
		Byte 3 Bit 1       = Pad 15
		Byte 3 Bit 0       = Pad 16
	*/
	err = binary.Read(rdr, binary.BigEndian, &sbyte)
	if err != nil {
		return errors.WithMessage(err, "failed to read HID packet")
	}
	packet.unpackbools(sbyte, []button2.Button{
		button2.PadC1,
		button2.PadC2,
		button2.PadC3,
		button2.PadC4,
		button2.PadD1,
		button2.PadD2,
		button2.PadD3,
		button2.PadD4,
	})

	// The boolean state for the other buttons are sent via Byte 4 & Byte 5.
	/*
		Byte 4 Bit 7 (MSB) = Shift Key
		Byte 4 But 6       = Reverse Key
		Byte 4 Bit 5       = Type Key
		Byte 4 Bit 4       = Size Key
		Byte 4 Bit 3       = Browse Key
		Byte 4 Bit 2       =
		Byte 4 Bit 1       =
		Byte 4 Bit 0       =
	*/
	err = binary.Read(rdr, binary.BigEndian, &sbyte)
	if err != nil {
		return errors.WithMessage(err, "failed to read HID packet")
	}
	packet.unpackbools(sbyte, []button2.Button{
		button2.Shift,
		button2.Reverse,
		button2.Type,
		button2.Size,
		button2.Browse,
		button2.Dial,
	})

	/*
		Byte 5 Bit 7 (MSB) = Kill Key 1
		Byte 5 Bit 6            = Kill Key 2
		Byte 5 Bit 5            = Kill Key 3
		Byte 5 Bit 4            = Kill Key 4
		Byte 5 Bit 3            = Sync Key
		Byte 5 Bit 2            = Quant Key
		Byte 5 Bit 1            = Capture Key
		Byte 5 Bit 0            =
	*/
	err = binary.Read(rdr, binary.BigEndian, &sbyte)
	if err != nil {
		return errors.WithMessage(err, "failed to read HID packet")
	}
	packet.unpackbools(sbyte, []button2.Button{
		button2.Mute1,
		button2.Mute2,
		button2.Mute3,
		button2.Mute4,
		button2.Sync,
		button2.Quant,
		button2.Capture,
	})

	/*
		Rotary Encoder

		The 6th byte contains a wrapped 0..255 value for the rotary encoder. On reset this value is 0, each clockwise step
		increments the value by 1 up to a maximum of 0xFF (255). Incrementing past 255 results in wrap around to 0 and
		decrementing through 0 wraps to 255.
	*/
	err = binary.Read(rdr, binary.BigEndian, &packet.Dial)
	if err != nil {
		return errors.WithMessage(err, "failed to read HID packet")
	}

	/*
		Analog Inputs

		The analog inputs are sent using bytes 7 thru 22. Each analog input uses two bytes in little endian format

		[TODO: check this] The first byte gives the least significant 8 bits of resolution, the second byte contains
		the most significant 4 bits of the ADC in the lower 4 bits.

		ie; a decimal value of 4000, usually represented as 0x0FA0 in hexadecimal will be sent as the byte stream  {0xA0, 0x0F}
	*/
	twobytes := uint16(0)
	for idx, _ := range packet.Filters {
		err = binary.Read(rdr, binary.LittleEndian, &twobytes)
		if err != nil {
			return errors.WithMessage(err, "failed to read HID packet")
		}

		packet.Filters[idx] = unpackuint16(twobytes)
	}
	for idx, _ := range packet.Volumes {
		err = binary.Read(rdr, binary.LittleEndian, &twobytes)
		if err != nil {
			return errors.WithMessage(err, "failed to read HID packet")
		}

		packet.Volumes[idx] = unpackuint16(twobytes)
	}

	// log.Infof("Volumes: %v, Sliders: %v", packet.Volumes, packet.Filters)
	return nil
}

// func packbools(bools []bool) (packed byte) {
// 	for bit, bol := range bools {
// 		val := 0
// 		if bol {
// 			val = 1
// 		}
// 		packed |= val << (7 - bit)
// 	}
// 	return packed
// }

func unpackuint16(data uint16) uint16 {
	first := byte(data)
	second := byte(data >> 8)
	return uint16(first) + uint16(second&0x0F)<<8
}

func (i *InState) unpackbools(zebyte byte, buttons []button2.Button) {
	for bit, btn := range buttons {
		if zebyte>>(7-bit)&0x1 != 0 {
			i.PressedButtons[btn] = button2.Pushed
		} else {
			i.PressedButtons[btn] = button2.Released
		}
	}
}
