# CI_CD_PIPELINE.md

## Purpose

Defines **Continuous Integration and Continuous Deployment (CI/CD) pipeline** for BJDMS.

**Critical Constraint**: All builds and deployments execute on **remote server (grayhawks.com)** only. No local builds.

---

## Pipeline Architecture

```
Git Push → GitHub/GitLab → Webhook → grayhawks.com → Build → Test → Deploy
```

**Trigger**: Push to `main` branch or Pull Request

---

## CI Pipeline (Continuous Integration)

### On Every Pull Request

**Steps**:
1. **Checkout Code**
2. **Lint** - Check code style
3. **Unit Tests** - Run all unit tests
4. **Integration Tests** - Run with test database
5. **Security Scan** - gosec, npm audit, bandit
6. **Coverage Check** - Fail if <70%
7. **Build** - Verify code compiles

**GitHub Actions Example** (.github/workflows/ci.yml):
```yaml
name: CI Pipeline

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run golangci-lint
        run: golangci-lint run ./...
  
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: bjdms_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
      
      redis:
        image: redis:7
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run migrations
        run: |
          go run cmd/migrate/main.go up
        env:
          DATABASE_URL: postgresql://postgres:test@localhost:5432/bjdms_test
      
      - name: Run tests
        run: go test ./... -v -coverprofile=coverage.out
        env:
          DATABASE_URL: postgresql://postgres:test@localhost:5432/bjdms_test
          REDIS_URL: redis://localhost:6379
      
      - name: Check coverage
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$COVERAGE < 70" | bc -l) )); then
            echo "Coverage $COVERAGE% is below 70%"
            exit 1
          fi
  
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - name: Run gosec
        run: |
          go install github.com/securego/gosec/v2/cmd/gosec@latest
          gosec ./...
      - name: Run npm audit (if Node.js frontend)
        run: npm audit --audit-level=high
        working-directory: ./frontend
```

---

## CD Pipeline (Continuous Deployment)

### Deployment Trigger

**Manual Approval** (for production):
- PR merged to `main` → Manual deploy button
- Tagged release (v1.0.0) → Auto-deploy to production

**Auto-Deploy** (for staging):
- Push to `main` → Auto-deploy to staging

---

### Remote Deployment Script

**Webhook on grayhawks.com**:

**/opt/bjdms/scripts/deploy.sh**:
```bash
#!/bin/bash
set -e

echo "Starting deployment..."

# Navigate to project
cd /opt/bjdms

# Pull latest code
git pull origin main

# Backup database
./scripts/backup_database.sh

# Run migrations
docker-compose run --rm api go run cmd/migrate/main.go up

# Build new images
docker-compose build --no-cache

# Zero-downtime deploy (blue-green)
# Start green environment
docker-compose -f docker-compose.green.yml up -d

# Health check green
for i in {1..30}; do
  if curl -f http://localhost:3001/health; then
    echo "Green environment healthy"
    break
  fi
  sleep 2
done

# Switch NGINX to green
sed -i 's/localhost:3000/localhost:3001/' /etc/nginx/sites-available/bjdms
nginx -t && systemctl reload nginx

# Stop blue
docker-compose down

# Rename green to blue for next deployment
mv docker-compose.green.yml docker-compose.blue.yml
mv docker-compose.yml docker-compose.green.yml
mv docker-compose.blue.yml docker-compose.yml

echo "Deployment complete!"
```

**Make executable**:
```bash
chmod +x /opt/bjdms/scripts/deploy.sh
```

---

### GitHub Actions Deployment

**.github/workflows/deploy.yml**:
```yaml
name: Deploy to Production

on:
  workflow_dispatch:  # Manual trigger
  push:
    tags:
      - 'v*'  # Auto-deploy on version tags

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to grayhawks.com
        uses: appleboy/ssh-action@master
        with:
          host: grayhawks.com
          username: root
          port: 25920
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /opt/bjdms
            ./scripts/deploy.sh
      
      - name: Verify deployment
        run: |
          sleep 30
          curl -f https://grayhawks.com/api/v1/health || exit 1
      
      - name: Notify team
        if: success()
        run: echo "Deployment successful!"
      
      - name: Rollback on failure
        if: failure()
        uses: appleboy/ssh-action@master
        with:
          host: grayhawks.com
          username: root
          port: 25920
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /opt/bjdms
            git reset --hard HEAD~1
            docker-compose up -d
```

---

## Pipeline Secrets

**GitHub Secrets** (Settings → Secrets):
- `SSH_PRIVATE_KEY`: SSH key for grayhawks.com access
- `DATABASE_URL`: Production database connection string (for migrations)
- `JWT_SECRET`: JWT signing key
- `SMS_API_KEY`: SMS gateway credentials

**Security**: Never commit secrets to code. Use environment variables or secrets management.

---

## Build Artifacts

**Docker Images**:
- Built on remote server (grayhawks.com)
- Tagged with git commit SHA: `bjdms-api:abc1234`
- Stored in local Docker registry or Docker Hub (optional)

---

## Rollback Automation

**Automatic Rollback** (if health check fails):
```bash
# In deploy.sh, after health check
if ! curl -f http://localhost:3001/health; then
  echo "Health check failed, rolling back..."
  docker-compose -f docker-compose.green.yml down
  docker-compose up -d  # Blue environment still running
  exit 1
fi
```

---

## Deployment Notifications

**Slack/Discord/Email**:
```yaml
- name: Notify Slack
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    text: 'Deployment to production: ${{ job.status }}'
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

---

## Environment-Specific Pipelines

### Staging Pipeline

**Trigger**: Every push to `main`

**Steps**:
1. Deploy to `/opt/bjdms-staging`
2. Run smoke tests
3. Notify QA team

### Production Pipeline

**Trigger**: Manual approval or version tag

**Steps**:
1. Deploy to `/opt/bjdms`
2. Run smoke tests
3. Monitor for 30 minutes
4. Send notification to ops team

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
