#!/bin/bash

# Jira MCP Test Interface Startup Script

set -e

echo "🚀 Starting Jira MCP Test Interface"
echo "=================================="
echo

# Check if we're in the web directory
if [ ! -f "index.html" ]; then
    echo "❌ Error: Please run this script from the web directory"
    echo "   cd web && ./start.sh"
    exit 1
fi

# Function to check if a port is in use
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        return 0  # Port is in use
    else
        return 1  # Port is free
    fi
}

# Check if MCP server is running
echo "🔍 Checking for Jira MCP server..."
if check_port 8080; then
    echo "✅ Jira MCP server detected on port 8080"
else
    echo "⚠️  Jira MCP server not detected on port 8080"
    echo "   Please start it first:"
    echo "   cd .. && ./jira-mcp -http_port 8080"
    echo
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if test server port is available
TEST_PORT=3000
if check_port $TEST_PORT; then
    echo "⚠️  Port $TEST_PORT is already in use"
    read -p "Use a different port? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter port number: " TEST_PORT
    else
        echo "❌ Exiting..."
        exit 1
    fi
fi

echo "🌐 Starting test interface server on port $TEST_PORT..."
echo

# Start the server
if command -v go >/dev/null 2>&1; then
    echo "📦 Using Go server..."
    PORT=$TEST_PORT go run server.go
elif command -v python3 >/dev/null 2>&1; then
    echo "📦 Using Python server..."
    python3 -m http.server $TEST_PORT
elif command -v python >/dev/null 2>&1; then
    echo "📦 Using Python server..."
    python -m http.server $TEST_PORT
elif command -v node >/dev/null 2>&1; then
    echo "📦 Using Node.js server..."
    npx http-server -p $TEST_PORT
else
    echo "❌ Error: No suitable server found"
    echo "   Please install Go, Python, or Node.js"
    exit 1
fi
