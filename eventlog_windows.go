// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"errors"
	"sys/windows"
)

func (this Priority) EventlogPriority() uint16 {
	switch this.Bound() {
	case PriorityError:
		return windows.EVENTLOG_ERROR_TYPE
	case PriorityWarn:
		return windows.EVENTLOG_WARNING_TYPE
	case PriorityNotice:
		return windows.EVENTLOG_INFORMATION_TYPE
	case PriorityInfo:
		return windows.EVENTLOG_INFORMATION_TYPE
	case PriorityTrace:
		return windows.EVENTLOG_INFORMATION_TYPE
	}
}

// TODO Add implementation

func NewEventlogFacility(level Priority) (Facility, error) {
	return nil, errors.New("not implemented")
}

func init() {
	newSyslogFacility = NewEventlogFacility
}
