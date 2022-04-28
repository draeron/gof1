package event

import (
	"github.com/draeron/gof1/pkg/f1/button"
)

type Filter func(Event) bool

func IsOfType(types ...Type) Filter {
	return func(event Event) bool {
		for _, tp := range types {
			if event.Type == tp {
				return true
			}
		}
		return false
	}
}

func IsButtonOfType(types ...button.Button) Filter {
	return func(event Event) bool {
		for _, tp := range types {
			if event.Btn == tp {
				return true
			}
		}
		return false
	}
}

func FilterChannel(input <-chan Event, output chan Event, filter Filter) <-chan Event {
	if output == nil {
		output = make(chan Event)
	}
	go func() {
		for it := range input {
			if filter(it) {
				output <- it
			}
		}
	}()

	return output
}
