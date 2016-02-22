// Copyright 2016 Aleksey Blinov. All rights reserved.
// +build !slognoflags

package slog

import (
	"flag"
)

func init() {
	flag.StringVar(&rtLevel, "loglevel", "info", "set logging level; supported values are \"error\", \"warn\", \"notice\" and \"info\"")
	flag.UintVar(&rtTrace, "trace", 0, "enable trace logging with specified verbosity")
	flag.StringVar(&rtFormat, "logfmt", "simple", "set logging format; supported values are \"simple\", \"json\" and \"json-pretty\"")
	flag.StringVar(&rtLog, "log", "stderr", "send log output to specified file; recognized special values are \"stdout\", \"stderr\" and \"syslog\"")
}
