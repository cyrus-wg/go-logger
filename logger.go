package logger

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const (
	requestIdKey contextKey = "request_id"
	userKey      contextKey = "user"
)

const (
	requestIdContextKey = string(requestIdKey)
	userContextKey      = string(userKey)
)

type LoggerConfig struct {
	Development     bool
	RequestIDPrefix string
	FixedKeyValues  map[string]any
	ExtraFields     []string
}

type Logger struct {
	logger          *zap.SugaredLogger
	requestIDPrefix string
	fixedKeyValues  map[string]any
	extraFields     []string
	devMode         bool
}

func NewLogger(config LoggerConfig) (*Logger, error) {
	logger := &Logger{
		requestIDPrefix: config.RequestIDPrefix,
		extraFields:     config.ExtraFields,
		devMode:         config.Development,
		fixedKeyValues:  config.FixedKeyValues,
	}

	loggerConfig := zap.NewProductionConfig()
	if logger.devMode {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	loggerConfig.EncoderConfig.MessageKey = "message"
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zLogger, err := loggerConfig.Build(
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, err
	}

	logger.logger = zLogger.Sugar()
	return logger, nil
}

func (l *Logger) Debug(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Debugw(msg, combinedAttributes...)
}

func (l *Logger) Info(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Infow(msg, combinedAttributes...)
}

func (l *Logger) Warn(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Warnw(msg, combinedAttributes...)
}

func (l *Logger) Error(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Errorw(msg, combinedAttributes...)
}

func (l *Logger) Panic(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Panicw(msg, combinedAttributes...)
}

func (l *Logger) Fatal(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Fatalw(msg, combinedAttributes...)
}

func (l *Logger) Debugf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Debugw(msg, combinedAttributes...)
}

func (l *Logger) Infof(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Infow(msg, combinedAttributes...)
}

func (l *Logger) Warnf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Warnw(msg, combinedAttributes...)
}

func (l *Logger) Errorf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Errorw(msg, combinedAttributes...)
}

func (l *Logger) Panicf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Panicw(msg, combinedAttributes...)
}

func (l *Logger) Fatalf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	combinedAttributes := l.combineAttributes(ctx)
	l.logger.Fatalw(msg, combinedAttributes...)
}

func (l *Logger) Debugw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := l.combineAttributes(ctx, keysAndValues...)
	l.logger.Debugw(msg, combinedAttributes...)
}

func (l *Logger) Infow(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := l.combineAttributes(ctx, keysAndValues...)
	l.logger.Infow(msg, combinedAttributes...)
}

func (l *Logger) Warnw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := l.combineAttributes(ctx, keysAndValues...)
	l.logger.Warnw(msg, combinedAttributes...)
}

func (l *Logger) Errorw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := l.combineAttributes(ctx, keysAndValues...)
	l.logger.Errorw(msg, combinedAttributes...)
}

func (l *Logger) Panicw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := l.combineAttributes(ctx, keysAndValues...)
	l.logger.Panicw(msg, combinedAttributes...)
}

func (l *Logger) Fatalw(ctx context.Context, msg string, keysAndValues ...any) {
	combinedAttributes := l.combineAttributes(ctx, keysAndValues...)
	l.logger.Fatalw(msg, combinedAttributes...)
}

func (l *Logger) Flush() {
	l.logger.Sync()
}

func (l *Logger) IsDevMode() bool {
	return l.devMode
}

func (l *Logger) GenerateRequestID() string {
	return l.requestIDPrefix + uuid.New().String()
}

func (l *Logger) SetRequestID(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}

func (l *Logger) GetRequestID(ctx context.Context) (string, bool) {
	requestId, ok := ctx.Value(requestIdKey).(string)
	return requestId, ok
}

func (l *Logger) SetUser(ctx context.Context, user any) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func (l *Logger) GetUser(ctx context.Context) (any, bool) {
	user := ctx.Value(userKey)
	return user, user != nil
}

