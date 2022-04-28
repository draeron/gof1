package f1

import (
	event "github.com/draeron/gof1/pkg/f1/event"
)

type EventCallBack func(event.Event)

func (c *Controller) AddCallback(filter event.Filter, cb EventCallBack) {
	go func() {
		input := make(chan event.Event, 10)
		c.Subscribe(input)
		for evt := range input {
			if filter(evt) {
				cb(evt)
			}
		}
	}()
}

func (c *Controller) Subscribe(channel chan<- event.Event) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Infof("adding new event suscriber")
	c.subscribers = append(c.subscribers, channel)
}

func (c *Controller) sendToSubscribers(evt event.Event) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for _, channel := range c.subscribers {
		select {
		case channel <- evt:
		default:
			// full or closed channel
		}
	}
}
