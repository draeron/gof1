package device

import (
	button2 "github.com/draeron/gof1/pkg/f1/button"
	event2 "github.com/draeron/gof1/pkg/f1/event"
	"github.com/draeron/gopkgs/color"
)

type State struct {
	in  InState
	out OutState
}

type PadState struct {
	color.Color
	button2.PushState
}

type ButtonState struct {
	LEDIntensity
	button2.PushState
}

type RangeState uint16

func (s *State) Copy() State {
	// copy
	st := *s

	// create and cpy new map
	st.out.Functions = map[button2.Button]LEDIntensity{}
	for k, v := range st.out.Functions {
		st.out.Functions[k] = v
	}
	return st
}

func (s *State) Pads() (states [16]PadState) {
	for _, it := range button2.Pads() {
		states[it].PushState, _ = s.in.PressedButtons[it]
		states[it].Color = s.out.Pads[it]
	}
	return
}

func (s *State) Functions() (states [16]ButtonState) {
	for _, it := range button2.Functions() {
		states[it].PushState, _ = s.in.PressedButtons[it]
		states[it].LEDIntensity = s.out.Functions[it]
	}
	return
}

func (s *State) Volumes() (states [4]RangeState) {
	for it, _ := range button2.Volumes() {
		states[it] = RangeState(s.in.Volumes[it])
	}
	return
}

func (s *State) Knobs() (states [4]RangeState) {
	for it, _ := range button2.Knobs() {
		states[it] = RangeState(s.in.Filters[it])
	}
	return
}

func (s *State) eventFromDiff(current InState) []event2.Event {
	evts := []event2.Event{}
	previous := s.in

	if current.Dial != previous.Dial {
		evt := event2.Event{
			Btn: button2.Dial,
		}

		pd := previous.Dial
		cd := current.Dial
		if (pd == 255 && cd == 0) || (pd == 0 && cd == 255) {
			pd = current.Dial
			cd = previous.Dial
		}

		if cd > pd {
			evt.Type = event2.Increment
			evt.Value = int16(current.Dial)
		} else {
			evt.Type = event2.Decrement
			evt.Value = int16(current.Dial)
		}

		evts = append(evts, evt)
	}

	for key, state := range current.PressedButtons {
		if current.PressedButtons[key] != previous.PressedButtons[key] {
			evt := event2.Event{
				Btn: key,
			}

			if state == button2.Pushed {
				evt.Type = event2.Pressed
				evt.Value = 1
			} else {
				evt.Type = event2.Released
				evt.Value = 0
			}
			evts = append(evts, evt)
		}
	}

	const maxval = 4090

	for idx, value := range current.Volumes {
		if current.Volumes[idx] != previous.Volumes[idx] {
			evts = append(evts, event2.Event{
				Btn:   button2.Volume1 + button2.Button(idx),
				Type:  event2.Changed,
				Value: int16(float64(value) / maxval * 256),
			})
		}
	}

	for idx, value := range current.Filters {
		if current.Filters[idx] != previous.Filters[idx] {
			evts = append(evts, event2.Event{
				Btn:   button2.Filter1 + button2.Button(idx),
				Type:  event2.Changed,
				Value: int16(float32(value) / maxval * 256),
			})
		}
	}

	return evts
}
