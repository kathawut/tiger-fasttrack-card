# Environment Variables Configuration

This application is configured to use **only environment variables** - no `.env` files are loaded in production.

## 🌍 **DigitalOcean App Platform Setup**

All environment variables are configured in the DigitalOcean App Platform dashboard or via `.do/app.yaml`.

### **Required Environment Variables:**

```yaml
ENVIRONMENT=production
PORT=8080
DATABASE_URL=${DATABASE_URL}    # Auto-injected by DigitalOcean
DB_HOST=${DB_HOST}             # Auto-injected by DigitalOcean
DB_PORT=${DB_PORT}             # Auto-injected by DigitalOcean
DB_USER=${DB_USER}             # Auto-injected by DigitalOcean
DB_PASSWORD=${DB_PASSWORD}     # Auto-injected by DigitalOcean
DB_NAME=${DB_NAME}             # Auto-injected by DigitalOcean
DB_SSLMODE=require
JWT_SECRET=${JWT_SECRET}       # Set this in App Platform
LOG_LEVEL=info
```

### **Setting Environment Variables in DigitalOcean:**

1. **Via App Platform Dashboard:**
   - Go to your app in DigitalOcean
   - Navigate to Settings → App-Level Environment Variables
   - Add `JWT_SECRET` with a secure value

2. **Via doctl CLI:**
   ```bash
   doctl apps update YOUR_APP_ID --spec .do/app.yaml
   ```

### **Database Variables:**
- Database connection variables (`DATABASE_URL`, `DB_HOST`, etc.) are **automatically injected** by DigitalOcean when you have a managed database connected to your app.

### **Security Notes:**
- **JWT_SECRET**: Generate a strong, random secret for production
- **Database credentials**: Managed automatically by DigitalOcean
- **No `.env` files**: All configuration through environment variables only

## 🏗️ **Local Development**

For local development, use `docker-compose.dev.yml` which sets all environment variables directly in the compose file:

```bash
docker-compose -f docker-compose.dev.yml up
```

## 🔒 **Production Deployment**

The application will read all configuration from environment variables set in DigitalOcean App Platform. No local files needed.

## ✅ **Benefits of This Approach:**

- ✅ **Security**: No sensitive data in code repository
- ✅ **12-Factor App**: Environment-based configuration
- ✅ **Cloud-native**: Works seamlessly with App Platform
- ✅ **No file dependencies**: Purely environment-driven
