package zerologger

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/redselig/currencier/internal/domain/usecase"
	"github.com/redselig/currencier/internal/util"
)

var _ usecase.Logger = (*Logger)(nil)
var sb strings.Builder

type Logger struct {
	logger  *zerolog.Logger
	isDebug bool
}

func NewLogger(logWriter io.Writer, isDebug bool) *Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(logWriter).With().Timestamp().Logger()
	return &Logger{logger: &logger, isDebug: isDebug}
}

func (l *Logger) log(ctx context.Context, message string, args ...interface{}) {
	l.logger.Log().Str("Request id", util.GetRequestID(ctx)).Msgf(message, args...)
}

func (l *Logger) Log(ctx context.Context, message interface{}, args ...interface{}) {
	switch mess := message.(type) {
	case error:
		if l.isDebug {
			err, ok := errors.Cause(mess).(stackTracer)
			if ok {
				st := err.StackTrace()
				l.log(ctx, fmt.Sprintf("%+v", st), args...)
				return
			}
		}
		l.log(ctx, mess.Error(), args...)
	case string:
		l.log(ctx, mess, args...)
	default:
		l.log(ctx, fmt.Sprintf("debug message %v has unknown type %v", message, mess), args...)
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
