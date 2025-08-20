// internal/pretty_handler.go
package internal

import (
	"bytes"
	"context" // << добавили
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PrettyHandler struct {
	opts   Options
	mu     sync.Mutex
	attrs  []slog.Attr
	groups []string
}

func NewPrettyHandler(opts Options) *PrettyHandler {
	o := opts.WithDefaults()
	return &PrettyHandler{opts: o}
}

// ДОЛЖНО быть именно context.Context в сигнатуре
func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	min := slog.LevelInfo
	if h.opts.Level != nil {
		min = h.opts.Level.Level()
	}
	return level >= min
}

// И тут — context.Context
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var buf bytes.Buffer

	// Цветовой режим
	useColor := false
	switch h.opts.Color {
	case ColorAlways:
		useColor = true
	case ColorNever:
		useColor = false
	default:
		if f, ok := h.opts.Writer.(*os.File); ok {
			useColor = isTTY(f)
		}
	}

	// [time]
	ts := time.Now()
	if !r.Time.IsZero() {
		ts = r.Time
	}
	buf.WriteString("[")
	buf.WriteString(ts.Format(h.opts.TimeFormat))
	buf.WriteString("] ")

	// LEVEL
	levelText, levelColor := levelStyle(r.Level)
	if useColor && levelColor != "" {
		buf.WriteString(levelColor)
	}
	buf.WriteString(levelText)
	if useColor {
		buf.WriteString(reset)
	}
	buf.WriteString(" ")

	// msg
	if r.Message != "" {
		buf.WriteString(r.Message)
	}

	// source (короткий путь)
	if h.opts.AddSource {
		if framePC := r.PC; framePC != 0 {
			if fn := runtime.FuncForPC(framePC); fn != nil {
				file, line := fn.FileLine(framePC)
				short := filepath.Base(file)
				buf.WriteString(" ")
				if useColor {
					buf.WriteString(fgBlue)
				}
				buf.WriteString("@")
				buf.WriteString(short)
				buf.WriteString(":")
				buf.WriteString(strconv.Itoa(line))
				if useColor {
					buf.WriteString(reset)
				}
			}
		}
	}

	// общие attrs из With
	allAttrs := make([]slog.Attr, 0, len(h.attrs)+8)
	allAttrs = append(allAttrs, h.attrs...)

	// атрибуты записи
	r.Attrs(func(a slog.Attr) bool {
		allAttrs = append(allAttrs, a)
		return true
	})

	// печать key=val
	for _, a := range allAttrs {
		key := h.renderKey(a.Key)
		val := h.renderValue(a.Value)
		buf.WriteString(" ")
		if useColor {
			buf.WriteString(fgGray)
		}
		buf.WriteString(key)
		buf.WriteString("=")
		if useColor {
			buf.WriteString(reset)
		}
		buf.WriteString(val)
	}

	buf.WriteByte('\n')

	_, err := io.Copy(h.opts.Writer, &buf)
	return err
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// копируем существующие attrs + добавляем новые
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)

	// НЕ копируем mutex; создаём новый PrettyHandler
	return &PrettyHandler{
		opts:   h.opts,
		attrs:  newAttrs,
		groups: append([]string{}, h.groups...),
	}
}


func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &PrettyHandler{
		opts:   h.opts,
		attrs:  append([]slog.Attr{}, h.attrs...),
		groups: append(append([]string{}, h.groups...), name),
	}
}

func levelStyle(lvl slog.Level) (text, color string) {
	switch {
	case lvl <= slog.LevelDebug:
		return "DEBUG", fgGray
	case lvl < slog.LevelWarn:
		return "INFO", fgGreen
	case lvl < slog.LevelError:
		return "WARN", fgYellow
	default:
		return "ERROR", fgRed
	}
}

func (h *PrettyHandler) renderKey(k string) string {
	if len(h.groups) == 0 {
		return k
	}
	return strings.Join(append(append([]string{}, h.groups...), k), ".")
}

func (h *PrettyHandler) renderValue(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		s := v.String()
		if strings.IndexFunc(s, isSpace) >= 0 {
			return strconv.Quote(s)
		}
		return s
	case slog.KindInt64:
		return strconv.FormatInt(v.Int64(), 10)
	case slog.KindUint64:
		return strconv.FormatUint(v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.FormatFloat(v.Float64(), 'f', -1, 64)
	case slog.KindBool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format(h.opts.TimeFormat)
	case slog.KindGroup:
		var b strings.Builder
		b.WriteByte('{')
		needSpace := false
		for _, a := range v.Group() {
			if needSpace {
				b.WriteByte(' ')
			}
			b.WriteString(a.Key)
			b.WriteByte('=')
			b.WriteString(h.renderValue(a.Value))
			needSpace = true
		}
		b.WriteByte('}')
		return b.String()
	default:
		return fmt.Sprintf("%v", v.Any())
	}
}

func isSpace(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }
