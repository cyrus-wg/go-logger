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
    logger.InitLogger(logger.LoggerConfig{Development: true})
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

## HTTP Middleware Example

```go
import (
    "net/http"
    "github.com/cyrus-wg/go-logger"
)

func main() {
    logger.InitLogger(logger.LoggerConfig{Development: true})
    defer logger.Flush()
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logger.Info(r.Context(), "Request received")
        w.Write([]byte("Hello!"))
    })
    http.ListenAndServe(":8080", logger.LoggerMiddleware(mux))
}
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

- `InitLogger(config LoggerConfig) error`
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
- `SetUserIP(ctx, ip)`, `GetUserIP(ctx)`
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

- `LoggerMiddleware(next http.Handler) http.Handler`

## Example: Using Both Global and Instance Loggers

```go
import (
    "context"
    "github.com/cyrus-wg/go-logger"
)

func main() {
    // Global logger
    logger.InitLogger(logger.LoggerConfig{Development: true})
    defer logger.Flush()
    logger.Info(context.Background(), "Global logger in action")

    // Instance logger
    l, _ := logger.NewLogger(logger.LoggerConfig{RequestIDPrefix: "SVC-"})
    defer l.Flush()
    l.Infow(context.Background(), "Instance logger", "service", "api")
}
```

## Examples

See the `example/` directory for demos of API servers and more.
