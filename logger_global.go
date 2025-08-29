package logger

import (
	"context"
	"fmt"
	"net/http"
)

var loggerInstance *Logger

func InitGlobalLogger(config LoggerConfig) error {
	gLogger, err := NewLogger(config)
	if err != nil {
		return err
	}

	loggerInstance = gLogger
	return nil
}

func DestroyGlobalLogger() {
	if loggerInstance != nil {
		loggerInstance.Flush()
		loggerInstance = nil
	}
}

func GetGlobalLogger() *Logger {
	return loggerInstance
}

func Debug(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Debugw(msg, combinedAttributes...)
}

func Info(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Infow(msg, combinedAttributes...)
}

func Warn(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Warnw(msg, combinedAttributes...)
}

func Error(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Errorw(msg, combinedAttributes...)
}

func Panic(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Panicw(msg, combinedAttributes...)
}

func Fatal(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Fatalw(msg, combinedAttributes...)
}

func Debugf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Debugw(msg, combinedAttributes...)
}

func Infof(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Infow(msg, combinedAttributes...)
}

func Warnf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Warnw(msg, combinedAttributes...)
}

func Errorf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Errorw(msg, combinedAttributes...)
}

func Panicf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Panicw(msg, combinedAttributes...)
}

func Fatalf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := loggerInstance.combineAttributes(ctx)
	loggerInstance.logger.Fatalw(msg, combinedAttributes...)
}

func Debugw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := loggerInstance.combineAttributes(ctx, keysAndValues...)
	loggerInstance.logger.Debugw(msg, combinedAttributes...)
}

func Infow(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := loggerInstance.combineAttributes(ctx, keysAndValues...)
	loggerInstance.logger.Infow(msg, combinedAttributes...)
}

func Warnw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := loggerInstance.combineAttributes(ctx, keysAndValues...)
	loggerInstance.logger.Warnw(msg, combinedAttributes...)
}

func Errorw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := loggerInstance.combineAttributes(ctx, keysAndValues...)
	loggerInstance.logger.Errorw(msg, combinedAttributes...)
}

func Panicw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := loggerInstance.combineAttributes(ctx, keysAndValues...)
	loggerInstance.logger.Panicw(msg, combinedAttributes...)
}

func Fatalw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := loggerInstance.combineAttributes(ctx, keysAndValues...)
	loggerInstance.logger.Fatalw(msg, combinedAttributes...)
}

func Flush() {
	loggerInstance.Flush()
}

func IsDevMode() bool {
	return loggerInstance.IsDevMode()
}

func GenerateRequestID() string {
	return loggerInstance.GenerateRequestID()
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return loggerInstance.SetRequestID(ctx, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
	return loggerInstance.GetRequestID(ctx)
}

func SetUser(ctx context.Context, user any) context.Context {
	return loggerInstance.SetUser(ctx, user)
}

func GetUser(ctx context.Context) (any, bool) {
	return loggerInstance.GetUser(ctx)
}

func SetUserIP(ctx context.Context, userIP string) context.Context {
	return loggerInstance.SetUserIP(ctx, userIP)
}

func GetUserIP(ctx context.Context) (string, bool) {
	return loggerInstance.GetUserIP(ctx)
}

func GetExtraFields(ctx context.Context) (map[string]any, bool) {
	return loggerInstance.GetExtraFields(ctx)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return loggerInstance.LoggerMiddleware(next)
}
