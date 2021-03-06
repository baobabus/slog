// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	errSuccess     = errors.New("success")
	errEllipsis    = errors.New("...")
	errNoFacility  = errors.New("no facility")
	errNoFormatter = errors.New("no formatter")
)

type sLogger struct {
	facility  Facility
	level     Priority
	formatter Formatter
	logs      map[Priority]Log
}

func New(facility Facility, level Priority, formatter Formatter, filter []string) (Logger, error) {
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
		if p < PriorityTrace {
			logs[p] = &sLog{formatter: formatter, logger: l, scope: nil}
		} else {
			logs[p] = &sLog{formatter: formatter, logger: l, filter: filter, scope: nil}
		}
	}
	return &sLogger{facility: facility, level: level, formatter: formatter, logs: logs}, err
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

func (this *sLogger) On(err ...error) Selector {
	return &sSelector{this, err}
}

func (this *sLogger) Success() Selector {
	return &sSelector{this, []error{errSuccess}}
}

func (this *sLogger) With(err ...error) Selector {
	if err == nil || len(err) == 0 || (len(err) == 1 && err[0] == nil) {
		err = []error{errSuccess}
	}
	return &sSelector{this, err}
}

type sLog struct {
	formatter Formatter
	logger    *log.Logger
	filter    []string
	scope     []error
	soff      int
}

func (this *sLog) Printe(message string, v ...interface{}) {
	this.prints(2, message, v, []error{errEllipsis})
}

func (this *sLog) Prints(message string, v ...interface{}) {
	this.prints(2, message, v, this.scope)
}

func (this *sLog) Fatals(message string, v ...interface{}) {
	this.prints(2, message, v, this.scope)
	os.Exit(1)
}

func (this *sLog) Logger() *log.Logger {
	return this.filteredLogger(2)
}

func (this *sLog) ScopedLog(err ...error) Log {
	if err == nil || len(err) == 0 || (len(err) == 1 && err[0] == nil) {
		return drain
	}
	return &sLog{formatter: this.formatter, logger: this.logger, filter: this.filter, scope: err}
}

func (this *sLog) Offset(stackOffset int) Log {
	return &sLog{formatter: this.formatter, logger: this.logger, filter: this.filter, scope: this.scope, soff: this.soff + stackOffset}
}

func (this *sLog) filteredLogger(calldepth int) *log.Logger {
	if len(this.filter) == 0 {
		return this.logger
	}
	if _, file, _, ok := runtime.Caller(calldepth + this.soff); ok {
		for _, f := range this.filter {
			if strings.HasSuffix(file, f) {
				return this.logger
			}
		}
	}
	return dscrd
}

func (this *sLog) prints(calldepth int, message string, v []interface{}, err []error) error {
	if err == nil {
		err = this.scope
	}
	s := message
	if this.formatter != nil {
		s = this.formatter(message, v, err)
	}
	return this.filteredLogger(calldepth+1).Output(calldepth+this.soff+1, s)
}

func (this *sLog) Output(calldepth int, s string) error {
	return this.filteredLogger(calldepth).Output(calldepth+this.soff+1, s)
}

func (this *sLog) Printf(format string, v ...interface{}) {
	this.filteredLogger(2).Output(2, fmt.Sprintf(format, v...))
}

func (this *sLog) Print(v ...interface{}) {
	this.filteredLogger(2).Output(2, fmt.Sprint(v...))
}

func (this *sLog) Println(v ...interface{}) {
	this.filteredLogger(2).Output(2, fmt.Sprintln(v...))
}

var dscrd = log.New(ioutil.Discard, "", 0)
var drain = &sLog{formatter: nil, logger: dscrd, scope: nil}

type sSelector struct {
	*sLogger
	scope []error
}

func (this *sSelector) Log(pri Priority) Log {
	return this.logs[pri.Bound()].ScopedLog(this.scope...)
}

func (this *sSelector) Info() Log {
	return this.logs[PriorityInfo].ScopedLog(this.scope...)
}

func (this *sSelector) Notice() Log {
	return this.logs[PriorityNotice].ScopedLog(this.scope...)
}

func (this *sSelector) Warning() Log {
	return this.logs[PriorityWarn].ScopedLog(this.scope...)
}

func (this *sSelector) Error() Log {
	return this.logs[PriorityError].ScopedLog(this.scope...)
}

func (this *sSelector) Trace(detail int) Log {
	if PriorityTrace+Priority(detail-1) <= this.level {
		return this.logs[PriorityTrace].ScopedLog(this.scope...)
	} else {
		return drain
	}
}

func (this *sSelector) Prints(message string, v ...interface{}) {
	this.scopedLog().prints(2, message, v, nil)
}

func (this *sSelector) Fatals(message string, v ...interface{}) {
	this.scopedLog().prints(2, message, v, nil)
	if !this.isSuccess() {
		os.Exit(1)
	}
}

func (this *sSelector) Logger() *log.Logger {
	return this.scopedLog().Logger()
}

func (this *sSelector) isSuccess() bool {
	return this.scope == nil || len(this.scope) == 1 && (this.scope[0] == nil || this.scope[0] == errSuccess || this.scope[0] == errEllipsis)
}

func (this *sSelector) scopedLog() Log {
	if this.isSuccess() {
		return this.logs[PriorityNotice].ScopedLog(this.scope...)
	} else {
		return this.logs[PriorityError].ScopedLog(this.scope...)
	}
}
