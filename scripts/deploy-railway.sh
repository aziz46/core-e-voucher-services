#!/bin/bash

# Railway Deployment Script
# Automated setup and deployment ke Railway

set -e

echo "==================================="
echo "Core E-Voucher Railway Deployment"
echo "==================================="
echo ""

# Check prerequisites
check_prerequisites() {
    echo "Checking prerequisites..."
    
    if ! command -v git &> /dev/null; then
        echo "❌ Git is not installed"
        exit 1
    fi
    
    if ! command -v railway &> /dev/null; then
        echo "⚠️  Railway CLI not found. Installing..."
        npm install -g @railway/cli
    fi
    
    echo "✅ Prerequisites OK"
}

# Login to Railway
login_railway() {
    echo ""
    echo "Logging in to Railway..."
    railway login
    echo "✅ Railway login successful"
}

# Create project
create_project() {
    echo ""
    echo "Creating Railway project..."
    read -p "Enter project name (default: core-e-voucher): " project_name
    project_name=${project_name:-core-e-voucher}
    
    railway init
    echo "✅ Project created: $project_name"
}

# Create services
create_services() {
    echo ""
    echo "Setting up services..."
    
    echo "1/5 Adding PostgreSQL database..."
    railway add --template postgres
    
    echo "2/5 Adding Credit Service..."
    railway add --dockerfile ./cmd/credit-service/Dockerfile
    
    echo "3/5 Adding Billing Service..."
    railway add --dockerfile ./cmd/billing-service/Dockerfile
    
    echo "4/5 Adding PPOB Core Service..."
    railway add --dockerfile ./cmd/ppob-core/Dockerfile
    
    echo "5/5 Configuring environment variables..."
    railway variables
    
    echo "✅ All services created"
}

# Run migrations
run_migrations() {
    echo ""
    echo "Running database migrations..."
    
    railway database exec < migrations/001_init.sql
    echo "✅ Schema created"
    
    railway database exec < migrations/002_seed.sql
    echo "✅ Seed data inserted"
}

# Deploy
deploy() {
    echo ""
    echo "Deploying services..."
    
    railway up --service credit-service
    railway up --service billing-service
    railway up --service ppob-core
    
    echo "✅ Deployment complete"
}

# Health check
health_check() {
    echo ""
    echo "Running health checks..."
    
    sleep 30  # Wait for services to be ready
    
    project_id=$(railway variables | grep RAILWAY_PROJECT_ID | cut -d= -f2)
    
    services=("credit-service" "billing-service" "ppob-core")
    
    for service in "${services[@]}"; do
        url="https://$project_id-$service.up.railway.app/health"
        echo -n "Checking $service... "
        
        if curl -s -f $url > /dev/null; then
            echo "✅ OK"
        else
            echo "❌ FAILED"
            return 1
        fi
    done
    
    echo "✅ All services healthy"
}

# Main flow
main() {
    check_prerequisites
    login_railway
    create_project
    create_services
    run_migrations
    deploy
    health_check
    
    echo ""
    echo "==================================="
    echo "✅ Deployment Successful!"
    echo "==================================="
    echo ""
    echo "Next steps:"
    echo "1. View dashboard: railway dashboard"
    echo "2. View logs: railway logs"
    echo "3. Monitor services: railway status"
    echo ""
    echo "📚 Documentation:"
    echo "   Local dev: make up"
    echo "   API docs: docs/openapi.yaml"
    echo "   Deployment: docs/RAILWAY.md"
    echo ""
}

# Run main if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
