package log

import (
	"io"
	"log/slog"
)

type Config struct {
	handler   slog.Handler
	writer    io.Writer
	addSource bool
	level     string
}

type Option interface {
	apply(*Config)
}

type OptionFunc func(*Config)

func (o OptionFunc) apply(c *Config) {
	o(c)
}

//nolint:revive // usable configuration function
func WithHandler(handler slog.Handler) OptionFunc {
	if handler == nil {
		return func(*Config) {}
	}

	return func(config *Config) {
		config.handler = handler
	}
}

//nolint:revive // usable configuration function
func WithWriter(writer io.Writer) OptionFunc {
	if writer == nil {
		return func(*Config) {}
	}

	return func(config *Config) {
		config.writer = writer
	}
}

func WithLevel(level string) OptionFunc {
	if level == "" {
		return func(*Config) {}
	}

	return func(config *Config) {
		config.level = level
	}
}

//nolint:revive // usable configuration function
func WithSource() OptionFunc {
	return func(config *Config) {
		config.addSource = true
	}
}
