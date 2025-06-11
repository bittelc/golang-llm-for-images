# golang-ai-server

**Go-based AI Server with Comprehensive Debug Logging**

---

## Overview

`golang-ai-server` is a server application written in Go that provides AI processing capabilities with Ollama integration. It supports multi-modal input including text prompts and images, with comprehensive debug logging for development, troubleshooting, and monitoring.

---

## Features

- **Ollama Integration**: Direct integration with Ollama API for AI processing
- **Multi-modal Input**: Supports text prompts and image file processing
- **Comprehensive Logging**: Structured debug logging with multiple levels and formats
- **File Processing**: Base64 encoding for images and PDFs
- **Performance Monitoring**: Request timing and performance metrics
- **Development-Friendly**: Extensive debugging information for development

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.21+ recommended for slog support)
- [Ollama](https://ollama.ai/) running locally on port 11434
- Images or files to process (optional)

### Installation

Clone the repository:

```bash
git clone https://github.com/bittelc/golang-ai-server.git
cd golang-ai-server
```

Install dependencies:

```bash
go mod tidy
```

### Quick Start

Run with default settings:

```bash
go run main.go
```

Run with debug logging:

```bash
LOG_LEVEL=DEBUG go run main.go
```

Run with JSON logging format:

```bash
LOG_LEVEL=INFO LOG_FORMAT=json go run main.go
```

---

## Logging System

This application features a comprehensive logging system for debugging, monitoring, and development purposes.

### Log Levels

- **DEBUG**: Detailed debugging information, file operations, HTTP details
- **INFO**: General application flow and status (default)
- **WARN**: Warning conditions and incomplete responses  
- **ERROR**: Error conditions and failures

### Log Formats

- **text**: Human-readable text format (default)
- **json**: Structured JSON format for log aggregation

### Environment Variables

```bash
# Set log level
export LOG_LEVEL=DEBUG    # DEBUG, INFO, WARN, ERROR

# Set log format  
export LOG_FORMAT=json    # text, json
```

### What Gets Logged

1. **Application Lifecycle**
   - Startup and shutdown events
   - Configuration loading
   - Total execution time

2. **User Input Processing**
   - Prompt collection and validation
   - Image path parsing and validation
   - File reading and encoding operations

3. **File Operations**
   - File opening, reading, and encoding
   - File size tracking and performance
   - Base64 encoding progress

4. **HTTP Operations**
   - Ollama API requests and responses
   - Request/response timing and size
   - HTTP status codes and headers

5. **Error Handling**
   - Detailed error context and stack traces
   - Operation failure points
   - Recovery attempts

### Logging Demo

Run the interactive logging demonstration:

```bash
./demo_logging.sh
```

### Development Tips

View detailed logs during development:

```bash
LOG_LEVEL=DEBUG go run main.go 2>&1 | tee app.log
```

For production-like monitoring:

```bash
LOG_LEVEL=INFO LOG_FORMAT=json go run main.go
```

See [LOGGING.md](LOGGING.md) for complete logging documentation.

---

## Usage

1. **Start the application**:
   ```bash
   go run main.go
   ```

2. **Enter your prompt** when asked

3. **Optionally provide image paths** (comma-separated, max 5):
   ```
   /path/to/image1.jpg, /path/to/image2.png
   ```

4. **View the AI response** and processing logs

### Example Session

```bash
$ LOG_LEVEL=INFO go run main.go
User prompt: Describe this image
Path to images, separated by commas, limit of 5 (optional): test_image.txt
[Processing logs will appear here]
[AI response will appear here]
Completed in 2.5s
```

---

## Architecture

### Components

- **main.go**: Application entry point and orchestration
- **input/**: User input handling and file processing
- **ollama/**: Ollama API client and request handling  
- **logger/**: Centralized logging utilities and configuration

### File Structure

```
golang-ai-server/
├── main.go              # Main application
├── input/
│   └── input.go         # User input and file processing
├── ollama/
│   └── server.go        # Ollama API client
├── logger/
│   └── logger.go        # Logging utilities
├── LOGGING.md           # Logging documentation
├── demo_logging.sh      # Logging demonstration
└── test_image.txt       # Test file for logging demo
```

---

## Development

### Adding New Logging

When adding new functionality, use the logging utilities:

```go
// Log processing steps
logger.LogProcessingStep("operation_name", map[string]interface{}{
    "param1": value1,
    "param2": value2,
})

// Log errors with context
logger.LogError("operation_name", err, map[string]interface{}{
    "context1": value1,
    "context2": value2,
})

// Log file operations
logger.LogFileOperation("read_file", filePath, fileSize)
