package main

import (
	"fmt"
	"os"

	"go.arsenm.dev/logger"
	"go.arsenm.dev/logger/log"
)

func init() {
	l := logger.NewPretty(os.Stderr)

	if os.Getenv("LEMMY_REPLY_BOT_DEBUG") == "1" {
		l.Level = logger.LogLevelDebug
	}

	log.Logger = l
}

type retryableLogger struct{}

func (retryableLogger) Error(msg string, v ...any) {
	msgs := splitMsgs(v)
	evt := log.Error(msg)
	sendEvt(evt, msgs)
}

func (retryableLogger) Info(msg string, v ...any) {
	msgs := splitMsgs(v)
	evt := log.Info(msg)
	sendEvt(evt, msgs)

}

func (retryableLogger) Debug(msg string, v ...any) {
	msgs := splitMsgs(v)
	evt := log.Debug(msg)
	sendEvt(evt, msgs)
}

func (retryableLogger) Warn(msg string, v ...any) {
	msgs := splitMsgs(v)
	evt := log.Warn(msg)
	sendEvt(evt, msgs)
}

func splitMsgs(v []any) map[string]any {
	out := map[string]any{}

	for i, val := range v {
		if (i+1)%2 == 0 {
			continue
		}

		out[val.(string)] = v[i+1]
	}

	return out
}

func sendEvt(evt logger.LogBuilder, msgs map[string]any) {
	for name, val := range msgs {
		switch val := val.(type) {
		case int:
			evt = evt.Int(name, val)
		case string:
			evt = evt.Str(name, val)
		case fmt.Stringer:
			evt = evt.Stringer(name, val)
		default:
			evt = evt.Any(name, val)
		}
	}
	evt.Send()
}
