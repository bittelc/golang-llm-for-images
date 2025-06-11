package logger

import (
	"log/slog"
	"os"
	"strings"
)

// LogLevel represents the available log levels
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Config holds logging configuration
type Config struct {
	Level  LogLevel
	Format string // "text" or "json"
}

// InitLogger initializes the global logger with the specified configuration
func InitLogger(config Config) *slog.Logger {
	var level slog.Level
	switch config.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch strings.ToLower(config.Format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	default:
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// GetLogLevelFromEnv returns the log level from environment variable
func GetLogLevelFromEnv() LogLevel {
	envLevel := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch envLevel {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN", "WARNING":
		return LevelWarn
	case "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// GetLogFormatFromEnv returns the log format from environment variable
func GetLogFormatFromEnv() string {
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))
	if format == "json" {
		return "json"
	}
	return "text"
}

// NewDefaultConfig creates a default logging configuration
func NewDefaultConfig() Config {
	return Config{
		Level:  GetLogLevelFromEnv(),
		Format: GetLogFormatFromEnv(),
	}
}

// LogRequest logs details about an HTTP-like request
func LogRequest(method, url string, bodySize int, headers map[string]string) {
	slog.Debug("HTTP Request",
		"method", method,
		"url", url,
		"body_size", bodySize,
		"headers", headers)
}

// LogResponse logs details about an HTTP-like response
func LogResponse(statusCode int, bodySize int, duration string, headers map[string]string) {
	slog.Debug("HTTP Response",
		"status_code", statusCode,
		"body_size", bodySize,
		"duration", duration,
		"headers", headers)
}

// LogError logs an error with context
func LogError(operation string, err error, context map[string]interface{}) {
	args := []interface{}{"operation", operation, "error", err}
	for k, v := range context {
		args = append(args, k, v)
	}
	slog.Error("Operation failed", args...)
}

// LogFileOperation logs file operations
func LogFileOperation(operation, path string, size int64) {
	slog.Debug("File operation",
		"operation", operation,
		"path", path,
		"size_bytes", size)
}

// LogUserInput logs user input operations
func LogUserInput(inputType string, size int) {
	slog.Debug("User input received",
		"type", inputType,
		"size", size)
}

// LogProcessingStep logs a processing step
func LogProcessingStep(step string, details map[string]interface{}) {
	args := []interface{}{"step", step}
	for k, v := range details {
		args = append(args, k, v)
	}
	slog.Debug("Processing step", args...)
}