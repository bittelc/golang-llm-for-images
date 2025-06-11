#!/bin/bash

# Logging Demonstration Script for golang-ai-server
# This script demonstrates different logging levels and formats

echo "=== Golang AI Server Logging Demo ==="
echo ""

# Function to run the application with different log settings
run_demo() {
    local level=$1
    local format=$2
    local description=$3
    
    echo "--- $description ---"
    echo "LOG_LEVEL=$level LOG_FORMAT=$format"
    echo ""
    
    # Create a simple test input for demonstration
    echo -e "Test prompt for logging demo\n" | LOG_LEVEL=$level LOG_FORMAT=$format timeout 10s go run main.go 2>&1 | head -20
    
    echo ""
    echo "Press Enter to continue to next demo..."
    read
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "Error: Please run this script from the golang-ai-server directory"
    exit 1
fi

echo "This demo will show different logging levels and formats."
echo "Note: The application will timeout after 10 seconds for demo purposes."
echo ""
echo "Press Enter to start..."
read

# Demo 1: ERROR level logging
run_demo "ERROR" "text" "ERROR Level - Only errors will be shown"

# Demo 2: WARN level logging
run_demo "WARN" "text" "WARN Level - Warnings and errors"

# Demo 3: INFO level logging (default)
run_demo "INFO" "text" "INFO Level - General information (default)"

# Demo 4: DEBUG level logging
run_demo "DEBUG" "text" "DEBUG Level - Detailed debugging information"

# Demo 5: INFO level with JSON format
run_demo "INFO" "json" "INFO Level with JSON Format - Structured logging"

# Demo 6: DEBUG level with JSON format
run_demo "DEBUG" "json" "DEBUG Level with JSON Format - Full debug in structured format"

echo "=== Demo Complete ==="
echo ""
echo "Logging Features Demonstrated:"
echo "- Application lifecycle logging"
echo "- User input processing logs"
echo "- File operation logging"
echo "- HTTP request/response logging"
echo "- Error handling with context"
echo "- Structured logging with key-value pairs"
echo ""
echo "To use different log levels in your application:"
echo "  LOG_LEVEL=DEBUG go run main.go        # Detailed debugging"
echo "  LOG_LEVEL=INFO go run main.go         # General information (default)"
echo "  LOG_LEVEL=WARN go run main.go         # Warnings and errors only"
echo "  LOG_LEVEL=ERROR go run main.go        # Errors only"
echo ""
echo "To use JSON format (useful for log aggregation):"
echo "  LOG_FORMAT=json go run main.go"
echo ""
echo "See LOGGING.md for complete documentation."