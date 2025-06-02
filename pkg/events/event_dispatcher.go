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

func (ed *EventDispatcher) Dispatch(event IEvent) error {
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		for _, handler := range handlers {
			handler.Handle(event)
		}
		return nil
	}
	return nil // No handlers registered for this event
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

func (ed *EventDispatcher) Clear() error {

	ed.handlers = make(map[string][]IEventHandler)

	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler IEventHandler) bool {

	if handlers, ok := ed.handlers[eventName]; ok {
		return slices.Contains(handlers, handler)
	}

	return false
}
