package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cyrus-wg/go-logger"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <demo_name>")
		fmt.Println("Available demos: dev, prod")
		return
	}

	demoName := os.Args[1]

	switch demoName {
	case "global":
		RunGlobalLoggerDemo()
	case "instance":
		RunInstanceLoggerDemo()
	default:
		fmt.Printf("Unknown demo: %s\n", demoName)
		fmt.Println("Available demos: global, instance")
	}
}

func RunGlobalLoggerDemo() {
	dLogger := logger.GetGlobalLogger()
	if dLogger != nil {
		panic("Global logger already initialized")
	}

	err := logger.InitGlobalLogger(logger.LoggerConfig{
		FixedKeyValues: map[string]any{
			"test-type": "global",
		},
		RequestIDPrefix: "GLOBAL-",
	})

	if err != nil {
		panic(err)
	}

	defer logger.Flush()

	dLogger = logger.GetGlobalLogger()
	if dLogger == nil {
		panic("Failed to retrieve global logger")
	}

	ctx := context.Background()
	logger.Info(ctx, "Global logger initialized successfully")
	logger.Debug(ctx, "This is a debug message won't be shown in production mode")

	mux := http.NewServeMux()

	indexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Context(), "Received request for /")
		time.Sleep(500 * time.Millisecond)
		w.Write([]byte("Index endpoint"))
	})

	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Health check - logging is skipped by middleware but manual logging still works
		logger.Info(r.Context(), "Health check called (manual log)")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})

	metricsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Metrics endpoint - logging is skipped by middleware
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"requests": 100, "errors": 0}`))
	})

	mux.Handle("/", indexHandler)
	mux.Handle("/health", healthHandler)
	mux.Handle("/metrics", metricsHandler)

	port := 8080
	logger.Infof(ctx, "Starting Global Logger Demo on :%d", port)
	logger.Info(ctx, "Endpoints: / (logged), /health (skipped), /metrics (skipped)")

	// Use middleware with skip paths for health and metrics endpoints
	middleware := logger.LoggerMiddleware(true, true, "/health", "/metrics")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), middleware(mux)); err != nil {
		logger.Fatal(ctx, "Server failed to start:", err)
	}
}

func RunInstanceLoggerDemo() {
	loggerInstance, err := logger.NewLogger(logger.LoggerConfig{
		FixedKeyValues: map[string]any{
			"test-type": "instance",
		},
		RequestIDPrefix: "INST-",
		Development:     true,
	})

	if err != nil {
		panic(err)
	}

	defer loggerInstance.Flush()

	devMode := loggerInstance.IsDevMode()
	ctx := context.Background()
	loggerInstance.Infow(ctx, "Logger initialized successfully", "devMode", devMode)
	loggerInstance.Debug(ctx, "This is a debug message from instance logger")

	mux := http.NewServeMux()

	indexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggerInstance.Info(r.Context(), "Received request for /")
		time.Sleep(500 * time.Millisecond)
		w.Write([]byte("Index endpoint"))
	})

	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Health check - logging is skipped by middleware
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})

	readyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Readiness check - logging is skipped by middleware
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ready": true}`))
	})

	mux.Handle("/", indexHandler)
	mux.Handle("/health", healthHandler)
	mux.Handle("/ready", readyHandler)

	port := 8080
	loggerInstance.Infof(ctx, "Starting Instance Logger Demo on :%d", port)
	loggerInstance.Info(ctx, "Endpoints: / (logged), /health (skipped), /ready (skipped)")

	// Use middleware with skip paths for health and readiness endpoints
	middleware := loggerInstance.LoggerMiddleware(true, true, "/health", "/ready")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), middleware(mux)); err != nil {
		loggerInstance.Fatal(ctx, "Server failed to start:", err)
	}
}
