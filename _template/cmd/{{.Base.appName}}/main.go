package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	"go.uber.org/automaxprocs/maxprocs"

	"{{.Base.moduleName}}/internal/log"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)
		os.Exit(1)
	}
}

//nolint:unparam // currently run returns no error; remove this exclusion directive if otherwise
func run() error {
	logger := log.New(log.WithLevel(os.Getenv("LOG_LEVEL")))
	ctx := context.Background()

	_, err := maxprocs.Set(maxprocs.Logger(func(s string, i ...interface{}) {
		logger.DebugContext(ctx, fmt.Sprintf(s, i...))
	}))
	if err != nil {
		return fmt.Errorf("setting max procs: %w", err)
	}

	logger.InfoContext(ctx, "Hello world!", slog.String("location", "world"))

	return nil
}
