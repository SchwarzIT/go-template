package log

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	for _, testcase := range []struct {
		name  string
		level string
		wants *slog.Logger
	}{
		{
			name:  "Valid/LowercaseDebug",
			level: "debug",
			wants: slog.New(&SpanContextHandler{
				withSpanID: true,
				handler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
					AddSource: false,
					Level:     slog.LevelDebug,
				}),
			}),
		},
		{
			name:  "Valid/UppercaseError",
			level: "ERROR",
			wants: slog.New(&SpanContextHandler{
				withSpanID: true,
				handler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
					AddSource: false,
					Level:     slog.LevelError,
				}),
			}),
		},
		{
			name:  "Invalid/UnknownLevelReturnsDefaults",
			level: "SERVER_ON_FIRE",
			wants: slog.New(&SpanContextHandler{
				withSpanID: true,
				handler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
					AddSource: false,
					Level:     slog.LevelInfo,
				}),
			}),
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			require.Equal(t, testcase.wants, New(WithLevel(testcase.level)))
		})
	}
}

func TestNoOp(t *testing.T) {
	require.Equal(t, slog.New(noOpHandler{}), NoOp())
}

func TestNoOpHandler_Enabled(t *testing.T) {
	require.False(t, noOpHandler{}.Enabled(context.Background(), slog.LevelDebug))
}

func TestNoOpHandler_Handle(t *testing.T) {
	require.Nil(t, noOpHandler{}.Handle(context.Background(), slog.Record{}))
}

func TestNoOpHandler_WithAttrs(t *testing.T) {
	require.Equal(t, noOpHandler{}, noOpHandler{}.WithAttrs([]slog.Attr{slog.String("nope", "zero")}))
}

func TestNoOpHandler_WithGroup(t *testing.T) {
	require.Equal(t, noOpHandler{}, noOpHandler{}.WithGroup("nope"))
}
