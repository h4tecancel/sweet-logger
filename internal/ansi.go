package internal

import (
	"os"
)

const (
	reset = "\x1b[0m"

	fgGray   = "\x1b[90m"
	fgGreen  = "\x1b[32m"
	fgYellow = "\x1b[33m"
	fgRed    = "\x1b[31m"
	fgBlue   = "\x1b[34m"
)

func isTTY(f *os.File) bool {
	// Простая эвристика: для Auto будем считать, что если это stdout/stderr — красим.
	// Если захочешь строгую проверку TTY — можно позже добавить x/sys/term.
	return f == os.Stdout || f == os.Stderr
}
