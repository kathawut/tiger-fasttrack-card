# Tiger FastTrack Card API - Postman Collection

This directory contains the complete Postman collection and environment files for testing the Tiger FastTrack Card API.

## 📁 Files Included

1. **`Tiger_FastTrack_Card_API.postman_collection.json`** - Complete API collection
2. **`Tiger_FastTrack_Card_Development.postman_environment.json`** - Development environment
3. **`Tiger_FastTrack_Card_Production.postman_environment.json`** - Production environment template

## 🚀 Quick Setup

### 1. Import Collection
1. Open Postman
2. Click **Import** button
3. Select `Tiger_FastTrack_Card_API.postman_collection.json`
4. Collection will be imported with all endpoints

### 2. Import Environment
1. In Postman, go to **Environments**
2. Click **Import**
3. Select `Tiger_FastTrack_Card_Development.postman_environment.json`
4. Set as **Active Environment**

### 3. Start Testing
1. Make sure your API server is running on `http://localhost:8080`
2. Start with **Authentication > Login** or **Authentication > Super Admin Login** to get JWT token
3. Token will be automatically saved for subsequent requests

## 📚 Collection Structure

### 🔐 Authentication
- **Register User** - Create new user account
- **Login** - Authenticate and receive JWT token (auto-saves token)
- **Super Admin Login** - Login with super admin credentials (production deployment)
- **Refresh Token** - Refresh expired JWT token

### 👤 User Management
- **Get Profile** - Get current user information
- **Update Profile** - Update user details
- **Change Password** - Change user password

### 🎯 Card Management (Master Data)
- **Get All Cards** - List all available cards
- **Get Card by ID** - Get specific card details
- **Create Card** - Add new card (admin only)
- **Update Card** - Modify card information (admin only)
- **Delete Card** - Remove card (admin only)

### 📋 Card Owner Management
- **Register Single Card** - Register one card to user
- **Register Multiple Cards** - Register multiple cards in one request
- **Get Card Owner Profile** - Get first card (backward compatibility)
- **Get All Card Owner Profiles** - Get all user's cards
- **Update Card Owner** - Modify card owner information
- **Delete Card Owner** - Remove card registration
- **Get All Card Owners** - Admin view of all registrations

### 🔍 Validation & Search APIs (NEW)
- **Validate Duplicate Card Registration** - Check for duplicates before registration
- **Search by Card Name and Number** - Find owners by card details
- **Search by ID Card or Phone** - Find owners by personal information

## 🔧 Environment Variables

### Development Environment
```json
{
  "baseUrl": "http://localhost:8080",
  "authToken": "", // Auto-populated after login
  "refreshToken": "", // For token refresh
  "cardId": "1", // Default card ID for testing
  "cardOwnerId": "1", // Default card owner ID for testing
  "testEmail": "user@example.com",
  "testPassword": "password123",
  "adminEmail": "admin@example.com",
  "adminPassword": "admin123",
  "superAdminUsername": "fluke_tg",
  "superAdminPassword": "Ais@07Aut",
  "superAdminRole": "super_admin"
}
```

### Production Environment
- Update `baseUrl` to your production domain
- Set real email credentials
- Keep sensitive data secure
- **Super Admin Credentials**: `fluke_tg` / `Ais@07Aut` (change after deployment)

## 🔐 **Super Admin Deployment Setup**

### Production Deployment Credentials
When deploying to production, use these pre-configured super admin credentials:

**Username:** `fluke_tg`  
**Password:** `Ais@07Aut`  
**Role:** `super_admin`  

### Setup Steps
1. **Deploy Application** to production environment
2. **Run Super Admin Creation Script:**
   ```bash
   ./create_super_admin.sh https://your-production-domain.com
   ```
3. **Test Super Admin Access** in Postman:
   - Import production environment
   - Use "Super Admin Login" request
   - Verify admin endpoints work

