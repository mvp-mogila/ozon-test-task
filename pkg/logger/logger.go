package logger

import (
	"context"
	"log"
	"os"
)

type loggerContextKey struct{}

var loggerCtxKey = loggerContextKey{}

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func AddToContext(ctx context.Context, l *log.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

func getFromContext(ctx context.Context) *log.Logger {
	l, ok := ctx.Value(loggerCtxKey).(*log.Logger)
	if !ok {
		return log.Default()
	}
	return l
}

func Infof(ctx context.Context, msg string, args ...any) {
	getFromContext(ctx).Printf("INFO:  "+msg+"\n", args...)
}

func Errorf(ctx context.Context, msg string, args ...any) {
	getFromContext(ctx).Printf("ERROR:  "+msg+"\n", args...)
}
