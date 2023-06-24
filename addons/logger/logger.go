package logger

import (
	"main/core/addon/listener"
	"main/core/mylog"
	"main/core/whats"

	"go.mau.fi/whatsmeow/types/events"
)

func init() {

	listener.AddEvent("*", printLog)
}

func printLog(i interface{}, gwb *whats.GoWhatsBot) error {
	if !whats.EventIs(i, &events.Receipt{}) {
		mylog.Debug(">", whats.EventToString(i))
	}

	return nil
}
