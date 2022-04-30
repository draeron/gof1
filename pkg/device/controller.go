package device

import (
	"github.com/draeron/gof1/pkg/f1/event"
)

func (d *Device) EnableDebugLogger() {
	go func() {
		log.Debugf("enable debug logging of events")
		ch := make(chan event.Event, 20)
		d.Subscribe(ch)
		for evt := range ch {
			log.Debugf(evt.String())
		}
	}()
}

func (d *Device) Subscribe(channel chan<- event.Event) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	log.Infof("adding new event suscriber")
	d.subscribers = append(d.subscribers, channel)
}

func (d *Device) Close() {
	if d.device != nil {
		d.device.Close()
	}
}

func (d *Device) Name() string {
	return F1ProductName
}

func (d *Device) String() string {
	return F1ProductName
}
