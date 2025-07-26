#!/bin/bash

# DigitalOcean Deployment Script for Tiger FastTrack Card API
# This script helps deploy the application to DigitalOcean App Platform

set -e

echo "🚀 DigitalOcean Deployment for Tiger FastTrack Card API"
echo "======================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="tiger-fasttrack-card"
REGION="nyc1"  # Change to your preferred region
BRANCH="main"  # Change if deploying from different branch

# Check if doctl is installed
if ! command -v doctl &> /dev/null; then
    echo -e "${RED}❌ doctl (DigitalOcean CLI) is not installed${NC}"
    echo "Please install it from: https://docs.digitalocean.com/reference/doctl/how-to/install/"
    exit 1
fi

# Check if user is authenticated
if ! doctl auth list > /dev/null 2>&1; then
    echo -e "${RED}❌ Not authenticated with DigitalOcean${NC}"
    echo "Please run: doctl auth init"
    exit 1
fi

echo -e "${GREEN}✅ doctl is installed and authenticated${NC}"

# Function to create database if it doesn't exist
create_database() {
    echo -e "${BLUE}📊 Setting up database...${NC}"
    
    # Check if database already exists
    if doctl databases list | grep -q "tiger-card-db"; then
        echo -e "${YELLOW}⚠️  Database 'tiger-card-db' already exists${NC}"
        DB_ID=$(doctl databases list --format ID,Name --no-header | grep "tiger-card-db" | awk '{print $1}')
    else
        echo -e "${BLUE}🔨 Creating PostgreSQL database...${NC}"
        doctl databases create tiger-card-db \
            --engine pg \
            --version 15 \
            --region $REGION \
            --size db-s-1vcpu-1gb \
            --num-nodes 1
        
        # Wait for database to be ready
        echo -e "${YELLOW}⏳ Waiting for database to be ready...${NC}"
        sleep 60
        
        DB_ID=$(doctl databases list --format ID,Name --no-header | grep "tiger-card-db" | awk '{print $1}')
    fi
    
    echo -e "${GREEN}✅ Database ready with ID: $DB_ID${NC}"
    
    # Get database connection details
    DB_CONNECTION=$(doctl databases connection $DB_ID --format URI --no-header)
    echo -e "${GREEN}📋 Database connection string: $DB_CONNECTION${NC}"
}

# Function to deploy app
deploy_app() {
    echo -e "${BLUE}🚀 Deploying application...${NC}"
    
    # Check if app already exists
    if doctl apps list | grep -q "$APP_NAME"; then
        echo -e "${YELLOW}⚠️  App '$APP_NAME' already exists, updating...${NC}"
        APP_ID=$(doctl apps list --format ID,Spec.Name --no-header | grep "$APP_NAME" | awk '{print $1}')
        doctl apps update $APP_ID --spec .do/app.yaml
    else
        echo -e "${BLUE}🔨 Creating new app...${NC}"
        doctl apps create --spec .do/app.yaml
    fi
    
    echo -e "${GREEN}✅ Application deployed successfully${NC}"
}

# Function to setup environment variables
setup_env_vars() {
    echo -e "${BLUE}🔧 Setting up environment variables...${NC}"
    
    # Get app ID
    APP_ID=$(doctl apps list --format ID,Spec.Name --no-header | grep "$APP_NAME" | awk '{print $1}')
    
    if [ -z "$APP_ID" ]; then
        echo -e "${RED}❌ Could not find app ID${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}📋 App ID: $APP_ID${NC}"
    echo -e "${YELLOW}💡 Set environment variables in DigitalOcean App Platform dashboard${NC}"
    echo -e "${YELLOW}💡 Or use: doctl apps update $APP_ID --spec .do/app.yaml${NC}"
}

