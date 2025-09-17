#!/bin/bash

echo "🚀 Starting Redis tests with Docker..."

# Clean up any existing test containers
echo "🧹 Cleaning up existing containers..."
docker-compose down redis 2>/dev/null

# Start Redis container
echo "📦 Starting Redis container..."
docker-compose up -d redis

# Wait for Redis to be ready
echo "⏳ Waiting for Redis to be ready..."
sleep 3

# Check if Redis is responding
if docker-compose exec redis redis-cli ping > /dev/null 2>&1; then
    echo "✅ Redis is ready!"
else
    echo "❌ Redis failed to start"
    exit 1
fi

# Run the tests
echo "🧪 Running Redis tests..."
if go test -v ./internal/redis/; then
    echo "✅ All Redis tests passed!"
    TEST_STATUS=0
else
    echo "❌ Some Redis tests failed!"
    TEST_STATUS=1
fi

# Clean up
echo "🧹 Cleaning up..."
docker-compose down redis

echo "🏁 Redis test run complete!"
exit $TEST_STATUS