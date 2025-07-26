# DigitalOcean Deployment Checklist

## 🚀 Pre-Deployment Setup

### 1. DigitalOcean Account Setup
- [ ] Create DigitalOcean account
- [ ] Install `doctl` CLI tool
- [ ] Authenticate with DigitalOcean: `doctl auth init`
- [ ] Verify authentication: `doctl account get`

### 2. Repository Setup
- [ ] Code pushed to Git repository (GitHub/GitLab)
- [ ] Main branch is production-ready
- [ ] All environment variables configured
- [ ] Database migrations ready

### 3. Database Setup (Already Done)
- [x] DigitalOcean Managed PostgreSQL database created
- [x] Database credentials configured in .env
- [x] Connection string tested

## 🔧 Deployment Configuration

### 4. App Platform Configuration
- [x] `.do/app.yaml` created with correct specifications
- [x] `Dockerfile` optimized for production
- [x] `.dockerignore` configured
- [x] `docker-compose.yml` updated

### 5. Environment Variables
- [x] Production environment file (`.env.production`)
- [x] JWT secret configured
- [x] Database connection string set
- [x] SSL mode enabled for database

## 🚀 Deployment Process

### 6. Run Deployment Script
```bash
# Make sure you're in the project root
cd /Users/k.kumpor/Projects/tiger-fasttrack-card

# Run the deployment script
./deploy-digitalocean.sh
```

### 7. Choose Deployment Option
- Option 1: Full deployment (recommended for first deployment)
- Option 2: App only (if database already exists)
- Option 3: Super admin only (if app already deployed)
- Option 4: Health check only

## ✅ Post-Deployment Verification

### 8. Health Check
- [ ] API responds to health check: `GET /health`
- [ ] Database connection working
- [ ] SSL certificate valid
- [ ] CORS configured properly

### 9. Super Admin Setup
- [ ] Super admin user created successfully
  - Username: `fluke_tg`
  - Password: `Ais@07Aut`
  - Role: `super_admin`
- [ ] Super admin can login
- [ ] Admin endpoints accessible

### 10. API Testing
- [ ] Import production Postman environment
- [ ] Update base URL to live domain
- [ ] Test authentication endpoints
- [ ] Test CRUD operations
- [ ] Test new validation and search APIs

### 11. Security Verification
- [ ] HTTPS enabled
- [ ] Database connections use SSL
- [ ] JWT tokens working
- [ ] Rate limiting (if implemented)
- [ ] CORS properly configured

## 🔍 Monitoring & Maintenance

### 12. Application Monitoring
- [ ] Set up application logs monitoring
- [ ] Configure error alerts
- [ ] Monitor database performance
- [ ] Set up uptime monitoring

### 13. Backup & Recovery
- [ ] Database backups automated
- [ ] Application code backed up
- [ ] Recovery procedures documented

## 📱 Client Integration

### 14. Frontend Integration
- [ ] Update frontend API endpoints
- [ ] Test client-server communication
- [ ] Update CORS if needed
- [ ] Test authentication flow

### 15. Third-party Integrations
- [ ] Payment gateways (if applicable)
- [ ] External APIs configured
- [ ] Webhooks configured

## 🛠 Useful Commands

### DigitalOcean CLI Commands
```bash
# List all apps
doctl apps list

# Get app details
doctl apps get <APP_ID>

# View app logs
doctl apps logs <APP_ID>

# Update app
doctl apps update <APP_ID> --spec .do/app.yaml

# List databases
doctl databases list

# Get database connection info
doctl databases connection <DB_ID>
```

### Application Commands
```bash
# Build locally
go build -o tiger-fasttrack-card main.go

# Run locally
./tiger-fasttrack-card

# Test API
curl https://your-app-url.ondigitalocean.app/health

# Create super admin
./create_super_admin.sh https://your-app-url.ondigitalocean.app
```

## 🚨 Troubleshooting

### Common Issues
1. **Build Fails**
   - Check Go version in Dockerfile
   - Verify all dependencies in go.mod
   - Check for syntax errors

2. **Database Connection Issues**
   - Verify DATABASE_URL format
   - Check SSL mode settings
   - Ensure database is running

3. **Environment Variables Not Set**
   - Check .do/app.yaml configuration
   - Verify DigitalOcean app settings
   - Update via dashboard or CLI

4. **Health Check Fails**
   - Check if app is running on correct port
   - Verify health endpoint exists
   - Check application logs

5. **Super Admin Creation Fails**
   - Ensure app is fully deployed
   - Check API endpoint accessibility
   - Verify request format

### Debug Commands
```bash
# Check app status
doctl apps get <APP_ID> --format Status

# View recent logs
doctl apps logs <APP_ID> --tail

# Check database status
doctl databases get <DB_ID>
```

## 📞 Support Resources

- [DigitalOcean App Platform Documentation](https://docs.digitalocean.com/products/app-platform/)
- [DigitalOcean Managed Databases](https://docs.digitalocean.com/products/databases/)
- [doctl CLI Reference](https://docs.digitalocean.com/reference/doctl/)

---

## 🎉 Success Indicators

When deployment is successful, you should see:
- ✅ Health check returns `{"status": "ok"}`
- ✅ Super admin can login successfully
- ✅ Database connections working
- ✅ All API endpoints responding
- ✅ HTTPS certificate valid
- ✅ Application logs show no errors

**Your Tiger FastTrack Card API is now live on DigitalOcean! 🚀**