### Security Notes
- ⚠️ **Change password after initial deployment**
- 🔒 **Store credentials securely**
- 👥 **Limit super admin access to trusted personnel only**
- 📝 **Monitor super admin activities**

### Super Admin Capabilities
- Full access to all API endpoints
- Card management (create, update, delete)
- View all user registrations
- Admin-only operations
- System administration functions

---

## 🔑 Authentication Flow

### Automatic Token Management
The collection includes scripts that automatically:
1. **Save JWT token** after successful login
2. **Include Bearer token** in all protected endpoints
3. **Handle token refresh** when needed

### Manual Token Setup
If auto-save doesn't work:
1. Login via **Authentication > Login**
2. Copy `token` from response
3. Set `authToken` environment variable manually

## 📝 Testing Scenarios

### 1. User Registration & Authentication
```
1. POST /auth/register (create account)
2. POST /auth/login (get JWT token)
3. GET /users/profile (verify authentication)
```

### 2. Card Owner Registration
```
1. GET /cards (see available cards)
2. POST /card-owners/validate-duplicate (check for duplicates)
3. POST /card-owners/register (register single card)
4. GET /card-owners/profiles (view registered cards)
```

### 3. Multiple Card Registration
```
1. POST /card-owners/register-multiple (register multiple cards)
2. GET /card-owners/profiles (view all cards)
3. PUT /card-owners/{id} (update specific registration)
```

### 4. Search & Validation
```
1. POST /card-owners/validate-duplicate (validate before registration)
2. GET /card-owners/search/by-card?card_name=Premium (search by card name)
3. GET /card-owners/search/by-owner?id_card=ID123 (search by ID card)
```

## 🎯 New API Endpoints Usage

### Validate Duplicate Registration
```http
POST /api/v1/card-owners/validate-duplicate
Authorization: Bearer {{authToken}}
Content-Type: application/json

{
  "card_id": 1,
  "card_number": "123456789"
}
```

**Response:**
- `200 OK` - Registration is valid
- `409 Conflict` - Duplicate found

### Search by Card Name and Number
```http
GET /api/v1/card-owners/search/by-card?card_name=Premium&card_number=123
Authorization: Bearer {{authToken}}
```

**Features:**
- Partial matching (case-insensitive for card names)
- Optional parameters (can search by either or both)
- AND logic (both conditions must match if provided)

### Search by ID Card or Phone
```http
GET /api/v1/card-owners/search/by-owner?id_card=ID123&phone_number=555
Authorization: Bearer {{authToken}}
```

**Features:**
- Partial matching (case-insensitive for ID cards)
- OR logic (matches either condition)
- Requires at least one parameter

## 🔒 Security Notes

### Role-Based Access
- **Regular Users**: Can only see their own card registrations
- **Admin Users**: Can see all registrations and manage cards

### Protected Endpoints
All endpoints except `/health` and `/auth/*` require JWT authentication:
```
Authorization: Bearer <your-jwt-token>
```

## 🐛 Troubleshooting

### Common Issues

1. **401 Unauthorized**
   - Check if JWT token is set correctly
   - Login again to refresh token

2. **403 Forbidden**
   - Check user role (admin vs regular user)
   - Verify endpoint permissions

3. **409 Conflict (Duplicate Registration)**
   - Card number already registered for that card ID
   - Use different card number or different card

4. **400 Bad Request**
   - Check request body format
   - Verify required fields are provided

### Debug Steps
1. Check environment variables are set
2. Verify server is running on correct port
3. Test with `/health` endpoint first
4. Check console for error messages

## 📞 Support

For API issues or questions:
1. Check server logs for detailed error messages
2. Verify database connection
3. Test with curl commands for debugging
4. Review API documentation

## 🔄 Updates

This collection includes all endpoints as of the latest API version. When new endpoints are added:
1. Re-import the updated collection
2. Update environment variables if needed
3. Test new endpoints with appropriate data

---

**Happy Testing! 🚀**
