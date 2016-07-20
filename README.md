# slog
Structured logging for go.

[![Build Status (Linux)](https://travis-ci.org/baobabus/slog.svg?branch=master)](https://travis-ci.org/baobabus/slog)

# Usage

Example logging with global logger:
```go
package main

import (
	"errors"
	"github.com/baobabus/slog"
	"time"
)

func emit(i int) error {
	if i % 2 == 0 {
		return errors.New("Bad input")
	}
	return nil
}

func main () {
	slog.Info().Prints("Starting conditional", "time", time.Now())
	for i := 0; i < 4; i++ {
		err := emit(i)
		slog.On(err).Prints("Problem with emit():", "i", i)
	}
	slog.Info().Prints("Finished conditional", "time", time.Now())
	slog.Info().Prints("Starting unconditional", "time", time.Now())
	for i := 0; i < 4; i++ {
		err := emit(i)
		slog.With(err).Warning().Prints("Returned from emit():", "i", i)
	}
	slog.Info().Prints("Finished unconditional", "time", time.Now())
}
```

Output:
```
INFO 2016/07/20 11:00:53 Starting conditional time=2016-07-20T11:00:53-05:00
ERROR 2016/07/20 11:00:53 Problem with emit(): i=0 - error=Bad input
ERROR 2016/07/20 11:00:53 Problem with emit(): i=2 - error=Bad input
INFO 2016/07/20 11:00:53 Finished conditional time=2016-07-20T11:00:53-05:00
INFO 2016/07/20 11:00:53 Starting unconditional time=2016-07-20T11:00:53-05:00
WARNING 2016/07/20 11:00:53 Returned from emit(): i=0 - error=Bad input
WARNING 2016/07/20 11:00:53 Returned from emit(): i=1 - success
WARNING 2016/07/20 11:00:53 Returned from emit(): i=2 - error=Bad input
WARNING 2016/07/20 11:00:53 Returned from emit(): i=3 - success
INFO 2016/07/20 11:00:53 Finished unconditional time=2016-07-20T11:00:53-05:00
```
