package nightcrawler

import (
	"context"
	"log/slog"

	"github.com/ymatsukawa/nightcrawler/detector"
)

type slowQueryDetectHandler struct {
	SourceHandler slog.Handler
	LogInfo       detector.ParseInfo
}

func NewSlogHandler(sourceHandler slog.Handler, suppress []int) *slowQueryDetectHandler {
	return &slowQueryDetectHandler{
		SourceHandler: sourceHandler,
		LogInfo: detector.ParseInfo{
			PreviousLine: "",
			Suppress:     suppress,
		},
	}
}

func (h *slowQueryDetectHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.SourceHandler.Enabled(ctx, level)
}

func (h *slowQueryDetectHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slowQueryDetectHandler{
		SourceHandler: h.SourceHandler.WithAttrs(attrs),
		LogInfo:       h.LogInfo,
	}
}

func (h *slowQueryDetectHandler) WithGroup(name string) slog.Handler {
	return &slowQueryDetectHandler{
		SourceHandler: h.SourceHandler.WithGroup(name),
		LogInfo:       h.LogInfo,
	}
}

func (h *slowQueryDetectHandler) Handle(ctx context.Context, record slog.Record) error {
	prev := record.Message
	if class, ok := detector.CatchSlowQuery(record.Message, h.LogInfo); ok {
		record.AddAttrs(slog.String("slow_query", class))
	}

	err := h.SourceHandler.Handle(ctx, record)
	h.LogInfo.PreviousLine = prev

	return err
}
