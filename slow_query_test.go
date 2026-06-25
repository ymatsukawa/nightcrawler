package slow_query

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	d "github.com/ymatsukawa/slow_query/detector"
)

func TestNewSlogHandler(t *testing.T) {
	type Args struct {
		SourceHandler slog.Handler
		Suppress      []int
	}

	tests := []struct {
		Name string
		Args Args
	}{
		{
			Name: "Create new slog handler with suppress",
			Args: Args{
				SourceHandler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
				Suppress:      []int{},
			},
		},
		{
			Name: "Create new slog handler with suppress (with categories)",
			Args: Args{
				SourceHandler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
				Suppress:      []int{d.SelectMany},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			handler := NewSlogHandler(tt.Args.SourceHandler, tt.Args.Suppress)
			assert.NotNil(t, handler)
			assert.Equal(t, tt.Args.SourceHandler, handler.SourceHandler)
			assert.Equal(t, tt.Args.Suppress, handler.LogInfo.Suppress)
		})
	}
}
