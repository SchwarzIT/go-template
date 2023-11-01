package log

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/stretchr/testify/require"
)

func TestInterceptorLogger(t *testing.T) {
	testcases := []struct {
		Name string

		LoggerLevel slog.Leveler

		Level   logging.Level
		Message string
		Fields  []any

		Expected string
	}{{
		Name:        "Info log with fields",
		LoggerLevel: slog.LevelInfo,
		Level:       logging.LevelInfo,
		Message:     "Test info message",
		Fields: []any{
			"string", "hello world",
			"int", 42,
			"int16", int16(1),
			"int32", int32(2),
			"int64", int64(3),
			"bool", true,
			"float32", float32(3.141),
			"float64", 3.1415926,
			"any",
			struct{}{},
		},
		Expected: `","level":"INFO","msg":"Test info message","string":"hello world",` +
			`"int":42,"int16":1,"int32":2,"int64":3,"bool":true,"float32":3.1410000324249268,"float64":3.1415926,"any":{}}`,
	}, {
		Name:        "Debug log but deactivated",
		LoggerLevel: slog.LevelInfo,
		Level:       logging.LevelDebug,
		Message:     "debug log deactivated",
		Fields:      nil,
		Expected:    "",
	}, {
		Name:        "Warn log",
		LoggerLevel: slog.LevelInfo,
		Level:       logging.LevelWarn,
		Message:     "This is a warning log",
		Fields:      nil,
		Expected:    `,"level":"WARN","msg":"This is a warning log"}`,
	}, {
		Name:        "Error log",
		LoggerLevel: slog.LevelInfo,
		Level:       logging.LevelError,
		Message:     "This is an error log",
		Fields:      nil,
		Expected:    `,"level":"ERROR","msg":"This is an error log"}`,
	}, {
		Name:        "Invalid log level",
		LoggerLevel: slog.LevelInfo,
		Level:       42,
		Message:     "level 42!",
		Fields:      nil,
		Expected:    `,"level":"ERROR+34","msg":"level 42!"}`,
	}, {
		Name:        "Invalid field key",
		LoggerLevel: slog.LevelInfo,
		Level:       logging.LevelInfo,
		Message:     "test",
		Fields: logging.Fields{
			42, "test",
		},
		Expected: `,"level":"INFO","msg":"test","!BADKEY":42,"!BADKEY":"test"}`,
	}}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			logger := slog.New(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{
				AddSource: false, // so it can be tested without a filesystem path
				Level:     testcase.LoggerLevel,
			}))

			loggerFunc := InterceptorLogger(logger)

			loggerFunc.Log(context.Background(), testcase.Level, testcase.Message, testcase.Fields...)

			require.Contains(t, buffer.String(), testcase.Expected)
		})
	}
}
