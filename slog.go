// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"log"
)

type Accessor interface {
	Log(pri Priority) Log
	Info() Log
	Notice() Log
	Warning() Log
	Error() Log
	Trace(detail int) Log
}

type Logger interface {
	Accessor
	Level() Priority
	Formatter() Formatter
	On(err error) Selector
	With(err error) Selector
}

type Log interface {
	Printe(message string, v ...interface{})
	Prints(message string, v ...interface{})
	Fatals(message string, v ...interface{})
	Logger() *log.Logger
	ScopedLog(err error) Log
	prints(calldepth int, message string, v []interface{}, err error) error
	// Shortcuts to log.Logger
	Output(calldepth int, s string) error
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
}

type Selector interface {
	Accessor
	Prints(message string, v ...interface{})
	Fatals(message string, v ...interface{})
	Logger() *log.Logger
}

// Logging level expressed as a number.
// Positive numbers correspond to trace detail level.
// Default is value corresponds to Info logging level.
type Priority int

const (
	PriorityError Priority = iota
	PriorityWarn
	PriorityNotice
	PriorityInfo
	PriorityTrace
)

var priTags = map[Priority]string{
	PriorityError:  "ERROR ",
	PriorityWarn:   "WARNING ",
	PriorityNotice: "NOTICE ",
	PriorityInfo:   "INFO ",
	PriorityTrace:  "TRACE ",
}

const prioritiesCount = 5

func (this Priority) Bound() Priority {
	switch {
	case this < PriorityError:
		return PriorityError
	case this > PriorityTrace:
		return PriorityTrace
	}
	return this
}

func (this Priority) Tag() string {
	return priTags[this.Bound()]
}

func Info() Log {
	return SharedLogger().Info()
}

func Notice() Log {
	return SharedLogger().Notice()
}

func Warning() Log {
	return SharedLogger().Warning()
}

func Error() Log {
	return SharedLogger().Error()
}

func Trace(detail int) Log {
	return SharedLogger().Trace(detail)
}

func On(err error) Selector {
	return SharedLogger().On(err)
}

func With(err error) Selector {
	return SharedLogger().With(err)
}
