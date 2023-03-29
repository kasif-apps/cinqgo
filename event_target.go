package cinqgo

type Event struct {
	EventType string
	Detail    interface{}
}

type Listener struct {
	EventType  string
	CallbackFn *func(e Event)
}

type EventTarget struct {
	listeners *[]Listener
}

func (target EventTarget) AddEventListener(event_type string, callback_fn *func(e Event)) {
	listener := Listener{
		EventType:  event_type,
		CallbackFn: callback_fn,
	}

	*target.listeners = append(*target.listeners, listener)
}

func (target EventTarget) RemoveEventListener(event_type string, callback_fn *func(e Event)) {
	listeners := *target.listeners
	for i, listener := range *target.listeners {
		if listener.EventType == event_type && listener.CallbackFn == callback_fn {
			listeners = remove(listeners, i)
		}
	}

	*target.listeners = listeners
}

func (target EventTarget) DispatchEvent(e Event) {
	for _, listener := range *target.listeners {
		if listener.EventType == e.EventType {
			fn := *listener.CallbackFn
			fn(e)
		}
	}
}

func NewEventTarget() EventTarget {
	return EventTarget{
		listeners: &[]Listener{},
	}
}

func NewEvent(event_type string, detail interface{}) Event {
	return Event{
		EventType: event_type,
		Detail:    detail,
	}
}

func remove[T interface{}](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
