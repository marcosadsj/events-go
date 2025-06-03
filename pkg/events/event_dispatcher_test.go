package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload any
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) GetPayload() any {
	return e.Payload
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event IEvent, wg *sync.WaitGroup) {
	// Handle the event
	wg.Done()
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.event = TestEvent{Name: "TestEvent", Payload: "TestPayload"}
	suite.event2 = TestEvent{Name: "TestEvent2", Payload: "TestPayload2"}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {

	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)

	suite.NoError(err, "Expected no error when registering a handler")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected one handler registered for the event")

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)

	suite.NoError(err, "Expected no error when registering a second handler for the same event")
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected two handlers registered for the event")

	assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.handlers[suite.event.GetName()][0], "Expected first handler to be registered correctly")

	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][1], "Expected second handler to be registered correctly")
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_AlreadyRegistered() {

	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)

	suite.NoError(err, "Expected no error when registering a handler")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected one handler registered for the event")

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)

	suite.Error(err, "Expected error when registering the same handler again")
	suite.Equal(ErrHandlerAlreadyRegistered, err, "Expected specific error for already registered handler")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected still only one handler registered for the event")
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)

	suite.NoError(err, "Expected no error when registering a handler")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected one handler registered for the event")

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)

	suite.NoError(err, "Expected no error when registering a second handler for a different event")
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected one handler registered for the second event")

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)

	suite.NoError(err, "Expected no error when registering a third handler for the second event")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]), "Expected one handler registered for the second event")

	err = suite.eventDispatcher.Clear()

	suite.NoError(err, "Expected no error when clearing all handlers")
	suite.Equal(0, len(suite.eventDispatcher.handlers), "Expected all handlers to be cleared")
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {

	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)

	suite.NoError(err, "Expected no error when registering a handler")

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)

	suite.NoError(err, "Expected no error when registering a second handler for the same event")
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected two handlers registered for the event")

	hasHandler := suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler)

	suite.True(hasHandler, "Expected handler to be registered for the event")

	hasHandler = suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2)

	suite.True(hasHandler, "Expected handler2 to be registered for the event")

	hasHandler = suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler3)

	suite.False(hasHandler, "Expected handler3 not to be registered for the event")
}

type MockEventHandler struct {
	mock.Mock
}

func (m *MockEventHandler) Handle(event IEvent, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	eh := &MockEventHandler{}

	eh.On("Handle", &suite.event)

	eh2 := &MockEventHandler{}
	eh2.On("Handle", &suite.event)

	eh3 := &MockEventHandler{}
	eh3.On("Handle", &suite.event2)

	err := suite.eventDispatcher.Register(suite.event.GetName(), eh)

	suite.NoError(err, "Expected no error when registering a handler")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected one handler registered for the event")

	err = suite.eventDispatcher.Register(suite.event.GetName(), eh2)

	suite.NoError(err, "Expected no error when registering a second handler for the same event")
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected two handlers registered for the event")

	err = suite.eventDispatcher.Register(suite.event2.GetName(), eh3)

	suite.NoError(err, "Expected no error when registering a third handler for the second event")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]), "Expected one handler registered for the second event")

	err = suite.eventDispatcher.Dispatch(&suite.event)

	suite.NoError(err, "Expected no error when dispatching the event")

	err = suite.eventDispatcher.Dispatch(&suite.event)

	suite.NoError(err, "Expected no error when dispatching the event")

	err = suite.eventDispatcher.Dispatch(&suite.event2)

	suite.NoError(err, "Expected no error when dispatching a different event")

	// eh.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "Handle", 2)
	eh2.AssertNumberOfCalls(suite.T(), "Handle", 2)

	eh3.AssertNumberOfCalls(suite.T(), "Handle", 1)

	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected two handlers still registered for the event after dispatching")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]), "Expected one handler still registered for the second event after dispatching")
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.NoError(err, "Expected no error when registering a handler")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected one handler registered for the event")

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.NoError(err, "Expected no error when registering a second handler for a different event")
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected one handler registered for the second event")

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.NoError(err, "Expected no error when registering a third handler for the second event")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]), "Expected one handler registered for the second event")

	err = suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	suite.NoError(err, "Expected no error when removing a handler")
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected one handler remaining for the event after removal")

	err = suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler2)
	suite.NoError(err, "Expected no error when removing a second handler")
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event.GetName()]), "Expected no handlers remaining for the event after removing all handlers")

	err = suite.eventDispatcher.Remove(suite.event2.GetName(), &suite.handler3)
	suite.NoError(err, "Expected no error when removing a handler from a different event")
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event2.GetName()]), "Expected no handlers remaining for the second event after removing all handlers")

	err = suite.eventDispatcher.Remove(suite.event2.GetName(), &suite.handler)
	suite.Error(err, "Expected error when trying to remove a handler that is not registered for the event")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
