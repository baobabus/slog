// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"errors"
	"io/ioutil"
	"log"
)

var (
	errSuccess     = errors.New("success")
	errEllipsis    = errors.New("...")
	errNoFacility  = errors.New("no facility")
	errNoFormatter = errors.New("no formatter")
)

type sLogger struct {
	level     Priority
	formatter Formatter
	logs      map[Priority]Log
}

func New(facility Facility, level Priority, formatter Formatter) (Logger, error) {
	if facility == nil {
		return nil, errNoFacility
	}
	if formatter == nil {
		return nil, errNoFormatter
	}
	ls, err := facility.OpenLogs(level)
	if err != nil {
		return nil, err
	}
	logs := make(map[Priority]Log, len(ls))
	for p, l := range ls {
		logs[p] = &sLog{formatter: formatter, logger: l, scope: nil}
	}
	return &sLogger{level: level, formatter: formatter, logs: logs}, err
}

func (this *sLogger) Level() Priority {
	return this.level
}

func (this *sLogger) Formatter() Formatter {
	return this.formatter
}

func (this *sLogger) Log(pri Priority) Log {
	return this.logs[pri.Bound()]
}

func (this *sLogger) Info() Log {
	return this.logs[PriorityInfo]
}

func (this *sLogger) Notice() Log {
	return this.logs[PriorityNotice]
}

func (this *sLogger) Warning() Log {
	return this.logs[PriorityWarn]
}

func (this *sLogger) Error() Log {
	return this.logs[PriorityError]
}

func (this *sLogger) Trace(detail int) Log {
	if PriorityTrace+Priority(detail-1) <= this.level {
		return this.logs[PriorityTrace]
	} else {
		return drain
	}
}

func (this *sLogger) On(err error) Selector {
	return &sSelector{this, err}
}

func (this *sLogger) With(err error) Selector {
	if err == nil {
		err = errSuccess
	}
	return &sSelector{this, err}
}

type sLog struct {
	formatter Formatter
	logger    *log.Logger
	scope     error
}

func (this *sLog) Printe(message string, v ...interface{}) {
	this.prints(2, message, v, errEllipsis)
}

func (this *sLog) Prints(message string, v ...interface{}) {
	this.prints(2, message, v, this.scope)
}

func (this *sLog) Logger() *log.Logger {
	return this.logger
}

func (this *sLog) ScopedLog(err error) Log {
	if err == nil {
		return drain
	}
	return &sLog{formatter: this.formatter, logger: this.logger, scope: err}
}

func (this *sLog) prints(calldepth int, message string, v []interface{}, err error) error {
	if err == nil {
		err = this.scope
	}
	s := message
	if this.formatter != nil {
		s = this.formatter(message, v, err)
	}
	return this.logger.Output(calldepth+1, s)
}

var drain = &sLog{formatter: nil, logger: log.New(ioutil.Discard, "", 0), scope: nil}

type sSelector struct {
	*sLogger
	scope error
}

func (this *sSelector) Log(pri Priority) Log {
	return this.logs[pri.Bound()].ScopedLog(this.scope)
}

func (this *sSelector) Info() Log {
	return this.logs[PriorityInfo].ScopedLog(this.scope)
}

func (this *sSelector) Notice() Log {
	return this.logs[PriorityNotice].ScopedLog(this.scope)
}

func (this *sSelector) Warning() Log {
	return this.logs[PriorityWarn].ScopedLog(this.scope)
}

func (this *sSelector) Error() Log {
	return this.logs[PriorityError].ScopedLog(this.scope)
}

func (this *sSelector) Trace(detail int) Log {
	if PriorityTrace+Priority(detail-1) <= this.level {
		return this.logs[PriorityTrace].ScopedLog(this.scope)
	} else {
		return drain
	}
}

func (this *sSelector) Prints(message string, v ...interface{}) {
	this.scopedLog().prints(2, message, v, nil)
}

func (this *sSelector) Logger() *log.Logger {
	return this.scopedLog().Logger()
}

func (this *sSelector) scopedLog() Log {
	if this.scope == nil || this.scope == errSuccess || this.scope == errEllipsis {
		return this.logs[PriorityInfo].ScopedLog(this.scope)
	} else {
		return this.logs[PriorityError].ScopedLog(this.scope)
	}
}
