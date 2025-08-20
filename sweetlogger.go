package sweetlogger

import (
	"log/slog"

	"github.com/h4tecancel/sweet-logger/internal"
)

type (
	Options   = internal.Options
	ColorMode = internal.ColorMode
)

const (
	ColorAuto   = internal.ColorAuto
	ColorAlways = internal.ColorAlways
	ColorNever  = internal.ColorNever
)

// New строит *slog.Logger с нашим PrettyHandler'ом.
func New(opts Options) *slog.Logger {
	h := internal.NewPrettyHandler(opts)
	return slog.New(h)
}
