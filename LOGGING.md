# Logging Documentation

This document describes the comprehensive debug logging system implemented throughout the golang-ai-server codebase.

## Overview

The application now includes structured logging using Go's standard `log/slog` package with custom utilities for consistent logging across all components. The logging system provides detailed debug information for troubleshooting, performance monitoring, and development purposes.

## Configuration

### Environment Variables

The logging system can be configured using environment variables:

- `LOG_LEVEL`: Sets the logging level
  - `DEBUG`: Most verbose, includes all debug information
  - `INFO`: General information messages (default)
  - `WARN`: Warning messages only
  - `ERROR`: Error messages only

- `LOG_FORMAT`: Sets the output format
  - `text`: Human-readable text format (default)
  - `json`: JSON format for structured logging

### Examples

```bash
# Run with debug logging in text format
LOG_LEVEL=DEBUG LOG_FORMAT=text go run main.go

# Run with info logging in JSON format for production
LOG_LEVEL=INFO LOG_FORMAT=json go run main.go

# Run with default settings (INFO level, text format)
go run main.go
```

## Logging Features

### 1. Application Lifecycle Logging
- Application startup and shutdown
- Configuration loading
- Total execution time
- Version and build information

### 2. User Input Logging
- Prompt input collection
- Image path processing
- Input validation
- File size tracking

### 3. File Operations Logging
- File opening and reading
- File size and encoding operations
- Base64 encoding progress
- Error handling for file operations

### 4. HTTP Request/Response Logging
- Ollama API request details
- Request payload size
- Response status codes
- Response processing time
- Header information

### 5. Processing Step Logging
- Each major processing step
- Data transformation operations
- Encoding/decoding operations
- Error contexts with detailed information

## Logger Utilities

### Core Functions

- `logger.InitLogger(config)`: Initialize logging with configuration
- `logger.LogError(operation, err, context)`: Log errors with context
- `logger.LogFileOperation(operation, path, size)`: Log file operations
- `logger.LogUserInput(type, size)`: Log user input operations
- `logger.LogProcessingStep(step, details)`: Log processing steps
- `logger.LogRequest(method, url, size, headers)`: Log HTTP requests
- `logger.LogResponse(status, size, duration, headers)`: Log HTTP responses

### Example Usage

```go
// Log a processing step
logger.LogProcessingStep("encode_image", map[string]interface{}{
    "path": imagePath,
    "size": fileSize,
})

// Log an error with context
logger.LogError("file_read", err, map[string]interface{}{
    "path": filePath,
    "operation": "base64_encode",
})

// Log file operations
logger.LogFileOperation("read_pdf", pdfPath, fileSize)
```

## Log Levels by Component

### Main Application (main.go)
- **INFO**: Application start/stop, request completion
- **DEBUG**: Configuration details, timing information

### Input Package (input/input.go)
- **INFO**: User input collection status
- **DEBUG**: Detailed input processing, file encoding
- **ERROR**: Input validation failures, file read errors

### Ollama Package (ollama/server.go)
- **INFO**: API request initiation and completion
- **DEBUG**: Request/response details, HTTP headers
- **ERROR**: Network failures, JSON parsing errors
- **WARN**: Incomplete responses

## Troubleshooting with Logs

### Common Debug Scenarios

1. **File Reading Issues**
   ```
   LOG_LEVEL=DEBUG go run main.go
   # Look for "Failed to read file" or "File operation" messages
   ```

2. **Network Problems**
   ```
   LOG_LEVEL=DEBUG go run main.go
   # Look for "HTTP request failed" or response status codes
   ```

3. **Performance Issues**
   ```
   LOG_LEVEL=INFO go run main.go
   # Check duration fields in log messages
   ```

### Log Message Structure

All log messages follow a consistent structure:
- **Timestamp**: When the event occurred
- **Level**: Log level (DEBUG, INFO, WARN, ERROR)
- **Message**: Human-readable description
- **Fields**: Structured key-value pairs with context

Example:
```
2024-01-15T10:30:45.123Z INFO Application completed successfully total_duration=2.5s
2024-01-15T10:30:45.124Z DEBUG File operation operation=read_image path=/path/to/image.jpg size_bytes=1024000
```

## Performance Impact

The logging system is designed to have minimal performance impact:
- Debug logging only activates when LOG_LEVEL=DEBUG
- Structured logging is efficient
- File operations and network calls are tracked for performance monitoring
- No logging occurs in hot loops

## Production Considerations

For production deployments:
1. Use `LOG_LEVEL=INFO` or higher
2. Consider `LOG_FORMAT=json` for log aggregation systems
3. Monitor log volume and rotate logs appropriately
4. Use structured logging fields for alerting and monitoring

## Integration with Monitoring

The structured logging format makes it easy to integrate with monitoring systems:
- **Prometheus**: Extract metrics from duration and count fields
- **ELK Stack**: Parse JSON logs for detailed analysis
- **Grafana**: Create dashboards from log metrics
- **AlertManager**: Set up alerts based on error patterns

## Development Tips

When adding new logging:
1. Use appropriate log levels
2. Include relevant context in structured fields
3. Follow the existing naming conventions
4. Test logging at different levels
5. Avoid logging sensitive information (like API keys)

Example of good logging practice:
```go
logger.LogProcessingStep("validate_input", map[string]interface{}{
    "input_type": "image_path",
    "path_count": len(paths),
    "max_allowed": 5,
})
```
