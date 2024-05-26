package logger

import (
	"book-catalog/internal/app/common"
	"book-catalog/internal/config"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const requestID = "requestID"

// LoggerImpl is a struct that represents a logger
type LoggerImpl struct {
	log *zerolog.Logger
}

// Nop returns a no operation logger
func Nop() *LoggerImpl {
	nop := zerolog.Nop()
	return &LoggerImpl{
		log: &nop,
	}
}

// NewLogger creates a new logger
func NewLogger(cfg *config.Config) *LoggerImpl {
	zerolog.SetGlobalLevel(cfg.Log.Level)

	var logger zerolog.Logger

	if cfg.Log.JSON {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	return &LoggerImpl{
		log: &logger,
	}
}

// New creates a new logger
func (l *LoggerImpl) New(name string) Logger {
	newLogger := l.log.With().Str("logger", name).Logger()
	return &LoggerImpl{
		log: &newLogger,
	}
}

// Logger returns the logger
func (l *LoggerImpl) Logger() zerolog.Logger {
	return *l.log
}

// Inf returns an info event
func (l *LoggerImpl) Inf() Event {
	return &ServiceEvent{
		Event: l.log.Info(), //nolint:zerologlint // not relevant
	}
}

// Wrn returns a warning event
func (l *LoggerImpl) Wrn() Event {
	return &ServiceEvent{
		Event: l.log.Warn(), //nolint:zerologlint // not relevant
	}
}

// Err returns an error event
func (l *LoggerImpl) Err(err error) Event {
	return &ServiceEvent{
		Event: l.log.Error().Err(err), //nolint:zerologlint // not relevant
	}
}

// Dbg returns a debug event
func (l *LoggerImpl) Dbg() Event {
	return &ServiceEvent{
		Event: l.log.Debug(), //nolint:zerologlint // not relevant
	}
}

// Trc returns a trace event
func (l *LoggerImpl) Trc() Event {
	return &ServiceEvent{
		Event: l.log.Trace(), //nolint:zerologlint // not relevant
	}
}

// Ftl returns a fatal event
func (l *LoggerImpl) Ftl() Event {
	return &ServiceEvent{
		Event: l.log.Fatal(), //nolint:zerologlint // not relevant
	}
}

// ServiceEvent is a struct that represents a logger event
type ServiceEvent struct {
	*zerolog.Event
}

// Ctx returns an event with context
func (e *ServiceEvent) Ctx(ctx context.Context) Event {
	id := common.GetRequestID(ctx)
	if id != nil {
		e.Event = e.Event.Str(requestID, id.(string))
	}
	return e
}

// Msg returns a message
func (e *ServiceEvent) Msg(msg string, args ...any) {
	e.Event.Msgf(msg, args...)
}

// Err returns an event with error
func (e *ServiceEvent) Err(err error) Event {
	e.Event = e.Event.Err(err)
	return e
}

// Values returns an event with pair (key=value)
func (e *ServiceEvent) Values(values ...any) Event {
	if len(values) == 0 {
		return e
	}

	for i := 0; i < len(values); i += 2 {
		name := fmt.Sprintf("%v", values[i])
		var value any
		if i+1 < len(values) {
			value = values[i+1]
		} else {
			value = nil
		}
		e.Event = e.Event.Interface(name, value)
	}

	return e
}

// Fatalf returns a fatal event with formatted message
func (l *LoggerImpl) Fatalf(format string, v ...any) {
	l.log.Fatal().Msgf(format, v...)
}

// Printf returns an event with formatted message
func (l *LoggerImpl) Printf(format string, v ...any) {
	l.log.Info().Msgf(format, v...)
}
