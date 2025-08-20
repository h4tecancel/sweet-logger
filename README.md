# sweet-logger

Easy pretty-handler for `log/slog`, focused on local development:
- Compact lines: `[time] LEVEL msg key=val ...`
- Colors by level (debug/trace muted, info green, warn yellow, error/critical red)
- Support for `With`/`WithGroup`, `AddSource`
- No external dependencies
  
## Install | Установка

```bash
go get github.com/h4tecancel/sweet-logger
