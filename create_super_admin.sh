#!/bin/bash

# Super Admin User Creation Script
# Run this script after deploying the application to create the super admin user

API_URL="${1:-http://localhost:8080}"
SUPER_ADMIN_EMAIL="fluke_tg@yourdomain.com"  # Update this with your actual domain
SUPER_ADMIN_PASSWORD="Ais@07Aut"
SUPER_ADMIN_USERNAME="fluke_tg"

echo "🔐 Creating Super Admin User for Tiger FastTrack Card API"
echo "========================================================"
echo "API URL: $API_URL"
echo "Email: $SUPER_ADMIN_EMAIL"
echo "Username: $SUPER_ADMIN_USERNAME"
echo "Role: super_admin"
echo ""

# Check if API is running
echo "1️⃣  Checking API health..."
HEALTH_RESPONSE=$(curl -s "$API_URL/health" -w "%{http_code}")
HTTP_CODE="${HEALTH_RESPONSE: -3}"

if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ API is running"
else
    echo "   ❌ API is not responding (HTTP: $HTTP_CODE)"
    echo "   Please ensure the API server is running at $API_URL"
    exit 1
fi

# Create super admin user
echo ""
echo "2️⃣  Creating super admin user..."

REGISTER_RESPONSE=$(curl -s -w "%{http_code}" -X POST "$API_URL/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$SUPER_ADMIN_EMAIL\",
        \"password\": \"$SUPER_ADMIN_PASSWORD\",
        \"name\": \"Super Administrator\",
        \"role\": \"super_admin\"
    }")

HTTP_CODE="${REGISTER_RESPONSE: -3}"
RESPONSE_BODY="${REGISTER_RESPONSE%???}"

echo "   HTTP Status: $HTTP_CODE"
echo "   Response: $RESPONSE_BODY"

if [ "$HTTP_CODE" = "201" ] || [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Super admin user created successfully"
elif [ "$HTTP_CODE" = "400" ] && [[ $RESPONSE_BODY == *"already exists"* ]]; then
    echo "   ⚠️  User already exists, continuing with login test..."
else
    echo "   ❌ Failed to create super admin user"
    exit 1
fi

# Test login with super admin credentials
echo ""
echo "3️⃣  Testing super admin login..."

LOGIN_RESPONSE=$(curl -s -w "%{http_code}" -X POST "$API_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$SUPER_ADMIN_EMAIL\",
        \"password\": \"$SUPER_ADMIN_PASSWORD\"
    }")

HTTP_CODE="${LOGIN_RESPONSE: -3}"
RESPONSE_BODY="${LOGIN_RESPONSE%???}"

if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Super admin login successful"
    
    # Extract token for testing
    TOKEN=$(echo $RESPONSE_BODY | grep -o '"token":"[^"]*"' | sed 's/"token":"\(.*\)"/\1/')
    
    if [ -n "$TOKEN" ]; then
        echo "   🔑 JWT Token obtained: ${TOKEN:0:50}..."
        
        # Test admin endpoint
        echo ""
        echo "4️⃣  Testing admin access..."
        
        ADMIN_TEST=$(curl -s -w "%{http_code}" -X GET "$API_URL/api/v1/card-owners/all" \
            -H "Authorization: Bearer $TOKEN")
        
        HTTP_CODE="${ADMIN_TEST: -3}"
        
        if [ "$HTTP_CODE" = "200" ]; then
            echo "   ✅ Admin access confirmed"
        else
            echo "   ⚠️  Admin access test returned HTTP $HTTP_CODE"
        fi
    fi
else
    echo "   ❌ Super admin login failed (HTTP: $HTTP_CODE)"
    echo "   Response: $RESPONSE_BODY"
    exit 1
fi

echo ""
echo "🎉 Super Admin Setup Complete!"
echo "==============================="
echo ""
echo "📋 Super Admin Credentials:"
echo "   Email: $SUPER_ADMIN_EMAIL"
echo "   Username: $SUPER_ADMIN_USERNAME"
echo "   Password: $SUPER_ADMIN_PASSWORD"
echo "   Role: super_admin"
echo ""
echo "🔗 Next Steps:"
echo "   1. Update Postman environment with your production domain"
echo "   2. Use 'Super Admin Login' request in Postman collection"
echo "   3. Test all admin endpoints with super admin credentials"
echo "   4. Consider changing password after initial setup"
echo ""
echo "⚠️  Security Reminder:"
echo "   - Store credentials securely"
echo "   - Consider changing password after deployment"
echo "   - Limit super admin access to trusted personnel only"
echo ""
echo "Happy deploying! 🚀"
