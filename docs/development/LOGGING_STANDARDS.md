# LOGGING_STANDARDS.md

## Purpose

Defines **logging standards and best practices** for BJDMS to enable effective debugging, monitoring, and auditing.

---

## Logging Principles

1. **Structured Logging** - Use JSON format for machine readability
2. **Appropriate Levels** - Use correct log level (DEBUG, INFO, WARN, ERROR)
3. **Contextual Information** - Include request ID, user ID, relevant IDs
4. **No Sensitive Data** - Never log passwords, tokens, PII
5. **Performance Aware** - Avoid excessive logging in hot paths

---

## Log Levels

| Level | Usage | Production | Examples |
|-------|-------|------------|----------|
| DEBUG | Detailed troubleshooting info | NO | SQL queries, function entry/exit |
| INFO | Normal operations | YES | User login, committee created |
| WARN | Potential issues, degraded performance | YES | Slow query (>1s), high memory |
| ERROR | Failures, exceptions | YES | Database error, API call failed |
| FATAL | System cannot continue | YES | Database unavailable, config missing |

**Production Log Level**: INFO (DEBUG only in development/troubleshooting)

---

## Structured Logging Format

**JSON Format**:
```json
{
  "level": "info",
  "timestamp": "2026-01-21T15:00:00Z",
  "message": "User logged in",
  "user_id": "uuid",
  "phone": "+880171***5678",  // Masked
  "ip_address": "1.2.3.4",
  "request_id": "req-uuid",
  "duration_ms": 45
}
```

**Benefits**:
- Easily parseable by log aggregators (Loki, ELK)
- Filterable by fields
- Machine-readable

---

## Language-Specific Libraries

### Go

**Library**: `logrus` or `zap` (high-performance)

**Setup** (logrus):
```go
import (
    log "github.com/sirupsen/logrus"
)

func init() {
    log.SetFormatter(&log.JSONFormatter{})
    log.SetLevel(log.InfoLevel)
}

// Usage
log.WithFields(log.Fields{
    "user_id": userID,
    "request_id": requestID,
    "ip_address": r.RemoteAddr,
}).Info("User logged in")
```

---

### Python

**Library**: `structlog` or built-in `logging`

**Setup**:
```python
import structlog

structlog.configure(
    processors=[
        structlog.stdlib.add_log_level,
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.JSONRenderer()
    ]
)

logger = structlog.get_logger()

# Usage
logger.info("user_logged_in", user_id=user_id, request_id=request_id)
```

---

### TypeScript/JavaScript

**Library**: `winston` or `pino`

**Setup** (winston):
```typescript
import winston from 'winston';

const logger = winston.createLogger({
  level: 'info',
  format: winston.format.json(),
  transports: [
    new winston.transports.File({ filename: 'error.log', level: 'error' }),
    new winston.transports.File({ filename: 'combined.log' }),
  ],
});

// Usage
logger.info('User logged in', {
  user_id: userId,
  request_id: requestId,
  ip_address: req.ip
});
```

---

## What to Log

### Authentication Events

```go
log.WithFields(log.Fields{
    "event": "login_success",
    "user_id": user.ID,
    "phone": maskPhone(user.Phone),  // +880171***5678
    "ip_address": r.RemoteAddr,
    "request_id": requestID,
}).Info("User logged in")

log.WithFields(log.Fields{
    "event": "login_failed",
    "phone": maskPhone(phone),
    "reason": "invalid_password",
    "ip_address": r.RemoteAddr,
}).Warn("Login attempt failed")
```

---

### API Requests

**Middleware** (Go):
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        requestID := generateRequestID()
        
        // Store request ID in context
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        r = r.WithContext(ctx)
        
        next.ServeHTTP(w, r)
        
        duration := time.Since(start)
        
        log.WithFields(log.Fields{
            "method": r.Method,
            "path": r.URL.Path,
            "duration_ms": duration.Milliseconds(),
            "request_id": requestID,
            "ip_address": r.RemoteAddr,
            "user_agent": r.UserAgent(),
        }).Info("API request completed")
    })
}
```

---

### Database Operations

```go
// Slow queries
start := time.Now()
rows, err := db.Query("SELECT * FROM users WHERE ...")
duration := time.Since(start)

