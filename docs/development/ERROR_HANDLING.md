# ERROR_HANDLING.md

## Purpose

Defines **error handling patterns and standards** for BJDMS to ensure consistent, debuggable, and user-friendly error management.

---

## Error Handling Principles

1. **Never Swallow Errors** - Always handle or propagate
2. **Add Context** - Wrap errors with meaningful information
3. **Log Appropriately** - Log errors at correct level
4. **User-Friendly Messages** - Don't expose internal details to users
5. **Fail Fast** - Return error immediately, don't continue with invalid state

---

## Error Types

### 1. Validation Errors (400 Bad Request)

**When**: User input is invalid

**Examples**:
- Invalid phone format
- Missing required field
- Age < 18

**Response**:
```json
{
  "error": "validation_error",
  "message": "Invalid input data",
  "details": [
    {
      "field": "phone",
      "error": "Phone must match Bangladesh format (+880XXXXXXXXXX)"
    },
    {
      "field": "dateOfBirth",
      "error": "Applicant must be at least 18 years old"
    }
  ]
}
```

---

### 2. Authentication Errors (401 Unauthorized)

**When**: User not authenticated or token invalid

**Examples**:
- Token expired
- Invalid credentials
- Missing Authorization header

**Response**:
```json
{
  "error": "unauthorized",
  "message": "Authentication required"
}
```

---

### 3. Authorization Errors (403 Forbidden)

**When**: User authenticated but lacks permission

**Examples**:
- User tries to create committee without permission
- User tries to view other jurisdiction's data

**Response**:
```json
{
  "error": "forbidden",
  "message": "You do not have permission to perform this action"
}
```

---

### 4. Not Found Errors (404)

**When**: Requested resource doesn't exist

**Response**:
```json
{
  "error": "not_found",
  "message": "Committee not found"
}
```

---

### 5. Conflict Errors (409)

**When**: Operation conflicts with current state

**Examples**:
- Duplicate active committee in jurisdiction
- Duplicate phone number

**Response**:
```json
{
  "error": "conflict",
  "message": "An active committee already exists in this jurisdiction"
}
```

---

### 6. Server Errors (500 Internal Server Error)

**When**: Unexpected system error

**Examples**:
- Database connection lost
- Unhandled exception

**Response**:
```json
{
  "error": "internal_server_error",
  "message": "An unexpected error occurred. Please try again later.",
  "request_id": "req-uuid"
}
```

