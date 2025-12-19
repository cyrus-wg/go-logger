# Go Logger

A flexible, context-aware logging framework for Go, wrapping [Uber's Zap](https://github.com/uber-go/zap) with both global and instance logger support. Designed for web applications and services needing structured, contextual, and performant logging.

## Features

- **Global logger**: Simple API, no need to manage logger instances
- **Instance logger**: Create multiple, isolated loggers with different configs
- **Context-aware logging**: Request ID, user, and IP extraction from context
- **HTTP middleware**: Automatic request logging and tracing
- **Development/Production modes**: Configurable log levels and output
- **Structured logging**: Key-value and formatted messages
- **Request tracing**: Automatic request ID generation
- **Real IP detection**: Extracts real client IP from headers
- **Async context support**: Preserve logging context in goroutines

## Installation

```bash
go get github.com/cyrus-wg/go-logger
```

## Quick Start

### Global Logger (Simple Usage)

```go
import (
    "context"
    "github.com/cyrus-wg/go-logger"
)

func main() {
    logger.InitGlobalLogger(logger.LoggerConfig{Development: true})
    defer logger.Flush()
    logger.Info(context.Background(), "Hello from global logger!")
}
```

### Instance Logger (Advanced Usage)

```go
import (
    "context"
    "github.com/cyrus-wg/go-logger"
)

func main() {
    myLogger, _ := logger.NewLogger(logger.LoggerConfig{RequestIDPrefix: "API-"})
    defer myLogger.Flush()
    myLogger.Infof(context.Background(), "Instance logger: %s", "custom config")
}
```

## HTTP Middleware

The logger provides configurable HTTP middleware for automatic request logging and tracing.

### Basic Middleware Usage

```go
import (
    "net/http"
    "github.com/cyrus-wg/go-logger"
)

func main() {
    logger.InitGlobalLogger(logger.LoggerConfig{Development: true})
    defer logger.Flush()
    
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logger.Info(r.Context(), "Request received")
        w.Write([]byte("Hello!"))
    })
    
    // Basic middleware - logs completion time only
    middleware := logger.LoggerMiddleware(false, true)
    http.ListenAndServe(":8080", middleware(mux))
}
```

### Advanced Request Logging

```go
// Enable detailed request logging + completion time
middleware := logger.LoggerMiddleware(true, true)
http.ListenAndServe(":8080", middleware(mux))
```

### Middleware Configuration Options

The `LoggerMiddleware` function accepts the following parameters:

- **`logRequestDetails bool`**: Logs comprehensive request information
- **`logCompleteTime bool`**: Logs request completion with latency
- **`skipPaths ...string`**: Variadic paths to skip logging (e.g., health checks)

```go
func LoggerMiddleware(logRequestDetails bool, logCompleteTime bool, skipPaths ...string) func(http.Handler) http.Handler
```

#### Configuration Examples

```go
// Option 1: Minimal logging (latency only)
logger.LoggerMiddleware(false, true)

// Option 2: Full request details + latency
logger.LoggerMiddleware(true, true)

// Option 3: Request details only (no latency)
logger.LoggerMiddleware(true, false)

// Option 4: No automatic logging (manual logging only)
logger.LoggerMiddleware(false, false)

// Option 5: Full logging but skip endpoints (metrics, health, favicon)
logger.LoggerMiddleware(true, true, "/health", "/metrics", "/favicon.ico")
```

### Skipping Paths

Use the `skipPaths` parameter to exclude specific paths from request logging. This is useful for:

- **Health check endpoints** that are called frequently by load balancers
- **Metrics endpoints** that generate excessive log noise
- **Static assets** like favicon.ico

```go
// Skip health checks and metrics - these won't generate logs
middleware := logger.LoggerMiddleware(true, true, "/health", "/healthz", "/metrics")
http.ListenAndServe(":8080", middleware(mux))
```

**Note:** Skipped paths still get a request ID assigned to the context, so manual logging within those handlers will still include the request ID.

### What Gets Logged

When `logRequestDetails = true`, the middleware logs:

#### Basic Request Info
- HTTP method, URL, path, query parameters
- Protocol version, host header

#### Client Information  
- Real client IP (with proxy header parsing)
- User agent, referer
- Remote address

#### Content Information
- Content type, content length
- Accept headers (accept, accept-encoding, accept-language)

#### Security Context
- Origin header for CORS analysis

#### Infrastructure Headers
- Load balancer headers (`X-Forwarded-*`)
- Proxy headers (`X-Real-IP`, `X-Client-IP`)

### Instance Logger Middleware

```go
myLogger, _ := logger.NewLogger(logger.LoggerConfig{
    RequestIDPrefix: "API-",
    Development: true,
})

middleware := myLogger.LoggerMiddleware(true, true)
http.ListenAndServe(":8080", middleware(mux))
```

## API Overview

### LoggerConfig

```go
type LoggerConfig struct {
    Development     bool
    RequestIDPrefix string
    FixedKeyValues  map[string]any
    ExtraFields     []string
}
```

### Global Logger Functions

- `InitGlobalLogger(config LoggerConfig) error`
- `Flush()`
- `Info(ctx, args...)`, `Debug`, `Warn`, `Error`, `Panic`, `Fatal`
- `Infof(ctx, format, args...)`, ...
- `Infow(ctx, msg, keysAndValues...)`, ...

### Instance Logger Methods

- `NewLogger(config LoggerConfig) (*Logger, error)`
- `(*Logger) Info(ctx, args...)`, ...
- `(*Logger) Infof(ctx, format, args...)`, ...
- `(*Logger) Infow(ctx, msg, keysAndValues...)`, ...
- `(*Logger) Flush()`

### Context Utilities

- `SetRequestID(ctx, id)`, `GetRequestID(ctx)`
- `SetUser(ctx, user)`, `GetUser(ctx)`
- `GenerateRequestID()`

### Async Context Support

- `DetachContext(ctx)` - Create detached context for goroutines
- `WithTimeout(ctx, timeout)` - Detached context with timeout

## Async Context Example

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    ctx = logger.SetUser(ctx, "john.doe")
    
    logger.Info(ctx, "Starting request")
    
    // Start async task that outlives the request
    go func() {
        // Preserve logging context even after request ends
        asyncCtx := logger.DetachContext(ctx)
        
        time.Sleep(5 * time.Second)
        // This still has the request ID and user info
        logger.Info(asyncCtx, "Async task completed")
    }()
    
    w.Write([]byte("Request handled"))
}
```

### Middleware

- `LoggerMiddleware(logRequestDetails bool, logCompleteTime bool, skipPaths ...string) func(http.Handler) http.Handler`

## Example: Using Both Global and Instance Loggers

```go
import (
    "context"
    "github.com/cyrus-wg/go-logger"
)

