package events

import "time"

type IEvent interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() any
}

type IEventHandler interface {
	Handle(event IEvent)
}

type IEventDispatcher interface {
	Register(eventName string, handler IEventHandler) error
	Dispatch(event IEvent) error
	Remove(eventName string, handler IEventHandler) error
	Has(eventName string, handler IEventHandler) bool
	Clear() error
}
