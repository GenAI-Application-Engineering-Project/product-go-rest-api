package logger

import (
	"io"
	"strings"
	"time"

	"product-services/internal/interfaces"

	"github.com/rs/zerolog"
)

type DefaultLogger struct {
	logger zerolog.Logger
}

func NewLogger(appEnv string, svc string, w io.Writer) interfaces.AppLogger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(w).With().
		Str("service", svc).
		Timestamp().
		CallerWithSkipFrameCount(3).
		Logger()

	// Set log level based on environment
	appEnv = strings.ToLower(appEnv)
	if appEnv == "production" || appEnv == "prod" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return &DefaultLogger{logger: logger}
}

func (l *DefaultLogger) Logger() zerolog.Logger {
	return l.logger
}

func (l *DefaultLogger) Fatal(err error, msg string) {
	l.logger.Fatal().Err(err).Msg(msg)
}
