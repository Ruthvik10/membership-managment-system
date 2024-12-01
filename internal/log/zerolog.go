package log

import (
	"io"

	"github.com/rs/zerolog"
)

type ZLogger struct {
	zerolog.Logger
}

func NewZLogger(out io.Writer) *ZLogger {
	return &ZLogger{Logger: zerolog.New(out).With().Timestamp().Logger()}
}

func (l *ZLogger) WriteInfo(msg string, fields map[string]interface{}) {
	l.Info().Fields(fields).Msg(msg)
}

func (l *ZLogger) WriteError(msg string, err error, fields map[string]interface{}) {
	l.Error().Fields(fields).Err(err).Msg(msg)
}

func (l *ZLogger) WriteFatal(msg string, err error, fields map[string]interface{}) {
	l.Fatal().Fields(fields).Err(err).Msg(msg)
}
