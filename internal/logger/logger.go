package logger

import (
	"context"

	"github.com/rs/zerolog"
)

// Logger is an interface that represents a logger
type Logger interface {
	New(name string) Logger
	Logger() zerolog.Logger

	Trc() Event
	Dbg() Event
	Inf() Event
	Wrn() Event
	Err(err error) Event
	Ftl() Event

	// goose
	Fatalf(format string, v ...any)
	Printf(format string, v ...any)
}

type Event interface {
	Ctx(ctx context.Context) Event
	Err(err error) Event
	Values(values ...any) Event
	Msg(msg string, args ...any)
}
