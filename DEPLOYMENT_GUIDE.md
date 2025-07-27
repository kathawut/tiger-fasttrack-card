# Tiger FastTrack Card API - Deployment Guide

## 🚀 Production Deployment Setup

### 1. Environment Variables
Ensure these environment variables are set in your production environment:

```bash
# Database Configuration
DATABASE_URL=postgres://username:password@host:port/database
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=tiger_fasttrack_card
DB_SSLMODE=verify-full

# Application Configuration
PORT=8080
ENVIRONMENT=production
JWT_SECRET=your-super-secure-jwt-secret-key

# CORS Configuration (if needed)
ALLOWED_ORIGINS=https://your-frontend-domain.com
```

### 2. Database Setup
Run the following commands to set up the database:

```bash
# 1. Run database migrations
go run main.go migrate

# 2. Create super admin user (manual SQL or via API)
```

### 3. Super Admin User Creation

**Method 1: Direct API Call (Recommended)**
```bash
# Register super admin user via API
curl -X POST https://your-production-domain.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "fluke_tg@yourdomain.com",
    "password": "Ais@07Aut",
    "name": "Super Administrator",
    "role": "super_admin"
  }'
```

**Method 2: Direct Database Insert**
```sql
-- Insert super admin user directly into database
INSERT INTO users (email, password, name, role, is_active, created_at, updated_at)
VALUES (
    'fluke_tg@yourdomain.com',
    '$2a$10$encrypted_password_hash', -- Use bcrypt to hash 'Ais@07Aut'
    'Super Administrator',
    'super_admin',
    true,
    NOW(),
    NOW()
);
```

### 4. Production Security Checklist

#### ✅ Environment Security
- [ ] Strong JWT secret (minimum 32 characters)
- [ ] Database credentials secured
- [ ] HTTPS enabled
- [ ] CORS properly configured

#### ✅ Database Security
- [ ] Database user has minimal required permissions
- [ ] Database connection uses SSL
- [ ] Regular backups configured
- [ ] Database firewall rules in place

#### ✅ Application Security
- [ ] Rate limiting enabled
- [ ] Input validation active
- [ ] Error messages don't expose sensitive data
- [ ] Logging configured (without sensitive data)

### 5. Postman Environment Setup

The production Postman environment includes:

```json
{
  "baseUrl": "https://your-production-domain.com",
  "superAdminUsername": "fluke_tg",
  "superAdminPassword": "Ais@07Aut",
  "superAdminRole": "super_admin"
}
```

**Usage in Postman:**
1. Import `Tiger_FastTrack_Card_Production.postman_environment.json`
2. Update `baseUrl` to your actual production domain
3. Use super admin credentials for testing admin endpoints

### 6. Health Check & Monitoring

**Health Check Endpoint:**
```bash
curl https://your-production-domain.com/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "message": "Tiger FastTrack Card API is running"
}
```

### 7. Deployment Commands

```bash
# Build the application
go build -o tiger-fasttrack-card main.go

# Run with production environment
ENVIRONMENT=production ./tiger-fasttrack-card

# Or with systemd service
sudo systemctl start tiger-fasttrack-card
sudo systemctl enable tiger-fasttrack-card
```

### 8. Initial Test Sequence

After deployment, run these tests:

```bash
# 1. Health check
curl https://your-domain.com/health

# 2. Super admin login
curl -X POST https://your-domain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"fluke_tg@yourdomain.com","password":"Ais@07Aut"}'

# 3. Create test card
curl -X POST https://your-domain.com/api/v1/cards \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"card_name":"Test Card","card_image":"test.jpg","card_quantity":100}'

# 4. Test validation endpoint
curl -X POST https://your-domain.com/api/v1/card-owners/validate-duplicate \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"card_id":1,"card_number":"TEST123"}'
```

### 9. Backup & Recovery

**Database Backup:**
```bash
# Daily backup script
pg_dump $DATABASE_URL > backup_$(date +%Y%m%d).sql
```

**Application Logs:**
```bash
# Monitor application logs
tail -f /var/log/tiger-fasttrack-card.log
```

### 10. Troubleshooting

**Common Issues:**

1. **Database Connection Failed**
   - Check DATABASE_URL format
   - Verify network connectivity
   - Check SSL configuration

2. **JWT Token Issues**
   - Verify JWT_SECRET is set
   - Check token expiration time
   - Validate token format

3. **CORS Errors**
   - Update ALLOWED_ORIGINS
   - Check preflight requests
   - Verify header configuration

**Debug Commands:**
```bash
# Check environment variables
printenv | grep -E "(DB_|JWT_|PORT)"

# Test database connection
psql $DATABASE_URL -c "SELECT 1"

# Check application status
systemctl status tiger-fasttrack-card
```

---

## 🔐 Super Admin Credentials (PRODUCTION)

**Username:** fluke_tg  
**Password:** Ais@07Aut  
**Role:** super_admin  

⚠️ **Security Note:** Change these credentials after initial deployment and setup!

---

## 📞 Support

For deployment issues:
1. Check application logs
2. Verify environment variables
3. Test database connectivity
4. Review security configurations

**Remember to update the production domain in all configuration files before deployment!**
