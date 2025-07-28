# Health Check Troubleshooting Guide

## 🚨 **Health Check Failures**

If you're getting "Readiness probe failed" errors, here are the fixes applied:

### **✅ Fixed Issues:**

1. **SSL Mode Mismatch**
   - Changed from `verify-full` to `require` for better compatibility
   - Updated in: `config.go`, `docker-compose.yml`, `.do/app.yaml`

2. **Health Check Command**
   - Changed from `wget` to `curl` for better compatibility
   - Updated Dockerfile to use Alpine with curl
   - Added proper health check configuration

3. **Startup Logging**
   - Added detailed logging to help diagnose startup issues
   - Logs database connection status
   - Logs migration status

### **🔧 Health Check Configuration:**

```yaml
# In .do/app.yaml
health_check:
  http_path: /health

# In Dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl --fail --silent http://localhost:8080/health || exit 1
```

### **📊 Health Check Response:**

```json
{
  "status": "ok",
  "message": "Tiger FastTrack Card API is running",
  "timestamp": "2025-07-28T08:30:00Z",
  "version": "1.0.0"
}
```

### **🔍 Debugging Steps:**

1. **Check Application Logs:**
   ```bash
   doctl apps logs YOUR_APP_ID
   ```

2. **Verify Environment Variables:**
   - Ensure `JWT_SECRET` is set in App Platform
   - Check database connection variables are auto-injected

3. **Test Health Endpoint Locally:**
   ```bash
   curl http://localhost:8080/health
   ```

4. **Database Connection:**
   - Ensure managed database is running
   - Check SSL mode is `require` not `verify-full`
   - Verify database credentials

### **⚡ Quick Fixes:**

1. **Redeploy with fixes:**
   ```bash
   git add .
   git commit -m "fix: health check and SSL configuration"
   git push origin main
   ```

2. **Or update via doctl:**
   ```bash
   doctl apps update YOUR_APP_ID --spec .do/app.yaml
   ```

### **🚀 Expected Startup Logs:**

```
Starting Tiger FastTrack Card API...
Environment: production
Port: 8080
Database Host: your-db-host
Database SSL Mode: require
Connecting to database...
Database connection established successfully
Running database migrations...
Database migrations completed successfully
Server starting on port 8080
```

### **❌ Common Issues:**

- **Database connection timeout**: Check firewall/network settings
- **SSL verification failed**: Use `require` instead of `verify-full`
- **Missing JWT_SECRET**: Set in App Platform environment variables
- **Port binding**: Ensure app listens on `0.0.0.0:8080`

All these issues have been fixed in the latest deployment configuration.
