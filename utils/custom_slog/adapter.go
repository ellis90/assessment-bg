package custom_slog

import (
	"context"
	"github.com/jackc/pgx/v4"
	"golang.org/x/exp/slog"
)

type Logger struct {
	l *slog.Logger
}

func NewLogger(l *slog.Logger) *Logger {
	return &Logger{l: l}
}

func (l *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]any) {
	logArrs := make([]slog.Attr, 0, len(data))

	for k, v := range data {
		logArrs = append(logArrs, slog.Any(k, v))
	}

	log := l.l.With(logArrs)

	switch level {
	case pgx.LogLevelTrace:
		log.Debug("PGX_LOG_LEVEL")
	case pgx.LogLevelDebug:
		log.Debug(msg)
	case pgx.LogLevelInfo:
		log.Info(msg)
	case pgx.LogLevelWarn:
		log.With(msg)
	case pgx.LogLevelError:
		log.Error(msg)
	default:
	}
}
