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
		slog.With(err).Prints("Returned from emit():", "i", i)
	}
	slog.Info().Prints("Finished unconditional", "time", time.Now())
}