# Function to create super admin user
create_super_admin() {
    echo -e "${BLUE}👑 Creating super admin user...${NC}"
    
    # Get app URL
    APP_URL=$(doctl apps get $(doctl apps list --format ID,Spec.Name --no-header | grep "$APP_NAME" | awk '{print $1}') --format LiveURL --no-header)
    
    if [ -z "$APP_URL" ]; then
        echo -e "${RED}❌ Could not get app URL${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}🌐 App URL: $APP_URL${NC}"
    
    # Wait for app to be ready
    echo -e "${YELLOW}⏳ Waiting for application to be ready...${NC}"
    sleep 30
    
    # Create super admin user
    echo -e "${BLUE}👤 Creating super admin user...${NC}"
    
    RESPONSE=$(curl -s -w "%{http_code}" -X POST "$APP_URL/api/v1/auth/register" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "fluke_tg@yourdomain.com",
            "password": "Ais@07Aut",
            "name": "Super Administrator",
            "role": "super_admin"
        }')
    
    HTTP_CODE="${RESPONSE: -3}"
    
    if [ "$HTTP_CODE" = "201" ] || [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✅ Super admin user created successfully${NC}"
    elif [ "$HTTP_CODE" = "400" ]; then
        echo -e "${YELLOW}⚠️  Super admin user might already exist${NC}"
    else
        echo -e "${RED}❌ Failed to create super admin user (HTTP: $HTTP_CODE)${NC}"
    fi
}

# Function to run health check
health_check() {
    echo -e "${BLUE}🏥 Running health check...${NC}"
    
    APP_URL=$(doctl apps get $(doctl apps list --format ID,Spec.Name --no-header | grep "$APP_NAME" | awk '{print $1}') --format LiveURL --no-header)
    
    if [ -z "$APP_URL" ]; then
        echo -e "${RED}❌ Could not get app URL${NC}"
        exit 1
    fi
    
    HEALTH_RESPONSE=$(curl -s "$APP_URL/health")
    
    if echo "$HEALTH_RESPONSE" | grep -q "ok"; then
        echo -e "${GREEN}✅ Health check passed${NC}"
        echo -e "${GREEN}🎉 Deployment successful!${NC}"
        echo -e "${GREEN}🌐 Your API is live at: $APP_URL${NC}"
    else
        echo -e "${RED}❌ Health check failed${NC}"
        echo -e "${RED}Response: $HEALTH_RESPONSE${NC}"
    fi
}

# Main deployment flow
echo -e "${BLUE}🔍 Pre-deployment checks...${NC}"

# Check if .do/app.yaml exists
if [ ! -f ".do/app.yaml" ]; then
    echo -e "${RED}❌ .do/app.yaml not found${NC}"
    echo "Please ensure you're in the project root directory"
    exit 1
fi

echo -e "${GREEN}✅ Configuration files found${NC}"

# Ask user what they want to do
echo -e "${BLUE}What would you like to do?${NC}"
echo "1) Full deployment (database + app + super admin)"
echo "2) Deploy app only"
echo "3) Create super admin user only"
echo "4) Health check only"
echo "5) Exit"

read -p "Choose option (1-5): " option

case $option in
    1)
        create_database
        deploy_app
        setup_env_vars
        echo -e "${YELLOW}⏳ Waiting for deployment to complete...${NC}"
        sleep 120
        create_super_admin
        health_check
        ;;
    2)
        deploy_app
        setup_env_vars
        ;;
    3)
        create_super_admin
        ;;
    4)
        health_check
        ;;
    5)
        echo -e "${BLUE}👋 Goodbye!${NC}"
        exit 0
        ;;
    *)
        echo -e "${RED}❌ Invalid option${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}🎉 Deployment process completed!${NC}"
echo ""
echo -e "${BLUE}📋 Next steps:${NC}"
echo "1. Update your Postman environment with the live URL"
echo "2. Test the API endpoints"
echo "3. Monitor the application logs"
echo "4. Set up monitoring and alerts"
echo ""
echo -e "${BLUE}🔧 Useful commands:${NC}"
echo "- View app logs: doctl apps logs \$(doctl apps list --format ID,Spec.Name --no-header | grep $APP_NAME | awk '{print \$1}')"
echo "- View app info: doctl apps get \$(doctl apps list --format ID,Spec.Name --no-header | grep $APP_NAME | awk '{print \$1}')"
echo "- Update app: doctl apps update \$(doctl apps list --format ID,Spec.Name --no-header | grep $APP_NAME | awk '{print \$1}') --spec .do/app.yaml"
echo ""
echo -e "${GREEN}Happy deploying! 🚀${NC}"
