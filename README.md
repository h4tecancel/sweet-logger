# sweet-logger

Лёгкий цветной pretty-handler для `log/slog`, ориентирован на локальную разработку:
- Компактные строки: `[time] LEVEL msg key=val ...`
- Цвета по уровню (debug/trace приглушённые, info зелёный, warn жёлтый, error/critical красный)
- Поддержка `With`/`WithGroup`, `AddSource`
- Без внешних зависимостей

## Установка

```bash
go get github.com/h4tecancel/sweet-logger