# Deployment ke Railway

Panduan step-by-step untuk deploy Core E-Voucher PPOB ke Railway.

## 📋 Prerequisites

1. **Railway Account** - Daftar di https://railway.app
2. **GitHub Account** - Repo harus di GitHub
3. **Railway CLI** (optional) - https://docs.railway.app/cli/installation

## 🚀 Deployment Steps

### Step 1: Push Repository ke GitHub

```bash
git init
git add .
git commit -m "Initial commit: PPOB platform setup"
git branch -M main
git remote add origin https://github.com/aziz46/core-e-voucher-services.git
git push -u origin main
```

### Step 2: Login ke Railway

Buka https://railway.app/dashboard dan login dengan GitHub account.

### Step 3: Create New Project

1. Klik **"Create a New Project"**
2. Pilih **"Deploy from GitHub repo"**
3. Connect GitHub account (jika belum)
4. Select repository: `core-e-voucher-services`
5. Klik **"Deploy"**

### Step 4: Create PostgreSQL Database

1. Di project dashboard, klik **"Add Service"**
2. Pilih **"PostgreSQL"** dari list
3. Railway akan auto-generate password
4. Klik **"Deploy"**
5. Copy environment variables yang digenerate (akan otomatis linked)

### Step 5: Configure Services

Setiap service perlu dikonfigurasi terpisah.

#### **Service 1: Credit Service**

1. Klik **"Add Service"** → **"GitHub Repo"**
2. Select branch: `main`
3. Setting konfigurasi:
   - **Name**: `credit-service`
   - **Dockerfile**: `./cmd/credit-service/Dockerfile`
   - **Root Directory**: `.` (root)
   - **Port**: `8080`
   - **Memory**: 512MB
   - **CPU**: 0.5

4. Environment Variables:
   ```
   SERVER_PORT=8080
   SERVER_HOST=0.0.0.0
   LOG_LEVEL=info
   ```

5. Klik **"Deploy"**

#### **Service 2: Billing Service**

Ulangi step yang sama dengan:
- **Name**: `billing-service`
- **Dockerfile**: `./cmd/billing-service/Dockerfile`
- **Port**: `8080`

#### **Service 3: PPOB Core**

Ulangi step yang sama dengan:
- **Name**: `ppob-core`
- **Dockerfile**: `./cmd/ppob-core/Dockerfile`
- **Port**: `8080`

### Step 6: Link Database ke Services

1. Untuk setiap service (credit, billing, ppob):
   - Buka service setting
   - Tab **"Variables"**
   - Klik **"Link"** database yang sudah dibuat
   - Railway auto-populate: `DATABASE_URL`, `DB_HOST`, dll

2. Jika perlu custom variables, tambahkan:
   ```
   DB_NAME=e_voucher
   DB_SSLMODE=disable
   LOG_LEVEL=info
   ```

### Step 7: Run Migrations

1. Buka PostgreSQL service
2. Tab **"Data"**
3. Klik **"Execute"** atau buka terminal
4. Copy-paste content dari `migrations/001_init.sql`
5. Run migration
6. Repeat untuk `migrations/002_seed.sql`

Atau gunakan **Railway CLI**:
```bash
railway database exec < migrations/001_init.sql
railway database exec < migrations/002_seed.sql
```

### Step 8: Verify Deployment

Check service health:

```bash
# Credit Service
curl https://credit-service-[project-id].railway.app/health

# Billing Service
curl https://billing-service-[project-id].railway.app/health

# PPOB Core
curl https://ppob-core-[project-id].railway.app/health
```

Expected response:
```json
{"status": "ok"}
```

---

## 📋 Railway Configuration Files

### `railway.json` (Optional)

Buat file `railway.json` di root untuk standardisasi:

```json
{
  "services": [
    {
      "name": "postgres",
      "image": "postgres:13"
    },
    {
      "name": "credit-service",
      "dockerfile": "./cmd/credit-service/Dockerfile",
      "port": 8080
    },
    {
      "name": "billing-service",
      "dockerfile": "./cmd/billing-service/Dockerfile",
      "port": 8080
    },
    {
      "name": "ppob-core",
      "dockerfile": "./cmd/ppob-core/Dockerfile",
      "port": 8080
    }
  ]
}
```

### `Procfile` (Alternative)

Jika ingin menggunakan built-in Go detection:

```
web: ./cmd/ppob-core/ppob-core
credit: ./cmd/credit-service/credit-service
billing: ./cmd/billing-service/billing-service
```

---

## 🔐 Environment Variables di Railway

Railway auto-provide database credentials. Mapping:

| Railway Variable | Go Code | Usage |
|-----------------|---------|-------|
| `DATABASE_URL` | Manual parsing | Full connection string |
| `PGHOST` | `cfg.Database.Host` | Database host |
| `PGPORT` | `cfg.Database.Port` | Database port |
| `PGUSER` | `cfg.Database.User` | Database user |
| `PGPASSWORD` | `cfg.Database.Password` | Database password |
| `PGDATABASE` | `cfg.Database.DBName` | Database name |

