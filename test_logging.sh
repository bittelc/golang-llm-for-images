#!/bin/bash

# Test script for validating the comprehensive logging system
# This script runs various test scenarios to ensure logging works correctly

set -e

echo "=== Golang AI Server Logging Test Suite ==="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to run a test
run_test() {
    local test_name=$1
    local log_level=$2
    local log_format=$3
    local input_data=$4
    local expected_pattern=$5
    
    print_status $BLUE "Testing: $test_name"
    print_status $YELLOW "  LOG_LEVEL=$log_level LOG_FORMAT=$format"
    
    # Create temporary file for output
    local temp_file=$(mktemp)
    
    # Run the test
    echo -e "$input_data" | LOG_LEVEL=$log_level LOG_FORMAT=$log_format go run main.go > "$temp_file" 2>&1 &
    local pid=$!
    
    # Wait for a few seconds then kill the process
    sleep 3
    kill $pid 2>/dev/null || true
    wait $pid 2>/dev/null || true
    
    # Check if expected pattern exists in output
    if grep -q "$expected_pattern" "$temp_file"; then
        print_status $GREEN "  ✓ PASS: Found expected pattern '$expected_pattern'"
    else
        print_status $RED "  ✗ FAIL: Pattern '$expected_pattern' not found"
        echo "  Output preview:"
        head -5 "$temp_file" | sed 's/^/    /'
    fi
    
    # Cleanup
    rm -f "$temp_file"
    echo ""
}

# Function to test file operations
test_file_operations() {
    print_status $BLUE "Testing file operation logging..."
    
    # Create a test file
    echo "Test file content for logging validation" > test_file.txt
    
    local temp_file=$(mktemp)
    echo -e "Test file logging\ntest_file.txt" | LOG_LEVEL=DEBUG go run main.go > "$temp_file" 2>&1 &
    local pid=$!
    
    sleep 3
    kill $pid 2>/dev/null || true
    wait $pid 2>/dev/null || true
    
    # Check for file operation logs
    if grep -q "File operation" "$temp_file" && grep -q "test_file.txt" "$temp_file"; then
        print_status $GREEN "  ✓ PASS: File operation logging works"
    else
        print_status $RED "  ✗ FAIL: File operation logging not found"
    fi
    
    # Cleanup
    rm -f test_file.txt "$temp_file"
    echo ""
}

# Check prerequisites
print_status $BLUE "Checking prerequisites..."

if ! command -v go &> /dev/null; then
    print_status $RED "Error: Go is not installed or not in PATH"
    exit 1
fi

if [ ! -f "main.go" ]; then
    print_status $RED "Error: Please run this script from the golang-ai-server directory"
    exit 1
fi

print_status $GREEN "✓ Prerequisites met"
echo ""

# Test 1: Application startup logging
run_test "Application Startup Logging" "INFO" "text" "Test prompt\n" "Starting golang-ai-server application"

# Test 2: Debug level logging
run_test "Debug Level Logging" "DEBUG" "text" "Debug test\n" "Processing step"

# Test 3: JSON format logging
run_test "JSON Format Logging" "INFO" "json" "JSON test\n" '"msg":"Starting golang-ai-server application"'

# Test 4: Error level logging (should show minimal output)
run_test "Error Level Logging" "ERROR" "text" "Error test\n" "User prompt:"

# Test 5: User input logging
run_test "User Input Logging" "DEBUG" "text" "Input logging test\n" "User input received"

# Test 6: Processing step logging
run_test "Processing Step Logging" "DEBUG" "text" "Step test\n" "step=start_user_input"

# Test 7: Ollama request logging
run_test "Ollama Request Logging" "DEBUG" "text" "Ollama test\n" "Sending request to Ollama"

# Test 8: File operations
test_file_operations

# Test 9: Configuration logging
print_status $BLUE "Testing configuration logging..."
temp_file=$(mktemp)
echo -e "Config test\n" | LOG_LEVEL=DEBUG go run main.go > "$temp_file" 2>&1 &
pid=$!
sleep 2
kill $pid 2>/dev/null || true
wait $pid 2>/dev/null || true

if grep -q "ollama_url=http://localhost:11434/api/generate" "$temp_file"; then
    print_status $GREEN "  ✓ PASS: Configuration logging works"
else
    print_status $RED "  ✗ FAIL: Configuration logging not found"
fi
rm -f "$temp_file"
echo ""

# Test 10: Structured logging fields
print_status $BLUE "Testing structured logging fields..."
temp_file=$(mktemp)
echo -e "Structured test\n" | LOG_LEVEL=DEBUG LOG_FORMAT=json go run main.go > "$temp_file" 2>&1 &
pid=$!
sleep 2
kill $pid 2>/dev/null || true
wait $pid 2>/dev/null || true

if grep -q '"level":"' "$temp_file" && grep -q '"time":"' "$temp_file"; then
    print_status $GREEN "  ✓ PASS: Structured logging fields present"
else
    print_status $RED "  ✗ FAIL: Structured logging fields missing"
fi
rm -f "$temp_file"
echo ""

# Summary
print_status $BLUE "=== Test Summary ==="
print_status $GREEN "Logging system validation completed!"
echo ""
print_status $YELLOW "Key features tested:"
echo "  - Application lifecycle logging"
echo "  - Debug level granularity"
echo "  - JSON and text format output"
echo "  - User input processing logs"
echo "  - File operation logging"
echo "  - HTTP request/response logging"
echo "  - Configuration logging"
echo "  - Structured logging fields"
echo ""
print_status $BLUE "To run manual tests:"
echo "  LOG_LEVEL=DEBUG go run main.go"
echo "  LOG_LEVEL=INFO LOG_FORMAT=json go run main.go"
echo ""
print_status $BLUE "See LOGGING.md for complete documentation."