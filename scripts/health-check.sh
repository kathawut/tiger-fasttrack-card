#!/bin/bash

# Health check script for Tiger FastTrack Card API

API_URL="${API_URL:-http://localhost:8080}"

echo "🔍 Checking API health at $API_URL..."

# Test health endpoint
response=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/health")

if [ "$response" = "200" ]; then
    echo "✅ API is healthy!"
    
    # Test a few more endpoints
    echo "🧪 Testing additional endpoints..."
    
    # Test cards endpoint
    cards_response=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/api/v1/cards")
    if [ "$cards_response" = "200" ]; then
        echo "✅ Cards endpoint working"
    else
        echo "❌ Cards endpoint failed (HTTP $cards_response)"
    fi
    
    echo "🎉 All tests passed!"
else
    echo "❌ API health check failed (HTTP $response)"
    exit 1
fi
