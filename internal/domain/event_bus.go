package domain

import (
	"go.uber.org/fx"
)

type Event interface{}

type EventBus struct {
	messageHandlers map[string][]func(event Event) error
}

func NewEventBus() *EventBus {
	return &EventBus{
		messageHandlers: make(map[string][]func(event Event) error),
	}
}

func (bus *EventBus) On(name string, handler func(event Event) error) {
	if handlers, ok := bus.messageHandlers[name]; ok {
		bus.messageHandlers[name] = append(handlers, handler)
	} else {
		bus.messageHandlers[name] = []func(event Event) error{handler}
	}
}
func (bus *EventBus) Emit(name string, event Event) {
	if handlers, ok := bus.messageHandlers[name]; ok {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

var Module = fx.Options(
	fx.Provide(NewEventBus),
)
