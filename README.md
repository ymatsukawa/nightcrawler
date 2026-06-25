# Nightcrawler - Slow query detect handler

## Description

Detect slow query doubted log from slog.

## Installation

```
go get -u github.com/ymatsukawa/nightcrawler
```

## Usage

```go
import (
  "os"
  "log/slog"

  "github.com/ymatsukawa/nightcrawler"
  d "github.com/ymatsukawa/nightcrawler/detector"
)

...

baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
// or
// baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

logger := slog.New(nightcrawler.NewSlogHandler(baseHandler, nil))
// or
// suppress := []int{d.SelectMany}
// logger := slog.New(nightcrawler.NewSlogHandler(baseHandler, suppress))
```

## Show case

**json**
```json
{
  "time": "2026-01-01:00:00.00000000Z",
  "level": "INFO",
  "msg": "SELECT id, name FROM \"users\" WHERE name LIKE '%a' LIMIT 10",
  "slow_query": "LIKE %..."
}
```

**text**
```txt
time=2026-01-01T12:00:00.000Z level=INFO msg="SELECT id, name FROM \"users\" WHERE name LIKE '%a' LIMIT 10" slow_query="LIKE %..."
```

No attribution of "slow_query" when no detection about slow query or supressing the target.

## LICENSE
MIT
