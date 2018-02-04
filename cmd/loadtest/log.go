package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
)

// MultiLogger logs to multiple log.Logger instances on a single Log() call.
type MultiLogger struct {
	log.Logger

	loggers []log.Logger
}

func NewMultiLogger(loggers ...log.Logger) *MultiLogger {
	return &MultiLogger{
		loggers: append([]log.Logger{}, loggers...),
	}
}

func (l *MultiLogger) Log(keyvals ...interface{}) error {
	kvs := keyvals[:len(keyvals):len(keyvals)]
	if len(kvs)%2 != 0 {
		kvs = append(kvs, log.ErrMissingValue)
	}

	var wg sync.WaitGroup
	for _, l := range l.loggers {
		wg.Add(1)
		go func(l log.Logger) {
			defer wg.Done()

			l.Log(kvs...)
		}(l)
	}
	wg.Wait()

	return nil
}

// InitLoggers is a helper function that initializes two log.Logger instances,
// which log to stdout and a log file in JSON format.
func InitLoggers(logFile io.Writer, debug bool) log.Logger {
	logfmt := term.NewLogger(log.NewSyncWriter(os.Stdout), log.NewLogfmtLogger, colorFn)
	logjson := log.NewJSONLogger(log.NewSyncWriter(logFile))

	var logger log.Logger
	logger = NewMultiLogger(logfmt, logjson)
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
	if debug {
		logger = level.NewFilter(logger, level.AllowDebug())
	} else {
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return logger
}

func colorFn(keyvals ...interface{}) term.FgBgColor {
	for i := 0; i < len(keyvals)-1; i += 2 {
		if keyvals[i] != "level" {
			continue
		}
		s, ok := keyvals[i+1].(fmt.Stringer)
		if !ok {
			return term.FgBgColor{}
		}
		switch s.String() {
		case "debug":
			return term.FgBgColor{Fg: term.Gray}
		case "info":
			return term.FgBgColor{Fg: term.White}
		case "warn":
			return term.FgBgColor{Fg: term.Yellow}
		case "error":
			return term.FgBgColor{Fg: term.Red}
		default:
			return term.FgBgColor{}
		}
	}
	return term.FgBgColor{}
}
