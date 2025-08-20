# sweet-logger

Easy pretty-handler for `log/slog`, focused on local development:
- Compact lines: `[time] LEVEL msg key=val ...`
- Colors by level (debug/trace muted, info green, warn yellow, error/critical red)
- Support for `With`/`WithGroup`, `AddSource`
- No external dependencies

  Easy pretty-handler для `log/slog`, ориентирован на локальную разработку:
- Компактные строки: `[time] LEVEL msg key=val ...`
- Цвета по уровню (debug/trace приглушённые, info зелёный, warn жёлтый, error/critical красный)
- Поддержка `With`/`WithGroup`, `AddSource`
- Без внешних зависимостей
  
## Install | Установка

```bash
go get github.com/h4tecancel/sweet-logger
