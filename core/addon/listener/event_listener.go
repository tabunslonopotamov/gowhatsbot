package listener

import (
	"main/core/mylog"
	"main/core/whats"
	"reflect"
)

func Call(i interface{}, gwb *whats.GoWhatsBot) {
	var eventString = reflect.TypeOf(i).String()

	for _, eventType := range []string{eventString, "*"} {
		if listener, ok := EventListener[eventType]; ok {
			for name, listen := range listener {
				if err := listen(i, gwb); err != nil {
					mylog.Error("EventListener", eventType, name)
				}
			}
		}
	}
}

type Listener func(interface{}, *whats.GoWhatsBot) error

var EventListener map[string][]Listener = map[string][]Listener{}

func AddEvent(eventType string, listener Listener) {

	if _, ok := EventListener[eventType]; !ok {
		EventListener[eventType] = []Listener{}
	}

	EventListener[eventType] = append(EventListener[eventType], listener)
}

func AddEventTypeOf(t interface{}, listener Listener) {
	var eventString = reflect.TypeOf(t).String()

	AddEvent(eventString, listener)
}
