package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestNewSpanContextHandler(t *testing.T) {
	jsonHandler := slog.NewJSONHandler(io.Discard, nil)

	for _, testcase := range []struct {
		name       string
		handler    slog.Handler
		withSpanID bool
		wants      slog.Handler
	}{
		{
			name:       "WithHandlerAndSpanID",
			handler:    jsonHandler,
			withSpanID: true,
			wants: &SpanContextHandler{
				withSpanID: true,
				handler:    jsonHandler,
			},
		},
		{
			name:    "WithNilHandler",
			handler: nil,
			wants: &SpanContextHandler{
				handler: defaultHandler(),
			},
		},
		{
			name:    "WithNilInterfaceHandler",
			handler: slog.Handler(nil),
			wants: &SpanContextHandler{
				handler: defaultHandler(),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			handler := NewSpanContextHandler(testcase.handler, testcase.withSpanID)
			require.Equal(t, testcase.wants, handler)
		})
	}
}

func TestSpanContextHandler_Handle(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewSpanContextHandler(slog.NewJSONHandler(buf, nil), true))
	testTraceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	testSpanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}

	for _, testcase := range []struct {
		name            string
		message         string
		withSpanContext bool
	}{
		{
			name:    "WithoutTrace",
			message: "test log event",
		},
		{
			name:            "WithTrace",
			message:         "test log event",
			withSpanContext: true,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			buf.Reset()

			ctx := context.Background()
			if testcase.withSpanContext {
				ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: testTraceID,
					SpanID:  testSpanID,
				}))
			}

			logger.InfoContext(ctx, testcase.message)

			var m map[string]any

			require.NoError(t, json.Unmarshal(buf.Bytes(), &m))

			if traceID, ok := m[traceIDKey]; ok && testcase.withSpanContext {
				require.Equal(t, testTraceID.String(), traceID)
			}

			if spanID, ok := m[spanIDKey]; ok && testcase.withSpanContext {
				require.Equal(t, testSpanID.String(), spanID)
			}
		})
	}
}

func TestSpanContextHandler_WithAttrs(t *testing.T) {
	handler1 := NewSpanContextHandler(slog.NewJSONHandler(io.Discard, nil), false)
	handler2 := NewSpanContextHandler(slog.NewJSONHandler(io.Discard, nil), false)
	handler2wants := NewSpanContextHandler(
		slog.NewJSONHandler(io.Discard, nil).
			WithAttrs([]slog.Attr{slog.String("key", "value"), slog.Int("num", 0)}),
		false,
	)

	for _, testcase := range []struct {
		name    string
		handler slog.Handler
		attrs   []slog.Attr
		wants   slog.Handler
	}{
		{
			name:    "NoAttrs",
			handler: handler1,
			wants:   handler1,
		},
		{
			name:    "AddAttrs",
			handler: handler2,
			attrs:   []slog.Attr{slog.String("key", "value"), slog.Int("num", 0)},
			wants:   handler2wants,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			require.Equal(t,
				testcase.wants,
				testcase.handler.WithAttrs(testcase.attrs),
			)
		})
	}
}

func TestSpanContextHandler_WithGroup(t *testing.T) {
	handler := NewSpanContextHandler(slog.NewJSONHandler(io.Discard, nil), false)
	handlerWants := NewSpanContextHandler(
		slog.NewJSONHandler(io.Discard, nil).WithGroup("test"),
		false,
	)

	for _, testcase := range []struct {
		name    string
		handler slog.Handler
		group   string
		wants   slog.Handler
	}{
		{
			name:    "AddGroup",
			handler: handler,
			group:   "test",
			wants:   handlerWants,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			require.Equal(t,
				testcase.wants,
				testcase.handler.WithGroup(testcase.group),
			)
		})
	}
}
