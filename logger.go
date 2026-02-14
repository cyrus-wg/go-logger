package logger

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
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
	loggerConfig.EncoderConfig.TimeKey = "@timestamp"
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

func (l *Logger) LoggerMiddleware(logRequestDetails bool, logCompleteTime bool, bypassList ...BypassRequestLogging) func(next http.Handler) http.Handler {
	compiledBypassList := compileBypassPatterns(bypassList)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			requestId := l.GenerateRequestID()
			r = r.WithContext(l.SetRequestID(r.Context(), requestId))

			shouldSkipLogging := shouldBypassMiddlewareLogging(compiledBypassList, r.URL.Path, r.Method)

			if logRequestDetails && !shouldSkipLogging {
				requestData := map[string]any{
					// Basic request info
					"method":       r.Method,
					"url":          r.URL.String(),
					"path":         r.URL.Path,
					"query_params": r.URL.RawQuery,
					"protocol":     r.Proto,
					"host":         r.Host,

					// Client information
					"user_ip":     getRealUserIP(r),
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

			latency := time.Since(startTime)

			if logCompleteTime && !shouldSkipLogging {
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

type BypassRequestLogging struct {
	Path    string
	Methods string         // Comma-separated methods (e.g., "GET,POST"), empty means all methods
	IsRegex bool           // If true, Path is treated as a regex pattern; otherwise Ant-style pattern
	regex   *regexp.Regexp // Pre-compiled regex (internal use)
}

// compileBypassPatterns pre-compiles regex patterns automatically.
// This is called internally when middleware is initialized.
func compileBypassPatterns(patterns []BypassRequestLogging) []BypassRequestLogging {
	compiled := make([]BypassRequestLogging, len(patterns))
	for i, p := range patterns {
		compiled[i] = p
		if p.IsRegex {
			if re, err := regexp.Compile("^" + p.Path + "$"); err == nil {
				compiled[i].regex = re
			}
		}
	}
	return compiled
}

func shouldBypassMiddlewareLogging(bypassList []BypassRequestLogging, path string, method string) bool {
	for _, bypass := range bypassList {
		if !matchMethod(bypass.Methods, method) {
			continue
		}

		if bypass.IsRegex {
			if matchRegex(&bypass, path) {
				return true
			}
		} else {
			if matchAntPattern(bypass.Path, path) {
				return true
			}
		}
	}

	return false
}

// matchMethod checks if the request method matches the allowed methods.
// Empty methods string means all methods are allowed.
func matchMethod(methods string, method string) bool {
	if methods == "" {
		return true
	}

	// Fast path for single method (no comma)
	if !strings.Contains(methods, ",") {
		return strings.EqualFold(strings.TrimSpace(methods), method)
	}

	// Split and check each method
	for m := range strings.SplitSeq(methods, ",") {
		if strings.EqualFold(strings.TrimSpace(m), method) {
			return true
		}
	}
	return false
}

// matchAntPattern matches path against Ant-style pattern (Spring Security style)
// Supports:
//   - ? matches one character
//   - * matches zero or more characters within a path segment
//   - ** matches zero or more path segments
func matchAntPattern(pattern, path string) bool {
	// Exact match
	if pattern == path {
		return true
	}

	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	return matchAntParts(patternParts, pathParts)
}

func matchAntParts(patternParts, pathParts []string) bool {
	patternIdx, pathIdx := 0, 0

	for patternIdx < len(patternParts) && pathIdx < len(pathParts) {
		patternPart := patternParts[patternIdx]

		if patternPart == "**" {
			// ** at the end matches everything
			if patternIdx == len(patternParts)-1 {
				return true
			}

			// Try to match ** with varying number of path segments
			for i := pathIdx; i <= len(pathParts); i++ {
				if matchAntParts(patternParts[patternIdx+1:], pathParts[i:]) {
					return true
				}
			}
			return false
		}

		if !matchSegment(patternPart, pathParts[pathIdx]) {
			return false
		}

		patternIdx++
		pathIdx++
	}

	// Handle trailing ** in pattern
	for patternIdx < len(patternParts) && patternParts[patternIdx] == "**" {
		patternIdx++
	}

	return patternIdx == len(patternParts) && pathIdx == len(pathParts)
}

// matchSegment matches a single path segment against a pattern segment
// Supports ? (one char) and * (zero or more chars within segment)
func matchSegment(pattern, segment string) bool {
	if pattern == "*" {
		return true
	}

	return matchWildcard(pattern, segment)
}

// matchWildcard matches string against pattern with ? and * wildcards
func matchWildcard(pattern, str string) bool {
	pLen, sLen := len(pattern), len(str)
	pIdx, sIdx := 0, 0
	starIdx, matchIdx := -1, 0

	for sIdx < sLen {
		if pIdx < pLen && (pattern[pIdx] == '?' || pattern[pIdx] == str[sIdx]) {
			pIdx++
			sIdx++
		} else if pIdx < pLen && pattern[pIdx] == '*' {
			starIdx = pIdx
			matchIdx = sIdx
			pIdx++
		} else if starIdx != -1 {
			pIdx = starIdx + 1
			matchIdx++
			sIdx = matchIdx
		} else {
			return false
		}
	}

	for pIdx < pLen && pattern[pIdx] == '*' {
		pIdx++
	}

	return pIdx == pLen
}

// matchRegex matches path against a regex pattern
// Uses pre-compiled regex if available, otherwise compiles on the fly
func matchRegex(bypass *BypassRequestLogging, path string) bool {
	// Use pre-compiled regex if available
	if bypass.regex != nil {
		return bypass.regex.MatchString(path)
	}

	// Fallback: compile and match (less efficient)
	matched, err := regexp.MatchString("^"+bypass.Path+"$", path)
	if err != nil {
		return false
	}
	return matched
}
