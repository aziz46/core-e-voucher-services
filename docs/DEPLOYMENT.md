# Deployment Guides

Dokumentasi deployment untuk berbagai platform.

## 📋 Supported Platforms

1. **[Railway](./RAILWAY.md)** - Recommended (easiest setup)
2. **[Docker Compose](../docker-compose.yml)** - Local development & testing
3. **Kubernetes** - Enterprise (TBD)
4. **Manual VPS** - Self-hosted (TBD)

## 🚀 Quick Start by Platform

### Railway (Recommended for MVP)
```bash
# See: ./RAILWAY.md
```
- ✅ Zero DevOps knowledge needed
- ✅ Auto-scaling & monitoring
- ✅ 5 minute setup
- ✅ $5/month minimum

### Docker Compose (Local / Self-hosted)
```bash
make up
# See: ../docker-compose.yml & ../Makefile
```
- ✅ Works anywhere Docker runs
- ✅ Full control
- ❌ Manual scaling
- 💰 Free (your own infrastructure)

### Kubernetes (Enterprise)
```bash
kubectl apply -f deploy/k8s/
```
- ✅ Production-grade scaling
- ✅ Auto-healing & load balancing
- ❌ Complex setup
- 💰 Varies

---

## 📊 Platform Comparison

| Aspect | Railway | Docker Compose | Kubernetes |
|--------|---------|---|---|
| Setup Time | 5 min | 5 min | 1-2 hour |
| Learning Curve | 0 | Low | High |
| Cost | $5-20/mo | Free* | $50-300+/mo |
| Auto-scaling | ✅ | ❌ | ✅ |
| Auto-SSL | ✅ | ❌ | ✅ (with ingress) |
| Monitoring | ✅ | ❌ | Need add-on |
| CI/CD Integration | ✅ | ❌ | ✅ |

*requires your own server

---

## 🔄 Deployment Pipeline

```
┌─────────────────┐
│  Git Push       │
│  (main branch)  │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  GitHub Actions CI                  │
│  - Build & Test                     │
│  - Code Quality Checks              │
│  - Build Docker Images              │
│  - Push to Registry                 │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Deploy to Staging                  │
│  - Run Migrations                   │
│  - Deploy Services                  │
│  - Health Checks                    │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Smoke Tests                        │
│  - API Tests                        │
│  - Integration Tests                │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Manual Approval (for production)   │
└────────┬────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Deploy to Production               │
│  - Blue-Green Deployment            │
│  - Database Migrations              │
│  - Health Verification              │
└─────────────────────────────────────┘
```

---

## 📝 Environment Configuration

### Development (.env.example)
```env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
```

### Staging (Railway)
```env
SERVER_PORT=8080
DB_HOST=<railway-postgres>
DB_PORT=5432
ENVIRONMENT=staging
LOG_LEVEL=debug
```

### Production (Railway / Self-hosted)
```env
SERVER_PORT=8080
DB_HOST=<production-postgres>
DB_PORT=5432
ENVIRONMENT=production
LOG_LEVEL=warn
```

---

## 🔐 Secrets Management

### Local Development
- Use `.env` file (in `.gitignore`)
- Never commit secrets

### Railway
- Use Railway Variables dashboard
- Automatically encrypted at rest
- Linked to services

### GitHub Actions (for CI/CD)
- Use GitHub Secrets
- Reference as: `${{ secrets.SECRET_NAME }}`

---

## ✅ Pre-deployment Checklist

- [ ] All tests passing
- [ ] No TypeErrors or linting errors
- [ ] Database migrations reviewed
- [ ] Environment variables configured
- [ ] Docker images build successfully
- [ ] Health check endpoints working
- [ ] Documentation updated
- [ ] Changelog updated

---

## 📞 Support

- **Railway Docs**: https://docs.railway.app
- **GitHub Actions**: https://docs.github.com/en/actions
- **Docker Compose**: https://docs.docker.com/compose/

---

See specific deployment guides:
- [Railway Deployment](./RAILWAY.md)
- [Local Development](../README.md)
