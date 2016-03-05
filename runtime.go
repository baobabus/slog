// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"os"
	"strings"
	"sync"
)

var (
	rtLevel        = "info"
	rtTrace   uint = 0
	rtModules      = ""
	rtFormat       = "simple"
	rtLog          = "stderr"
)

var newSyslogFacility func(Priority) (Facility, error)

var (
	sharedFacilityMu = &sync.Mutex{}
	sharedFacility   Facility
	sharedLoggerMu   = &sync.Mutex{}
	sharedLogger     Logger
)

func SharedLogger() Logger {
	sharedLoggerMu.Lock()
	defer sharedLoggerMu.Unlock()
	if sharedLogger == nil {
		sharedLogger, _ = New(SharedFacility(), DefaultLevel(), DefaultFormatter(), DefaultFilter())
	}
	return sharedLogger
}

func SharedFacility() Facility {
	sharedFacilityMu.Lock()
	defer sharedFacilityMu.Unlock()
	if sharedFacility == nil {
		switch rtLog {
		case "stdout":
			sharedFacility, _ = NewStdFacility(os.Stdout)
		case "stderr":
			sharedFacility, _ = NewStdFacility(os.Stderr)
		case "syslog":
			if newSyslogFacility != nil {
				sharedFacility, _ = newSyslogFacility(DefaultLevel())
			}
		default:
			sharedFacility, _ = NewFileFacility(rtLog)
		}
	}
	return sharedFacility
}

func DefaultLevel() Priority {
	if rtTrace > 0 {
		return Priority(PriorityTrace + Priority(rtTrace-1))
	}
	switch rtLevel {
	case "error":
		return PriorityError
	case "warn":
		return PriorityWarn
	case "notice":
		return PriorityNotice
	case "info", "":
		return PriorityInfo
	default:
		return PriorityInfo
	}
}

func DefaultFilter() []string {
	res := make([]string, 0)
	if len(rtModules) > 0 {
		res = strings.Split(rtModules, ",")
		for i, _ := range res {
			res[i] = strings.TrimSpace(res[i])
		}
	}
	return res
}

func DefaultFormatter() Formatter {
	switch rtFormat {
	case "simple":
		return SimpleFormatter
	case "json":
		return CompactJsonFormatter
	case "json-pretty":
		return PrettyJsonFormatter
	default:
		return SimpleFormatter
	}
}
