package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_ "go.uber.org/automaxprocs"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)
		os.Exit(1)
	}
}

func getLogger() (*zap.Logger, error) {
	logLevel := zapcore.InfoLevel
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		var err error
		logLevel, err = zapcore.ParseLevel(levelStr)
		if err != nil {
			return nil, err
		}
	}

	logConf := zap.NewProductionConfig()
	logConf.Level = zap.NewAtomicLevelAt(logLevel)

	logger, err := logConf.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func run() error {
	logger, err := getLogger()
	if err != nil {
		return err
	}

	defer logger.Sync()

	logger.Info("Hello world!", zap.String("location", "world"))

	return nil
}
