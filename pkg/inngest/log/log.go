package log

import (
	"context"
	"os"

	"github.com/inngest/inngest/pkg/cli"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	DefaultLevel = zerolog.InfoLevel
)

var (
	prettyFormatter = zerolog.ConsoleWriter{
		Out:         os.Stderr,
		FormatLevel: func(i any) string { return "" },
	}
)

type loggerKey struct{}

func With(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func From(ctx context.Context) *zerolog.Logger {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		return Default()
	}

	l := logger.(zerolog.Logger)
	return &l
}

func New(lvl zerolog.Level) *zerolog.Logger {
	l := zerolog.New(os.Stderr).Level(lvl).With().Timestamp().Logger()
	if cli.IsTtyLogger() {
		l = l.Output(prettyFormatter)
	}
	return &l
}

func Copy(l zerolog.Logger) zerolog.Logger {
	c := l.Output(os.Stderr)
	if cli.IsTtyLogger() {
		c = c.Output(prettyFormatter)
	}
	return c
}

func Default() *zerolog.Logger {
	lvl, err := zerolog.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	return New(lvl)
}