if duration > time.Second {
    log.WithFields(log.Fields{
        "query": "SELECT * FROM users WHERE ...",
        "duration_ms": duration.Milliseconds(),
    }).Warn("Slow database query detected")
}
```

---

### Business Logic Events

```go
log.WithFields(log.Fields{
    "event": "committee_created",
    "committee_id": committee.ID,
    "jurisdiction_id": committee.JurisdictionID,
    "committee_type": committee.Type,
    "created_by": userID,
}).Info("Committee created")

log.WithFields(log.Fields{
    "event": "complaint_submitted",
    "complaint_id": complaint.ID,
    "is_anonymous": complaint.IsAnonymous,
    "target_user_id": complaint.TargetUserID,
}).Info("Complaint submitted")
```

---

### Errors

```go
log.WithFields(log.Fields{
    "error": err.Error(),
    "request_id": requestID,
    "user_id": userID,
    "context": "creating committee",
}).Error("Failed to create committee")

// Include stack trace for unexpected errors
log.WithFields(log.Fields{
    "error": err.Error(),
    "stack": string(debug.Stack()),
}).Fatal("Unrecoverable error")
```

---

## What NOT to Log

### Sensitive Data

**NEVER log**:
- Passwords (plain or hashed)
- JWT tokens
- API keys
- Full credit card numbers
- National ID numbers (NID)
- Full phone numbers (mask: +880171***5678)
- Personal addresses

**Masking Function** (Go):
```go
func maskPhone(phone string) string {
    if len(phone) < 8 {
        return "***"
    }
    return phone[:6] + "***" + phone[len(phone)-4:]
}

// +8801712345678 → +880171***5678
```

---

### High-Frequency Operations

**Avoid logging**:
- Every GET request (only log errors or slow requests)
- Cache hits (too noisy)
- Heartbeat/health checks

---

## Log Retention

| Log Type | Retention | Storage |
|----------|-----------|---------|
| Application logs | 30 days | Local disk + S3 |
| Error logs | 90 days | S3 |
| Access logs (NGINX) | 90 days | Local disk |
| Audit logs | 5 years | Database table |

**Rotation**:
```bash
# logrotate config
/opt/bjdms/logs/*.log {
    daily
    rotate 30
    compress
    missingok
    notifempty
    create 0644 root root
}
```

---

## Centralized Logging

**Tools**: Loki (lightweight), ELK Stack (full-featured)

**Log Shipping** (Docker setup):
```yaml
# docker-compose.yml
services:
  api:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
        labels: "service=api"
  
  loki:
    image: grafana/loki:latest
    ports:
      - "3100:3100"
```

**Query Logs** (Loki):
```
{service="api"} |= "error" | json
```

---

## Performance Considerations

### Avoid String Formatting in Hot Paths

**Bad** (always evaluates):
```go
log.Debug(fmt.Sprintf("Processing user %s with data %v", userID, userData))
```

**Good** (only if DEBUG enabled):
```go
if log.GetLevel() >= log.DebugLevel {
    log.Debug("Processing user", "user_id", userID, "data", userData)
}
```

---

### Batch Logging (for high-volume events)

```go
type LogBuffer struct {
    entries []LogEntry
    mu      sync.Mutex
}

func (b *LogBuffer) Add(entry LogEntry) {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    b.entries = append(b.entries, entry)
    
    if len(b.entries) >= 100 {
        b.flush()
    }
}

func (b *LogBuffer) flush() {
    for _, entry := range b.entries {
        log.Info(entry)
    }
    b.entries = nil
}
```

---

## Debugging with Logs

### Request ID Tracing

**Generate unique request ID**:
```go
func generateRequestID() string {
    return uuid.New().String()
}
```

**Include in all logs for that request**:
```go
logger := log.WithField("request_id", requestID)
logger.Info("Starting request")
// ... later in code
logger.Error("Request failed")
```

**Filter logs by request ID**:
```bash
grep "req-abc-123" /opt/bjdms/logs/api.log
```

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
