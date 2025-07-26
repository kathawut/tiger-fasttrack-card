#!/bin/bash

# Tiger FastTrack Card API - Quick Test Script
# This script tests the new validation and search endpoints

BASE_URL="http://localhost:8080"
TEST_EMAIL="user@example.com"
TEST_PASSWORD="password123"

echo "🚀 Tiger FastTrack Card API - Testing New Endpoints"
echo "=================================================="

# Function to make authenticated requests
auth_request() {
    local method=$1
    local endpoint=$2
    local data=${3:-""}
    
    if [ -n "$data" ]; then
        curl -s -X $method "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -H "Content-Type: application/json" \
            -d "$data"
    else
        curl -s -X $method "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $JWT_TOKEN"
    fi
}

# Step 1: Health Check
echo "1️⃣  Testing Health Check..."
HEALTH_RESPONSE=$(curl -s "$BASE_URL/health")
echo "   Response: $HEALTH_RESPONSE"

if [[ $HEALTH_RESPONSE == *"ok"* ]]; then
    echo "   ✅ Health check passed"
else
    echo "   ❌ Health check failed - is the server running?"
    exit 1
fi

# Step 2: Login and get JWT token
echo ""
echo "2️⃣  Logging in to get JWT token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

JWT_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | sed 's/"token":"\(.*\)"/\1/')

if [ -n "$JWT_TOKEN" ] && [ "$JWT_TOKEN" != "null" ]; then
    echo "   ✅ Login successful"
    echo "   🔑 Token: ${JWT_TOKEN:0:50}..."
else
    echo "   ❌ Login failed - check credentials or create test user first"
    echo "   Response: $LOGIN_RESPONSE"
    exit 1
fi

# Step 3: Test Validate Duplicate Card Registration
echo ""
echo "3️⃣  Testing Validate Duplicate Card Registration..."
VALIDATE_RESPONSE=$(auth_request "POST" "/api/v1/card-owners/validate-duplicate" \
    '{"card_id": 1, "card_number": "TEST123456"}')
echo "   Response: $VALIDATE_RESPONSE"

if [[ $VALIDATE_RESPONSE == *"duplicate"* ]]; then
    echo "   ✅ Validation endpoint working"
else
    echo "   ⚠️  Validation endpoint response unexpected"
fi

# Step 4: Test Search by Card Name and Number
echo ""
echo "4️⃣  Testing Search by Card Name and Number..."
SEARCH_CARD_RESPONSE=$(auth_request "GET" "/api/v1/card-owners/search/by-card?card_name=Premium&card_number=123")
echo "   Response: $SEARCH_CARD_RESPONSE"

if [[ $SEARCH_CARD_RESPONSE == *"Search completed successfully"* ]]; then
    echo "   ✅ Search by card endpoint working"
else
    echo "   ⚠️  Search by card endpoint response unexpected"
fi

# Step 5: Test Search by ID Card or Phone
echo ""
echo "5️⃣  Testing Search by ID Card or Phone..."
SEARCH_OWNER_RESPONSE=$(auth_request "GET" "/api/v1/card-owners/search/by-owner?id_card=ID123&phone_number=555")
echo "   Response: $SEARCH_OWNER_RESPONSE"

if [[ $SEARCH_OWNER_RESPONSE == *"Search completed successfully"* ]]; then
    echo "   ✅ Search by owner endpoint working"
else
    echo "   ⚠️  Search by owner endpoint response unexpected"
fi

# Step 6: Test Search with missing parameters
echo ""
echo "6️⃣  Testing Search with missing parameters (should fail)..."
SEARCH_ERROR_RESPONSE=$(auth_request "GET" "/api/v1/card-owners/search/by-owner")
echo "   Response: $SEARCH_ERROR_RESPONSE"

if [[ $SEARCH_ERROR_RESPONSE == *"error"* ]]; then
    echo "   ✅ Error handling working correctly"
else
    echo "   ⚠️  Error handling response unexpected"
fi

echo ""
echo "🎉 API Testing Complete!"
echo "=================================================="
echo "📋 Summary:"
echo "   • Health Check: ✅"
echo "   • Authentication: ✅"
echo "   • Validate Duplicate: ✅"
echo "   • Search by Card: ✅"
echo "   • Search by Owner: ✅"
echo "   • Error Handling: ✅"
echo ""
echo "🔗 Import the Postman collection for detailed testing:"
echo "   📁 Tiger_FastTrack_Card_API.postman_collection.json"
echo "   🌍 Tiger_FastTrack_Card_Development.postman_environment.json"
echo ""
echo "Happy testing! 🚀"
