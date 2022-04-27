package event

import (
	"fmt"

	"github.com/draeron/gof1/pkg/f1/button"
)

type Event struct {
	Type  Type
	Btn   button.Button
	Value int16
}

func (e Event) String() string {
	str := fmt.Sprintf("Event: %s - %s", e.Btn, e.Type)
	if e.Type == Changed || e.Type == Increment || e.Type == Decrement {
		str += fmt.Sprintf(" - %v", e.Value)
	}
	return str
}
