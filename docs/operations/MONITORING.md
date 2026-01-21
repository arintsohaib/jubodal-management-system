# MONITORING.md

## Purpose

Defines **monitoring, observability, and alerting strategy** for BJDMS running on arint.win.

**Goal**: Detect issues before users experience them, maintain 99.5% uptime.

---

## Monitor

ing Stack

**Recommended Tools** (lightweight, Docker-compatible):
- **Prometheus**: Metrics collection
- **Grafana**: Visualization dashboards
- **Loki**: Log aggregation
- **Alertmanager**: Alert routing (SMS/email)

**Alternative** (simpler): Uptime monitoring services (UptimeRobot, Healthchecks.io)

---

## Key Metrics

### Application Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| API Response Time (p95) | <200ms | >500ms |
| Error Rate | <0.1% | >1% |
| Request Rate | - | >1000 req/s (DDoS?) |
| Active Users (concurrent) | - | Monitor trend |

### Infrastructure Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| CPU Usage | <70% | >85% for 5 min |
| Memory Usage | <80% | >90% for 5 min |
| Disk Usage | <80% | >90% |
| Network I/O | - | Monitor spikes |

### Database Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| Connection Pool Usage | <80% | >90% |
| Query Duration (p95) | <100ms | >500ms |
| Slow Queries | 0 | >5 per hour |
| Deadlocks | 0 | >0 |

### Redis Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| Memory Usage | <2GB | >3GB |
| Hit Rate | >90% | <80% |
| Connected Clients | <100 | >200 |

---

## Health Checks

### Application Health Endpoint

**URL**: `https://arint.win/api/v1/health`

**Response**:
```json
{
  "status": "healthy",
  "timestamp": "2026-01-21T14:55:00Z",
  "services": {
    "database": "connected",
    "redis": "connected",
    "opensearch": "connected",
    "storage": "accessible"
  },
  "uptime": 345600
}
```

**Check Frequency**: Every 60 seconds

**Alert if**: Status not "healthy" for 2 consecutive checks

---

### Database Health

```sql
SELECT 1;  -- Simple connectivity check
SELECT COUNT(*) FROM users WHERE is_active = true;  -- Data accessibility
```

---

### Redis Health

```bash
redis-cli PING  # Expected: PONG
```

---

## Logging Strategy

### Log Levels

| Level | Usage | Examples |
|-------|-------|----------|
| ERROR | System failures, critical issues | Database connection lost, API crash |
| WARN | Potential issues, degraded performance | Slow query (>1s), high memory usage |
| INFO | Normal operations | User login, committee created |
| DEBUG | Development/troubleshooting | SQL queries, function calls |

**Production**: Log INFO and above only.

---

### Structured Logging (JSON)

```json
{
  "level": "info",
  "timestamp": "2026-01-21T14:55:00Z",
  "message": "User logged in",
  "user_id": "uuid",
  "ip_address": "1.2.3.4",
  "user_agent": "Mozilla/5.0...",
  "request_id": "req-uuid"
}
```

**Benefits**: Easily parseable, filterable, searchable.

---

### Log Retention

- **Application Logs**: 30 days
- **Access Logs**: 90 days
- **Error Logs**: 1 year
- **Audit Logs**: 5 years (separate database table)

---

## Alerting Rules

### Critical Alerts (Immediate SMS/Call)

1. **Database Down**: Cannot connect to PostgreSQL
2. **API Down**: Health check fails for 5 minutes
3. **Disk Full**: <5% disk space remaining
4. **Security Incident**: Multiple failed Super Admin login attempts

---

### High Priority (SMS/Email within 15 min)

1. **High Error Rate**: >1% of requests fail
2. **Slow Response**: p95 latency >1s
3. **Memory Critical**: >95% memory usage
4. **SSL Certificate Expiring**: <7 days until expiry

---

###Medium Priority (Email within 1 hour)

1. **Elevated Error Rate**: >0.5% errors
2. **High CPU**: >80% for 15 minutes
3. **Slow Queries**: >10 queries >1s per hour
4. **Failed Backup**: Daily backup did not complete

---

## Dashboards

### 1. System Overview

- Current active users
- Requests per minute
- Error rate %
- Response time (p50, p95, p99)
- CPU/Memory/Disk usage

### 2. Application Performance

- API endpoint response times
- Database query performance
- Redis cache hit rate
- File upload/download metrics

### 3. User Activity

- Logins per hour
- Activities submitted per day
- Complaints received per day
- Donations recorded per day

### 4. Security

- Failed login attempts
- Account lockouts
- Unauthorized access attempts
- IP blocklist size

---

## Prometheus Configuration Example

```yaml
scrape_configs:
  - job_name: 'api'
    static_configs:
      - targets: ['api:3000']
    metrics_path: '/metrics'
    scrape_interval: 15s

  - job_name: 'analytics'
    static_configs:
      - targets: ['analytics:8000']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

---

## Grafana Dashboard Panels

1. **Active Users** (Gauge)
2. **Request Rate** (Graph, requests/sec)
3. **Error Rate** (Graph, %)
4. **Response Time Distribution** (Heatmap)
5. **Database Connections** (Gauge)
6. **CPU Usage** (Graph, %)
7. **Memory Usage** (Graph, %)
8. **Disk Space** (Gauge, GB free)

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