**Log** (but don't expose to user):
```
ERROR: Database connection failed: dial tcp: connection refused
Request ID: req-uuid
Stack trace: ...
```

---

## Language-Specific Patterns

### Go

**Error Wrapping**:
```go
if err != nil {
    return nil, fmt.Errorf("failed to create committee: %w", err)
}
```

**Custom Errors**:
```go
var (
    ErrDuplicateCommittee = errors.New("committee already exists in jurisdiction")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidPhone = errors.New("invalid phone format")
)

// Usage
if exists {
    return nil, ErrDuplicateCommittee
}

// Error checking
if errors.Is(err, ErrDuplicateCommittee) {
    // Handle duplicate
}
```

**Structured Errors**:
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Usage
return &ValidationError{Field: "phone", Message: "invalid format"}
```

---

### Python

**Try/Except**:
```python
try:
    committee = create_committee(jurisdiction_id)
except DuplicateCommitteeError as e:
    raise HTTPException(status_code=409, detail=str(e))
except DatabaseError as e:
    logger.error(f"Database error: {e}", exc_info=True)
    raise HTTPException(status_code=500, detail="Internal server error")
```

**Custom Exceptions**:
```python
class BJDMSError(Exception):
    """Base exception for BJDMS"""
    pass

class ValidationError(BJDMSError):
    def __init__(self, field: str, message: str):
        self.field = field
        self.message = message
        super().__init__(f"{field}: {message}")

class UnauthorizedError(BJDMSError):
    pass
```

---

### TypeScript/JavaScript

**Try/Catch**:
```typescript
try {
  const committee = await createCommittee(jurisdictionId);
  return committee;
} catch (error) {
  if (error instanceof DuplicateCommitteeError) {
    return res.status(409).json({
      error: 'conflict',
      message: error.message
    });
  }
  
  // Log and return generic error
  logger.error('Unexpected error:', error);
  return res.status(500).json({
    error: 'internal_server_error',
    message: 'An unexpected error occurred'
  });
}
```

**Custom Error Classes**:
```typescript
class BJDMSError extends Error {
  constructor(message: string, public code: string) {
    super(message);
    this.name = 'BJDMSError';
  }
}

class ValidationError extends BJDMSError {
  constructor(public field: string, message: string) {
    super(message, 'validation_error');
  }
}

// Usage
throw new ValidationError('phone', 'Invalid Bangladesh phone format');
```

---

## HTTP Error Response Format

**Standard Structure**:
```json
{
  "error": "error_code",
  "message": "Human-readable message",
  "details": [/* Optional array of field-specific errors */],
  "request_id": "req-uuid"  // For support/debugging
}
```

**Implementation** (Go):
```go
type ErrorResponse struct {
    Error     string          `json:"error"`
    Message   string          `json:"message"`
    Details   []FieldError    `json:"details,omitempty"`
    RequestID string          `json:"request_id"`
}

type FieldError struct {
    Field   string `json:"field"`
    Error   string `json:"error"`
}

func WriteError(w http.ResponseWriter, statusCode int, err error, requestID string) {
    resp := ErrorResponse{
        Error:     errorCode(err),
        Message:   err.Error(),
        RequestID: requestID,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteStatus(statusCode)
    json.NewEncoder(w).Encode(resp)
}
```

---

## Database Error Handling

### Connection Errors

```go
db, err := sql.Open("postgres", dsn)
if err != nil {
    log.Fatal("Failed to connect to database:", err)
}

// Retry logic
for i := 0; i < 3; i++ {
    err = db.Ping()
    if err == nil {
        break
    }
    time.Sleep(time.Second * 2)
}

if err != nil {
    return fmt.Errorf("database unavailable after retries: %w", err)
}
```

---

### Query Errors

```go
row := db.QueryRow("SELECT * FROM users WHERE id = $1", userID)
var user User
err := row.Scan(&user.ID, &user.Name, ...)

if err == sql.ErrNoRows {
    return nil, ErrUserNotFound  // Not Found
} else if err != nil {
    log.Error("Database query failed:", err)
    return nil, fmt.Errorf("failed to fetch user: %w", err)
}
```

---

### Transaction Errors

```go
tx, err := db.Begin()
if err != nil {
    return fmt.Errorf("failed to start transaction: %w", err)
}
defer tx.Rollback()  // Auto-rollback if commit not called

err = createCommittee(tx, ...)
if err != nil {
    return err  // Transaction auto-rolled back
}

err = tx.Commit()
if err != nil {
    return fmt.Errorf("failed to commit transaction: %w", err)
}
```

---

## Retry Logic

**When to Retry**:
- Network failures (temporary)
- Database connection timeouts
- External API rate limits

**When NOT to Retry**:
- Validation errors (will fail again)
- Authentication errors
- Business logic errors (duplicate committee)

**Implementation**:
```go
func retryWithBackoff(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        // Check if retryable
        if !isRetryable(err) {
            return err
        }
        
        // Exponential backoff
        wait := time.Duration(math.Pow(2, float64(i))) * time.Second
        time.Sleep(wait)
    }
    
    return fmt.Errorf("max retries exceeded")
}

func isRetryable(err error) bool {
    // Network errors, timeouts are retryable
    // Validation, auth errors are not
    return strings.Contains(err.Error(), "connection refused") ||
           strings.Contains(err.Error(), "timeout")
}
```

---

## Panic Recovery

**Go**: Use middleware to recover from panics

```go
func RecoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Error("Panic recovered:", err)
                log.Error("Stack trace:", string(debug.Stack()))
                
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}
```

---

## User-Facing Error Messages

**Rules**:
- Never expose stack traces to users
- Never expose database error details
- Never expose file paths, internal IDs
- Use generic messages for 500 errors

**Good**:
- "Invalid phone number format. Please use +880XXXXXXXXXX"
- "An active committee already exists in this jurisdiction"
- "You do not have permission to perform this action"

**Bad**:
- "SQL error: duplicate key value violates unique constraint 'committees_jurisdiction_id_idx'"
- "Panic: nil pointer dereference at /opt/bjdms/internal/service/committee. go:123"

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
