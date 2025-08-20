package internal

import (
	"io"
	"log/slog"
	"os"
)

type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

type Options struct {
	Writer     io.Writer
	Level      slog.Leveler
	AddSource  bool
	TimeFormat string
	Color      ColorMode
}

func (o Options) WithDefaults() Options {
	if o.Writer == nil {
		o.Writer = os.Stdout
	}
	if o.Level == nil {
		o.Level = slog.LevelInfo
	}
	if o.TimeFormat == "" {
		o.TimeFormat = "15:04:05.000"
	}
	return o
}