func main() {
    // Global logger
    logger.InitGlobalLogger(logger.LoggerConfig{Development: true})
    defer logger.Flush()
    logger.Info(context.Background(), "Global logger in action")

    // Instance logger
    l, _ := logger.NewLogger(logger.LoggerConfig{RequestIDPrefix: "SVC-"})
    defer l.Flush()
    l.Infow(context.Background(), "Instance logger", "service", "api")
}
```

## Real-World Examples

### Production API Server

```go
package main

import (
    "context"
    "net/http"
    "github.com/cyrus-wg/go-logger"
)

func main() {
    // Initialize logger for production
    logger.InitGlobalLogger(logger.LoggerConfig{
        Development: false,  // Production mode
        RequestIDPrefix: "PROD-",
        FixedKeyValues: map[string]any{
            "service": "user-api",
            "version": "v1.0.0",
        },
    })
    defer logger.Flush()

    mux := http.NewServeMux()
    
    // API endpoint with context logging
    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // Add user context
        userID := r.Header.Get("X-User-ID")
        if userID != "" {
            ctx = logger.SetUser(ctx, userID)
        }
        
        logger.Infow(ctx, "Processing user request", 
            "endpoint", "/users",
            "method", r.Method,
        )
        
        // Your business logic here
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"users": []}`))
        
        logger.Infow(ctx, "Request processed successfully")
    })
    
    // Health check endpoint (logging will be skipped)
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": "healthy"}`))
    })
    
    // Metrics endpoint (logging will be skipped)
    mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"requests": 1000}`))
    })

    // Use full request logging in production, skip health and metrics endpoints
    middleware := logger.LoggerMiddleware(true, true, "/health", "/metrics")
    
    logger.Info(context.Background(), "Starting production server on :8080")
    http.ListenAndServe(":8080", middleware(mux))
}
```

### Development Server with Debug Logging

```go
func main() {
    // Development configuration
    logger.InitGlobalLogger(logger.LoggerConfig{
        Development: true,  // Debug logs enabled
        RequestIDPrefix: "DEV-",
    })
    defer logger.Flush()

    mux := http.NewServeMux()
    mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // Debug logs only show in development
        logger.Debug(ctx, "Debug: Processing request")
        logger.Infow(ctx, "Request details", "path", r.URL.Path)
        
        w.Write([]byte("Debug response"))
    })

    // Full logging for development debugging
    middleware := logger.LoggerMiddleware(true, true)
    http.ListenAndServe(":3000", middleware(mux))
}
```

### Multiple Service Instances

```go
func main() {
    // Create separate loggers for different services
    authLogger, _ := logger.NewLogger(logger.LoggerConfig{
        RequestIDPrefix: "AUTH-",
        FixedKeyValues: map[string]any{"service": "auth"},
    })
    
    userLogger, _ := logger.NewLogger(logger.LoggerConfig{
        RequestIDPrefix: "USER-",
        FixedKeyValues: map[string]any{"service": "users"},
    })

    // Different middleware configurations
    authMiddleware := authLogger.LoggerMiddleware(true, true)  // Full logging
    userMiddleware := userLogger.LoggerMiddleware(false, true) // Latency only

    // Setup routes with different loggers...
}
```

## Sample Log Output

### With Request Details (`logRequestDetails = true`)

```json
{
  "level": "INFO",
  "@timestamp": "2024-09-28T10:30:45.123Z",
  "message": "Incoming request",
  "request_id": "PROD-550e8400-e29b-41d4-a716-446655440000",
  "details": {
    "method": "POST",
    "url": "https://api.example.com/users?active=true",
    "path": "/users",
    "query_params": "active=true",
    "protocol": "HTTP/1.1",
    "host": "api.example.com",
    "user_ip": "203.0.113.1",
    "remote_addr": "10.0.0.1:54321",
    "user_agent": "Mozilla/5.0 (compatible; APIClient/1.0)",
    "content_type": "application/json",
    "content_length": 156,
    "origin": "https://app.example.com"
  }
}
```

### With Completion Logging (`logCompleteTime = true`)

```json
{
  "level": "INFO", 
  "@timestamp": "2024-09-28T10:30:45.256Z",
  "message": "Request completed",
  "request_id": "PROD-550e8400-e29b-41d4-a716-446655440000",
  "latency": "0.133"
}
```

## Examples

See the `example/` directory for complete demos of API servers and more usage patterns.
