package log

import (
	"context"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Fields = log.Fields

var traceIDKey int

func WithTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, &traceIDKey, uuid.New().String())
}

func TraceID(ctx context.Context) string {
	id, ok := ctx.Value(&traceIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

func Fatal(err error) {
	log.WithError(err).Fatal()
}

func Init(ctx context.Context) *log.Entry {
	id := TraceID(ctx)
	if id == "" {
		id = "<none>"
	}
	return log.WithField("traceID", id)
}

func WithError(ctx context.Context, err error) *log.Entry {
	return Init(ctx).WithError(err)
}

func WithField(ctx context.Context, key string, value interface{}) *log.Entry {
	return Init(ctx).WithField(key, value)
}

func WithFields(ctx context.Context, fields log.Fields) *log.Entry {
	return Init(ctx).WithFields(fields)
}

func Info(ctx context.Context, args ...interface{}) {
	Init(ctx).Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	Init(ctx).Infof(format, args...)
}
