package button

import (
	"sort"
)

//go:generate go-enum -f=$GOFILE --noprefix

/*
Button x ENUM(
	Filter1
	Filter2
	Filter3
	Filter4

	Volume1
	Volume2
	Volume3
	Volume4

	PadA1
	PadA2
	PadA3
	PadA4

	PadB1
	PadB2
	PadB3
	PadB4

	PadC1
	PadC2
	PadC3
	PadC4

	PadD1
	PadD2
	PadD3
	PadD4

	Mute1
	Mute2
	Mute3
	Mute4

	Dial
	SevenSegment

	Sync
	Quant
	Capture
	Reverse
	Type
	Size
	Browse
) */
type Button int

type Buttons []Button

func (b Button) Type() BtnType {
	switch {
	case b.IsPad(), b.IsMute(), b.IsFunctions(), b == SevenSegment:
		return Push

	case b.IsFader(), b.IsKnob():
		return Absolute

	case b == Dial:
		return Relative
	}
	panic("unknown button type")
}

func (b Button) IsPad() bool {
	return b >= PadA1 && b <= PadD4
}

func (b Button) IsKnob() bool {
	return b >= Filter1 && b <= Filter4
}

func (b Button) IsMute() bool {
	return b >= Mute1 && b <= Mute4
}

func (b Button) IsFader() bool {
	return b >= Volume1 && b <= Volume4
}

func (b Button) IsFunctions() bool {
	return b >= Sync && b <= Browse
}

func Values() (s Buttons) {
	for _, b := range _ButtonValue {
		s = append(s, b)
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i] > s[j]
	})
	return
}

func Pads() (s []Button) {
	for b := PadA1; b <= PadD4; b++ {
		s = append(s, b)
	}
	return
}

func Mutes() (s []Button) {
	for b := Mute1; b <= Mute4; b++ {
		s = append(s, b)
	}
	return
}

func Functions() (s []Button) {
	for b := Sync; b <= Browse; b++ {
		s = append(s, b)
	}
	return
}

func Volumes() (s []Button) {
	for b := Volume1; b <= Volume4; b++ {
		s = append(s, b)
	}
	return
}

func Knobs() (s []Button) {
	for b := Filter1; b <= Filter4; b++ {
		s = append(s, b)
	}
	return
}

func FromXY(x, y int) Button {
	if x < 0 || y < 0 || x > 4 || y > 4 {
		return Button(-1)
	}
	return PadA1 + Button(y*4)
}
