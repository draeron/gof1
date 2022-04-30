package f1

import (
	"github.com/draeron/gof1/pkg/f1/button"
	"github.com/draeron/gopkgs/color/7bits"
)

//go:generate go-enum -f=$GOFILE --noprefix
/*
	Mask x ENUM(
	MaskKnobs
	MaskPads
	MaskFunctions
	MaskVolumes
	MaskMutes
	MaskAll
)
*/
type MaskPreset int

type Mask map[button.Button]bool

/*
	Remove
*/
func (m Mask) Intersect(mapp ButtonStateMap) button.ColorMap {
	out := make(button.ColorMap)
	for k, v := range m {
		if v {
			if cl := mapp.Get(k); cl != nil {
				out[k] = seven_bits.FromColor(cl.Color)
			}
		}
	}
	return out
}

func (m Mask) MergePreset(masks ...MaskPreset) Mask {
	out := m
	for _, mask := range masks {
		out.Merge(mask.Mask())
	}
	return out
}

func (m Mask) Merge(masks ...Mask) Mask {
	out := m
	for _, mask := range masks {
		for b, v := range mask {
			out[b] = v
		}
	}
	return out
}

func (mp MaskPreset) Mask() Mask {
	m := Mask{}

	switch mp {
	case MaskAll:
		for _, b := range button.Values() {
			m[b] = true
		}
	case MaskFunctions:
		for _, b := range button.Functions() {
			m[b] = true
		}
	case MaskPads:
		for _, b := range button.Pads() {
			m[b] = true
		}
	case MaskKnobs:
		for _, b := range button.Knobs() {
			m[b] = true
		}
	case MaskVolumes:
		for _, b := range button.Volumes() {
			m[b] = true
		}
	case MaskMutes:
		for _, b := range button.Mutes() {
			m[b] = true
		}
	}

	return m
}
