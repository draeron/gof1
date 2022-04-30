package button

import (
	"github.com/draeron/gopkgs/color"
)

type ColorMap map[Button]color.Color

func (c ColorMap) ApplyFrom(other ColorMap) ColorMap {
	for k, v := range other {
		c[k] = v
	}
	return c
}

func (c ColorMap) DiffFrom(cmap ColorMap) ColorMap {
	out := ColorMap{}

	for btn, col := range c {
		if other, ok := cmap[btn]; ok {
			if !col.Equal(other) {
				out[btn] = col
			}
		} else { // missing color considered changed too
			out[btn] = col
		}
	}
	return out
}
