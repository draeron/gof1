package device

import (
	event2 "github.com/draeron/gof1/pkg/f1/event"
)

type EventCallBack func(event2.Event)

func (d *Device) AddCallback(filter event2.Filter, cb EventCallBack) {
	go func() {
		input := make(chan event2.Event, 10)
		d.Subscribe(input)
		for evt := range input {
			if filter(evt) {
				cb(evt)
			}
		}
	}()
}

func (d *Device) sendToSubscribers(evt event2.Event) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	for _, channel := range d.subscribers {
		select {
		case channel <- evt:
		default:
			// full or closed channel
		}
	}
}