Update code untuk Railway:

```go
// pkg/config/config.go
func LoadConfig() *Config {
    // ...
    
    // For Railway, prefer DATABASE_URL
    if os.Getenv("DATABASE_URL") != "" {
        // Parse DATABASE_URL dan extract components
        dsn := os.Getenv("DATABASE_URL")
        // Parse dan assign ke cfg.Database
    }
    
    // Fallback to individual env vars
    cfg.Database.Host = viper.GetString("PGHOST")
    cfg.Database.Port = viper.GetInt("PGPORT")
    // ...
    
    return cfg
}
```

---

## 🔗 Service Communication di Railway

Dalam Railway, services berkomunikasi via URL publik (tidak ada private network).

Update PPOB Core untuk call services:

```go
// internal/ppob/handler/transaction.go

// Di Railway, gunakan publik URL
creditURL := os.Getenv("CREDIT_SERVICE_URL")
if creditURL == "" {
    creditURL = "http://credit-service:8080" // Local fallback
}
```

Set environment variables di Railway:
```
CREDIT_SERVICE_URL=https://credit-service-[project].railway.app
BILLING_SERVICE_URL=https://billing-service-[project].railway.app
```

---

## 💾 Persistent Data di Railway

Railway provides:
- **PostgreSQL** - Auto-managed, persistent
- **Redis** (optional) - Available as add-on

Untuk MinIO (file storage), ada 2 opsi:

**Opsi 1: Gunakan Railway's S3-compatible storage**
```
Belum tersedia, tapi bisa gunakan AWS S3 atau DigitalOcean Spaces
```

**Opsi 2: Disable MinIO untuk MVP**
```go
// Simpan receipt path langsung di database untuk MVP
// Real implementation: integrate S3/Spaces later
```

---

## 🔄 CI/CD dengan Railway

Railway auto-deploy saat push ke main branch:

1. Push code ke GitHub:
```bash
git add .
git commit -m "fix: update config for railway"
git push
```

2. Railway otomatis:
   - Detect Dockerfile
   - Build image
   - Deploy ke staging/production
   - Run health checks

Monitor deployment di Railway dashboard.

---

## 📊 Monitoring & Logs

### Via Railway Dashboard
1. Login ke https://railway.app/dashboard
2. Select project
3. Select service
4. Tab **"Logs"** untuk view real-time logs
5. Tab **"Metrics"** untuk monitoring CPU/Memory

### Via Railway CLI
```bash
# Install CLI
npm install -g @railway/cli

# Login
railway login

# View logs
railway logs --service credit-service
railway logs --service billing-service
railway logs --service ppob-core

# View variables
railway variables
```

---

## 💰 Pricing & Cost Optimization

Railway pricing model:
- **$5/month** minimum spend (includes some free tier)
- **PostgreSQL**: $2-5/month depending on size
- **Services**: $0.15 per GB/hour usage

**Cost optimization:**
- Scale down services ketika tidak ada traffic
- Use PostgreSQL shared tier untuk MVP
- Monitor metrics dan adjust resource allocation

---

## 🆘 Troubleshooting

### Service not starting
```bash
# Check build logs
railway logs --service credit-service --raw

# Common issues:
# 1. Missing DATABASE_URL
# 2. Port already in use (Railway assigns automatically)
# 3. Build failed
```

### Database connection error
```bash
# Verify database is running
railway status

# Check database credentials
railway variables | grep PG

# Test connection
psql $DATABASE_URL
```

### Service timeout
```bash
# Increase timeout di Railway dashboard
# Service Settings → Network → Health Check Interval
```

---

## 🎯 Deployment Checklist

- [ ] GitHub repo created dan pushed
- [ ] Railway project created
- [ ] PostgreSQL database provisioned
- [ ] 3 services deployed (credit, billing, ppob-core)
- [ ] Database migrations applied
- [ ] Health checks passing
- [ ] Inter-service communication working
- [ ] Environment variables configured
- [ ] Logs monitoring setup
- [ ] Custom domain configured (optional)

---

## 🌐 Custom Domain (Optional)

1. Railway dashboard → Project → Settings
2. Domains → Add custom domain
3. Update DNS records di domain registrar
4. Wait untuk SSL certificate provisioning

---

## 📚 Useful Links

- **Railway Docs**: https://docs.railway.app
- **Railway CLI**: https://docs.railway.app/cli
- **Pricing**: https://railway.app/pricing
- **Go on Railway**: https://docs.railway.app/languages/go

---

## Next Steps

Setelah deploy berhasil:
1. Setup API monitoring (Sentry, etc)
2. Setup automated backups untuk PostgreSQL
3. Configure rate limiting di reverse proxy
4. Setup alerts untuk down-time
5. Integrate payment gateway (real provider)

---

Sudah siap deploy ke Railway! Cukup follow steps di atas. Ada yang perlu klarifikasi?
