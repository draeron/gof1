package layout

import (
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/draeron/gof1/pkg/f1"
	"github.com/draeron/gof1/pkg/f1/button"
	"github.com/draeron/gof1/pkg/f1/event"
	"github.com/draeron/gopkgs/color"
)

type Layout interface {
	Connect(controller f1.Controller)
	Disconnect()
	Activate()
	Deactivate()
}

type BasicLayout struct {
	DebugName  string
	state      f1.ButtonStateMap
	lastColors button.ColorMap
	controler  f1.Controller
	handlers   handlersMap
	enabled    atomic.Bool
	eventsCh   chan (event.Event)
	mask       f1.Mask
	mutex      sync.RWMutex
	ticker     *time.Ticker

	holdTimer        map[HandlerType]time.Duration
	holdTimerDefault time.Duration
}

type handlersMap map[HandlerType]HoldHandler

const DefaultHoldDuration = time.Millisecond * 250

func NewLayoutPreset(preset f1.MaskPreset) *BasicLayout {
	return NewLayout(preset.Mask())
}

func NewLayout(mask f1.Mask) *BasicLayout {
	l := &BasicLayout{
		state:            f1.NewButtonStateMap(),
		lastColors:       button.ColorMap{},
		handlers:         handlersMap{},
		mask:             mask,
		holdTimerDefault: DefaultHoldDuration,
		holdTimer:        map[HandlerType]time.Duration{},
	}
	l.state.SetColors(mask, color.Black) // allocated state
	return l
}

func (l *BasicLayout) Connect(controller f1.Controller) {
	l.mutex.Lock()
	l.controler = controller
	l.mutex.Unlock()

	if l.DebugName != "" {
		log.Infof("connecting layout %s to controller %s", l.DebugName, controller.Name())
	}

	go l.tickEvents()
	go l.tickUpdate()
}

func (l *BasicLayout) Disconnect() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.DebugName != "" {
		log.Infof("disconnecting layout %s from controller %s", l.DebugName, l.controler.Name())
	}

	close(l.eventsCh)
	l.controler = nil
	l.ticker.Stop()
	l.ticker = nil
}

/*
	When enabling a layout, it will transfert it's color state to
*/
func (l *BasicLayout) Activate() {
	l.enabled.Store(true)
}

/*
	When disabling a layout, any pressed state will be deleted
*/
func (l *BasicLayout) Deactivate() {
	l.enabled.Store(false)

	l.mutex.Lock()
	defer l.mutex.Unlock()

	// clear last displayed
	l.lastColors = button.ColorMap{}
	l.state.ResetPressed()
}

/*
	The handler will be
*/
func (l *BasicLayout) SetHandler(htype HandlerType, handler Handler) {
	l.handlers[htype] = func(layout *BasicLayout, btn button.Button, first bool) {
		if first {
			handler(layout, btn)
		}
	}
}

func (l *BasicLayout) SetHandlerHold(htype HandlerType, handler HoldHandler) {
	l.handlers[htype] = handler
}

func (l *BasicLayout) SetHoldTimer(htype HandlerType, duration time.Duration) {
	l.holdTimer[htype] = duration
}

func (l *BasicLayout) SetDefaultHoldTimer(duration time.Duration) {
	l.holdTimerDefault = duration
}

func (l *BasicLayout) HoldTime(btn button.Button) time.Duration {
	return l.state.HoldTime(btn)
}

func (l *BasicLayout) IsPressed(btn button.Button) bool {
	return l.state.IsPressed(btn)
}

func (l *BasicLayout) IsHold(btn button.Button, threshold time.Duration) bool {
	return l.state.IsHold(btn, threshold)
}

func (l *BasicLayout) UpdateDevice() error {
	if l.enabled.Load() {
		l.mutex.Lock()
		defer l.mutex.Unlock()

		if l.controler == nil {
			return nil
		}

		colors := l.mask.Intersect(l.state).DiffFrom(l.lastColors)

		if len(colors) > 0 {
			err := l.controler.SetPadColors(colors)
			l.lastColors = l.lastColors.ApplyFrom(colors)
			return err
		}
	}
	return nil
}

func (l *BasicLayout) Color(btn button.Button) color.Color {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.state.Color(btn)
}

func (l *BasicLayout) SetColorAll(col color.Color) error {
	l.state.SetAllColors(col)
	return nil
}

func (l *BasicLayout) SetColorMask(mask f1.MaskPreset, col color.Color) error {
	for b, _ := range mask.Mask() {
		l.state.SetColor(b, col)
	}
	return nil
}

func (l *BasicLayout) SetColorMany(btns []button.Button, color color.Color) error {
	for _, k := range btns {
		l.state.SetColor(k, color)
	}
	return nil
}

func (l *BasicLayout) SetColor(btn button.Button, color color.Color) error {
	l.state.SetColor(btn, color)
	return nil
}

func (l *BasicLayout) SetColors(set button.ColorMap) error {
	l.state.SetColorsMap(set)
	return nil
}

func (l *BasicLayout) tickEvents() {
	l.mutex.Lock()
	l.eventsCh = make(chan event.Event, 20)
	l.controler.Subscribe(l.eventsCh)
	l.mutex.Unlock()

	for e := range l.eventsCh {
		l.dispatch(e)
	}
}

func (l *BasicLayout) tickUpdate() {
	l.mutex.Lock()
	l.ticker = time.NewTicker(time.Second / 60)
	l.mutex.Unlock()

	for range l.ticker.C {
		l.UpdateDevice()
	}
}

func (l *BasicLayout) dispatch(e event.Event) {
	if !l.enabled.Load() || !l.mask[e.Btn] {
		return
	}
	var ht HandlerType

	switch {
	case e.Btn.IsPad():
		if e.Type == event.Pressed {
			ht = PadPressed
		} else {
			ht = PadReleased
		}

	case e.Btn.IsMute():
		if e.Type == event.Pressed {
			ht = MutePressed
		} else {
			ht = MuteReleased
		}

	case e.Btn.IsFunctions():
		if e.Type == event.Pressed {
			ht = FunctionsPressed
		} else {
			ht = FunctionsReleased
		}

	case e.Btn == button.Dial:
		if e.Type == event.Pressed {
			ht = DialPressed
		} else {
			ht = DialReleased
		}
	}

	l.mutex.Lock()
	if e.Type == event.Pressed {
		l.state.Press(e.Btn)
	}
	l.mutex.Unlock()

	if e.Type == event.Pressed {
		if handle, ok := l.handlers[ht+1]; ok {
			timer := l.holdTimerDefault
			if t, ok := l.holdTimer[ht+1]; ok {
				timer = t
			}
			go func() {
				first := true
				for {
					<-time.After(timer)
					if l.state.IsHold(e.Btn, timer) {
						handle(l, e.Btn, first)
						if first {
							first = false
						}
					} else {
						return
					}
				}
			}()
		}
	}

	if h, ok := l.handlers[ht]; ok {
		h(l, e.Btn, true)
	}

	l.mutex.Lock()
	if e.Type == event.Released {
		l.state.Release(e.Btn)
	}
	l.mutex.Unlock()
}