func (l *Logger) GetExtraFields(ctx context.Context) (map[string]any, bool) {
	if len(l.extraFields) == 0 {
		return nil, false
	}

	pairs := make(map[string]any)
	for _, field := range l.extraFields {
		if value := ctx.Value(field); value != nil {
			pairs[field] = value
		}
	}

	return pairs, true
}

// DetachContext creates a new background context with logging values copied from the original context.
// Use this when starting goroutines that may outlive the original request context.
func (l *Logger) DetachContext(ctx context.Context) context.Context {
	newCtx := context.Background()

	// Copy request ID
	if requestID, ok := l.GetRequestID(ctx); ok {
		newCtx = l.SetRequestID(newCtx, requestID)
	}

	// Copy user
	if user, ok := l.GetUser(ctx); ok {
		newCtx = l.SetUser(newCtx, user)
	}

	// Copy extra fields
	if extraFields, ok := l.GetExtraFields(ctx); ok {
		for key, value := range extraFields {
			newCtx = context.WithValue(newCtx, key, value)
		}
	}

	return newCtx
}

// WithTimeout creates a detached context with timeout that preserves logging values.
// Use this for async operations that need both timeout and logging context.
func (l *Logger) WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	detachedCtx := l.DetachContext(ctx)
	return context.WithTimeout(detachedCtx, timeout)
}

func (l *Logger) combineAttributes(ctx context.Context, keysAndValues ...any) []any {
	var combined []any

	for k, v := range l.fixedKeyValues {
		combined = append(combined, k, v)
	}
	if requestId, ok := l.GetRequestID(ctx); ok {
		combined = append(combined, requestIdContextKey, requestId)
	}
	if user, ok := l.GetUser(ctx); ok {
		combined = append(combined, userContextKey, user)
	}
	if extraFields, ok := l.GetExtraFields(ctx); ok {
		for k, v := range extraFields {
			combined = append(combined, k, v)
		}
	}

	combined = append(combined, keysAndValues...)
	return combined
}

func (l *Logger) LoggerMiddleware(logRequestDetails bool, logCompleteTime bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			requestId := l.GenerateRequestID()
			r = r.WithContext(l.SetRequestID(r.Context(), requestId))

			userIP := getRealUserIP(r)

			if logRequestDetails {
				requestData := map[string]any{
					// Basic request info
					"method":       r.Method,
					"url":          r.URL.String(),
					"path":         r.URL.Path,
					"query_params": r.URL.RawQuery,
					"protocol":     r.Proto,
					"host":         r.Host,

					// Client information
					"user_ip":     userIP,
					"remote_addr": r.RemoteAddr,
					"user_agent":  r.Header.Get("User-Agent"),
					"referer":     r.Header.Get("Referer"),

					// Request size and content
					"content_type":    r.Header.Get("Content-Type"),
					"content_length":  r.ContentLength,
					"accept":          r.Header.Get("Accept"),
					"accept_encoding": r.Header.Get("Accept-Encoding"),
					"accept_language": r.Header.Get("Accept-Language"),

					// Security headers
					"origin": r.Header.Get("Origin"),

					// Load balancer / proxy headers
					"x_forwarded_for":   r.Header.Get("X-Forwarded-For"),
					"x_forwarded_proto": r.Header.Get("X-Forwarded-Proto"),
					"x_forwarded_host":  r.Header.Get("X-Forwarded-Host"),
					"x_real_ip":         r.Header.Get("X-Real-IP"),
					"x_client_ip":       r.Header.Get("X-Client-IP"),
				}

				l.Infow(r.Context(), "Incoming request", "details", requestData)
			}

			next.ServeHTTP(w, r)

			latency := time.Since(now)

			if logCompleteTime {
				l.Infow(r.Context(), "Request completed", "latency", latency)
			}
		})
	}
}

func getRealUserIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	if xci := r.Header.Get("X-Client-IP"); xci != "" {
		return xci
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
