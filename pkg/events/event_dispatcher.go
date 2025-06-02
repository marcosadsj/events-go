package events

import (
	"errors"
	"slices"
)

var (
	ErrHandlerAlreadyRegistered = errors.New("handler already registered for this event")
)

type EventDispatcher struct {
	handlers map[string][]IEventHandler
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]IEventHandler),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler IEventHandler) error {

	if _, ok := ed.handlers[eventName]; !ok {
		ed.handlers[eventName] = []IEventHandler{}
	}

	if slices.Contains(ed.handlers[eventName], handler) {
		return ErrHandlerAlreadyRegistered
	}

	ed.handlers[eventName] = append(ed.handlers[eventName], handler)

	return nil
}
