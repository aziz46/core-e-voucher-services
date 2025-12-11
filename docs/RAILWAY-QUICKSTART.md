# Quick Start Guide untuk Railway Deployment

## 🚀 Railway Deployment - 5 Menit Setup

Railway adalah platform PaaS yang paling mudah untuk deploy aplikasi Go. Berikut step-by-stepnya:

### Prerequisites
- GitHub account & repo (push code)
- Railway account (daftar gratis di https://railway.app)
- ✅ Semua file sudah siap di repo!

### Option 1: Automated Setup (Recommended)

```bash
# 1. Make script executable
chmod +x scripts/deploy-railway.sh

# 2. Run deployment script
./scripts/deploy-railway.sh

# Script akan:
# ✅ Check prerequisites
# ✅ Login ke Railway
# ✅ Create project
# ✅ Setup PostgreSQL
# ✅ Deploy 3 services
# ✅ Run migrations
# ✅ Health check
```

### Option 2: Manual Setup via Dashboard

1. **Login ke Railway**
   - Buka https://railway.app/dashboard
   - Click "Login with GitHub"

2. **Create New Project**
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Authorize GitHub
   - Select: `aziz46/core-e-voucher-services`
   - Click "Deploy"

3. **Add PostgreSQL**
   - Click "Add Service"
   - Select "PostgreSQL"
   - Click "Deploy"

4. **Add Credit Service**
   - Click "Add Service" → "GitHub Repo"
   - Select main branch
   - Set Dockerfile: `./cmd/credit-service/Dockerfile`
   - Click "Deploy"
   - Wait for build & deploy

5. **Add Billing Service**
   - Repeat step 4
   - Set Dockerfile: `./cmd/billing-service/Dockerfile`

6. **Add PPOB Core**
   - Repeat step 4
   - Set Dockerfile: `./cmd/ppob-core/Dockerfile`
   - Set "Public Port" to expose API

7. **Link Database to Services**
   - For each service:
     - Click service → "Variables"
     - Click "Link" → select PostgreSQL
     - Railway auto-configures DATABASE_URL

8. **Run Migrations**
   - Click PostgreSQL service
   - Tab "Data"
   - Copy-paste dari migrations/001_init.sql
   - Execute
   - Repeat untuk 002_seed.sql

9. **Verify Deployment**
   ```bash
   # Get service URLs from Railway dashboard
   curl https://[service-id].up.railway.app/health
   ```

### Quick Reference Commands

```bash
# Using Railway CLI
railway login
railway status
railway logs --service credit-service
railway variables
railway up --service ppob-core
```

### Environment Variables di Railway

Railway automatically provides:
```
DATABASE_URL=postgres://...
PGHOST=...
PGPORT=...
PGUSER=...
PGPASSWORD=...
PGDATABASE=e_voucher
```

Anda tinggal set di Railway dashboard:
```
SERVER_PORT=8080
LOG_LEVEL=info
ENVIRONMENT=production
```

### Access Your Services

Setelah deploy berhasil:
```
PPOB Core: https://[project-id]-ppob-core.up.railway.app
Credit Service: https://[project-id]-credit-service.up.railway.app
Billing Service: https://[project-id]-billing-service.up.railway.app
PostgreSQL: Managed by Railway (no public URL)
```

### Test API

```bash
# Health check
curl https://[project-id]-ppob-core.up.railway.app/health

# Create transaction
curl -X POST https://[project-id]-ppob-core.up.railway.app/v1/tenant_001/transactions \
  -H "X-API-Key: sk_live_abc123xyz789" \
  -H "Content-Type: application/json" \
  -d '{
    "product_code": "PLN_PREPAID",
    "customer_no": "081234567890",
    "amount": 50000,
    "partner_id": "partner_001"
  }'
```

### Auto Deployment from GitHub

Railway automatically deploys setiap kali push ke main branch:

```bash
git add .
git commit -m "feat: add new endpoint"
git push

# Railway akan:
# 1. Pull latest code
# 2. Build Docker image
# 3. Run tests
# 4. Deploy to services
# 5. Run health checks
```

### Monitor & Debug

**Via Railway Dashboard**
- Click service → "Logs" untuk real-time logs
- "Metrics" untuk CPU/Memory/Bandwidth
- "Deployments" untuk deployment history

**Via Railway CLI**
```bash
railway logs --service ppob-core --follow
railway status
```

### Troubleshooting

**Service not starting?**
```bash
# Check logs
railway logs --service credit-service --raw
```

**Database connection error?**
```bash
# Verify DATABASE_URL
railway variables | grep DATABASE

# Check database status
railway status postgres
```

**Port issues?**
```bash
# Railway auto-assigns ports
# Check via: railway status
```

### Cost

Railway pricing:
- **Free tier**: $0 for testing
- **Minimum**: $5/month (includes some free credits)
- **PostgreSQL**: $2-5/month
- **Services**: ~$0.15/GB/hour

Total estimasi: $5-15/month untuk MVP.

### Next Steps

1. ✅ Deployed ke Railway
2. Setup custom domain (optional)
3. Enable automated backups
4. Setup error tracking (Sentry)
5. Configure rate limiting
6. Integrate with real payment provider

### Useful Links

- Railway Docs: https://docs.railway.app
- Railway CLI: https://docs.railway.app/cli
- Deployment Guide: docs/RAILWAY.md
- Full README: docs/README.md

---

**Selesai!** Aplikasi sudah live di Railway. Akses dari mana saja dengan HTTPS!
