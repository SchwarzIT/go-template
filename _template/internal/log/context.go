package log

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

// traceIDKey and spanIDKey define the key values for log attributes, for tracing metadata
//
// https://opentelemetry.io/docs/specs/otel/compatibility/logging_trace_context/
const (
	traceIDKey = "trace_id"
	spanIDKey  = "span_id"
)

const maxSpanContextAttrs = 2

// SpanContextHandler is a slog.Handler wrapper that adds trace data as log attributes on each
// Handle call, given that the input context to the method contains a valid trace.SpanContext.
type SpanContextHandler struct {
	withSpanID bool
	handler    slog.Handler
}

//nolint:gochecknoglobals // instantiated once on a package level to be used in comparisons in NewSpanContextHandler
var nilHandler = slog.Handler(nil)

// NewSpanContextHandler creates a SpanContextHandler from the input slog.Handler and optional
// configuration.
func NewSpanContextHandler(handler slog.Handler, withSpanID bool) slog.Handler {
	if handler == nil || handler == nilHandler {
		handler = defaultHandler()
	}

	return &SpanContextHandler{
		withSpanID: withSpanID,
		handler:    handler,
	}
}

func defaultHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
// It is called early, before any arguments are processed,
// to save effort if the log event should be discarded.
// If called from a Logger method, the first argument is the context
// passed to that method, or context.Background() if nil was passed
// or the method does not take a context.
// The context is passed so Enabled can use its values
// to make a decision.
func (h *SpanContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle handles the Record.
// It will only be called when Enabled returns true.
// The Context argument is as for Enabled.
// It is present solely to provide Handlers access to the context's values.
// Canceling the context should not affect record processing.
// (Among other things, log messages may be necessary to debug a
// cancellation-related problem.)
//
// Handle methods that produce output should observe the following rules:
//   - If r.Time is the zero time, ignore the time.
//   - If r.PC is zero, ignore it.
//   - Attr's values should be resolved.
//   - If an Attr's key and value are both the zero value, ignore the Attr.
//     This can be tested with attr.Equal(Attr{}).
//   - If a group's key is empty, inline the group's Attrs.
//   - If a group has no Attrs (even if it has a non-empty key),
//     ignore it.
func (h *SpanContextHandler) Handle(ctx context.Context, record slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		attrs := make([]slog.Attr, 1, maxSpanContextAttrs)
		attrs[0] = slog.String(traceIDKey, sc.TraceID().String())

		if h.withSpanID {
			attrs = append(attrs, slog.String(spanIDKey, sc.SpanID().String()))
		}

		record.AddAttrs(attrs...)
	}

	return h.handler.Handle(ctx, record)
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (h *SpanContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SpanContextHandler{
		withSpanID: h.withSpanID,
		handler:    h.handler.WithAttrs(attrs),
	}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// The keys of all subsequent attributes, whether added by With or in a
// Record, should be qualified by the sequence of group names.
//
// How this qualification happens is up to the Handler, so long as
// this Handler's attribute keys differ from those of another Handler
// with a different sequence of group names.
//
// A Handler should treat WithGroup as starting a Group of Attrs that ends
// at the end of the log event. That is,
//
//	logger.WithGroup("s").LogAttrs(level, msg, slog.Int("a", 1), slog.Int("b", 2))
//
// should behave like
//
//	logger.LogAttrs(level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
//
// If the name is empty, WithGroup returns the receiver.
func (h *SpanContextHandler) WithGroup(name string) slog.Handler {
	return &SpanContextHandler{
		withSpanID: h.withSpanID,
		handler:    h.handler.WithGroup(name),
	}
}
