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
	log.Error(msg).
		Str("method", msgs["method"].(string)).
		Stringer("url", msgs["url"].(fmt.Stringer)).
		Send()
}

func (retryableLogger) Info(msg string, v ...any) {
	msgs := splitMsgs(v)
	log.Info(msg).
		Str("method", msgs["method"].(string)).
		Stringer("url", msgs["url"].(fmt.Stringer)).
		Send()
}

func (retryableLogger) Debug(msg string, v ...any) {
	msgs := splitMsgs(v)
	log.Debug(msg).
		Str("method", msgs["method"].(string)).
		Stringer("url", msgs["url"].(fmt.Stringer)).
		Send()
}

func (retryableLogger) Warn(msg string, v ...any) {
	msgs := splitMsgs(v)
	log.Warn(msg).
		Str("method", msgs["method"].(string)).
		Stringer("url", msgs["url"].(fmt.Stringer)).
		Send()
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
