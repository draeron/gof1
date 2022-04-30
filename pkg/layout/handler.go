package layout

import (
	"github.com/draeron/gof1/pkg/f1/button"
)

//go:generate go-enum -f=$GOFILE --noprefix

type Handler func(layout *BasicLayout, btn button.Button)
type HoldHandler func(layout *BasicLayout, btn button.Button, first bool)

/*
	HandlerType x ENUM(
	FunctionsPressed
  FunctionsHold
	FunctionsReleased
	PadPressed
  PadHold
	PadReleased
	MutePressed
	MuteHold
	MuteReleased
	DialPressed
	DialHold
	DialReleased
)
*/
type HandlerType int

func (h HandlerType) IsPressed() bool {
	switch h {
	case PadPressed, FunctionsPressed, MutePressed, DialPressed:
		return true
	default:
		return false
	}
}

func (h HandlerType) IsReleased() bool {
	switch h {
	case PadReleased, FunctionsReleased, MuteReleased, DialReleased:
		return true
	default:
		return false
	}
}

func (h HandlerType) IsHold() bool {
	switch h {
	case PadHold, FunctionsHold, MuteHold, DialHold:
		return true
	default:
		return false
	}
}
