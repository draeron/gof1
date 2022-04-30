package f1

import (
	"github.com/draeron/gof1/pkg/f1/button"
	"github.com/draeron/gof1/pkg/f1/event"
	"github.com/draeron/gopkgs/color"
)

type Controller interface {
	EnableDebugLogger()
	Close()
	SetPadColors(sets button.ColorMap) error
	Subscribe(channel chan<- event.Event)
	String() string
	Name() string
}

type Colorer interface {
	SetPadColorAll(col color.Color) error
	SetPadColorMany(btns []button.Button, color color.Color) error
	SetPadColor(btn button.Button, color color.Color) error
	SetPadColors(sets button.ColorMap) error
}
