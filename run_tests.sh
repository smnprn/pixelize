#!/bin/bash

PROJECT_DIR=$(pwd)
GO_TEST_CMD="go test ./..."

mkdir -p "$TEST_OUTPUT_DIR"

echo "Running Go tests..."
if $GO_TEST_CMD 2>&1; then
    echo "✅ Tests passed."
else
    echo "❌ Tests failed. Check the terminal output for details."
    exit 1
fi

echo "✅ All tests passed successfully."