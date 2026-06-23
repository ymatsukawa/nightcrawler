# Slow query detect handler

## Description

Detect slow query from slog.

## Installation

```
go get -u github.com/ymatsukawa/slow_query
```

## Usage

```go
import (
  "os"
  "log/slog"

  "github.com/ymatsukawa/slow_query"
  c "github.com/ymatsukawa/slow_query/category"
)

...

baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
// or
// baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

logger := slog.New(slow_query.NewSlogHandler(baseHandler, nil))
// or
// suppress := []string{c.SelectMany}
// logger := slog.New(slow_query.NewSlogHandler(baseHandler, suppress))
```

## Show case

**json**
```json
{
  "time": "2026-01-01:00:00.00000000Z",
  "level": "INFO",
  "msg": "SELECT * FROM \"users\"",
  "slow_query": "Select many(select * OR no limit)"
}
```

**text**
```txt
time=2026-01-01T12:00:00.000Z level=INFO msg="SELECT * FROM \"users\"" slow_query="Select many(select * OR no limit)"
```

No attribution of "slow_query" when no detection about slow query.

## LICENSE
MIT
