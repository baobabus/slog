// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"log"
	"log/syslog"
)

func (this Priority) SyslogPriority() syslog.Priority {
	return syslog.Priority(this.Bound()+3) | syslog.LOG_USER
}

type fSyslog struct {
	writer *syslog.Writer
}

func NewSyslogFacility(level Priority) (Facility, error) {
	writer, err := syslog.New(level.SyslogPriority(), "")
	if err != nil {
		return nil, err
	}
	return &fSyslog{writer: writer}, nil
}

func (this *fSyslog) OpenLogs(level Priority) (map[Priority]*log.Logger, error) {
	res := make(map[Priority]*log.Logger, prioritiesCount)
	for pri := PriorityError; pri <= PriorityTrace; pri++ {
		if pri <= level {
			// syslog does it's own time and priority stamping,
			// although the priority is in numeric form, so we'll keep ours
			if pri < PriorityTrace {
				res[pri] = log.New(&lSyslog{writer: this.writer, priority: pri}, pri.Tag(), 0)
			} else {
				res[pri] = log.New(&lSyslog{writer: this.writer, priority: pri}, pri.Tag(), log.Lshortfile)
			}
		} else {
			res[pri] = drain.Logger()
		}
	}
	return res, nil
}

type lSyslog struct {
	writer   *syslog.Writer
	priority Priority
}

func (this lSyslog) Write(p []byte) (n int, err error) {
	err = nil
	switch this.priority {
	case PriorityError:
		err = this.writer.Err(string(p))
	case PriorityWarn:
		err = this.writer.Warning(string(p))
	case PriorityNotice:
		err = this.writer.Notice(string(p))
	case PriorityInfo:
		err = this.writer.Info(string(p))
	case PriorityTrace:
		err = this.writer.Debug(string(p))
	}
	if err != nil {
		n = 0
	} else {
		n = len(p)
	}
	return n, err
}

func init() {
	newSyslogFacility = NewSyslogFacility
}
