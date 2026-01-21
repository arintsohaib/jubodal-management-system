# CODE_STANDARDS.md

## Purpose

Defines **coding conventions and best practices** for BJDMS development across all languages (Go, Python, JavaScript/TypeScript).

**Goal**: Maintainable, readable, secure code that survives leadership changes and scales over 5+ years.

---

## General Principles

1. **Clarity over Cleverness** - Code should be obvious, not clever
2. **Consistency** - Follow existing patterns
3. **Security First** - Always validate input, never trust user data
4. **Performance Second** - Optimize only after profiling
5. **Document Why, Not What** - Code shows what, comments explain why
6. **No Dead Code** - Remove unused code immediately

---

## Language-Specific Standards

### Go (Backend API)

**Style Guide**: [Effective Go](https://go.dev/doc/effective_go) + [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

**Naming**:
```go
// Good: Descriptive, clear
func CreateCommittee(ctx context.Context, req CreateCommitteeRequest) (*Committee, error)

// Bad: Abbreviated, unclear
func CrtCmt(c context.Context, r CmtReq) (*Cmt, error)
```

**Error Handling**:
```go
// Good: Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to create committee: %w", err)
}

// Bad: Lose context
if err != nil {
    return nil, err
}
```

**Project Structure**:
```
cmd/
  api/          # Main API server
  migrate/      # Database migrations
internal/
  auth/         # Authentication logic
  committee/    # Committee domain
  handler/      # HTTP handlers
  repository/   # Database access
  service/      # Business logic
pkg/
  validation/   # Shared validation
  errors/       # Custom errors
```

**Database Access**:
- Use `pgx` or `sqlx` for PostgreSQL
- Always use parameterized queries (prevent SQL injection)
- Repository pattern for data access
```go
type CommitteeRepository interface {
    Create(ctx context.Context, committee *Committee) error
    GetByID(ctx context.Context, id uuid.UUID) (*Committee, error)
    // ...
}
```

---

### Python (Scripts, ML/AI if used)

**Style Guide**: [PEP 8](https://peps.python.org/pep-0008/)

**Formatting**: Use `black` (auto-formatter)

**Type Hints**:
```python
# Good: Type hints for clarity
def calculate_budget_variance(
    allocated: Decimal, 
    spent: Decimal
) -> Decimal:
    return allocated - spent

# Bad: No type hints
def calc_var(a, s):
    return a - s
```

**Project Structure**:
```
bjdms/
  api/
    endpoints/
    models/
    services/
  core/
    config.py
    dependencies.py
  tests/
```

---

### JavaScript/TypeScript (Frontend)

**Style Guide**: [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)

**Use TypeScript**: Strict mode enabled
```typescript
// tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true
  }
}
```

**Component Structure** (React/Next.js):
```typescript
// Good: Typed props, clear interface
interface CommitteeCardProps {
  committee: Committee;
  onEdit: (id: string) => void;
}

export const CommitteeCard: React.FC<CommitteeCardProps> = ({ 
  committee, 
  onEdit 
}) => {
  // ...
};

// Bad: Untyped, unclear
export const Card = (props) => {
  // ...
};
```

---

## Security Standards

### Input Validation

**Rule**: Validate EVERYTHING from users/external sources

```go
// Good: Validate phone format
func ValidatePhone(phone string) error {
    matched, _ := regexp.MatchString(`^\+880\d{10}$`, phone)
    if !matched {
        return errors.New("invalid Bangladesh phone number format")
    }
    return nil
}
```

**Use Libraries**:
- Go: `go-playground/validator`
- Python: `pydantic`
- TypeScript: `zod`, `yup`

---

### SQL Injection Prevention

**Always use parameterized queries**:

```go
// Good: Parameterized
db.Query("SELECT * FROM users WHERE phone = $1", phone)

// BAD: String concatenation (SQL INJECTION!)
db.Query("SELECT * FROM users WHERE phone = '" + phone + "'")
```

---

### XSS Prevention

```javascript
// Good: React escapes by default
<div>{userInput}</div>

// BAD: Dangerous HTML injection
<div dangerouslySetInnerHTML={{__html: userInput}} />
```

---

### Authentication Checks

```go
// EVERY protected endpoint
func HandleCommitteeCreate(w http.ResponseWriter, r *http.Request) {
    // 1. Extract JWT
    claims, err := auth.ExtractClaims(r)
    if err != nil {
        http.Error(w, "Unauthorized", 401)
        return
    }
    
    // 2. Check permission
    if !auth.HasPermission(claims, "committee.create") {
        http.Error(w, "Forbidden", 403)
        return
    }
    
    // 3. Proceed with business logic
    // ...
}
```

---

## Code Review Checklist

### Before Submitting PR

- [ ] Code compiles/runs without errors
- [ ] All tests pass
- [ ] No linter warnings
- [ ] Security: Input validated, no SQL injection, no XSS
- [ ] Error handling: All errors checked and wrapped
- [ ] Logging: Critical actions logged
- [ ] Documentation: Complex logic commented
- [ ] No secrets in code (use .env)

### Reviewer Checklist

- [ ] Code follows standards (this document)
- [ ] Logic is correct
- [ ] Edge cases handled
- [ ] Performance acceptable (no obvious N+1 queries)
- [ ] Tests cover new code
- [ ] Security: No vulnerabilities
- [ ] Backward compatible (or migration plan documented)

---

## Testing Standards

(See TESTING_STRATEGY.md for full details)

**Minimum Coverage**: 70% overall, 90% for critical modules (auth, finance)

**Test Naming**:
```go
func TestCreateCommittee_Success(t *testing.T) { ... }
func TestCreateCommittee_DuplicateJurisdiction_ReturnsError(t *testing.T) { ... }
```

---

## Performance Guidelines

### Database Queries

**Avoid N+1**:
```go
// Bad: N+1 query
for _, committee := range committees {
    members := db.Query("SELECT * FROM members WHERE committee_id = $1", committee.ID)
    // ...
}

// Good: JOIN or batch query
committees := db.Query(`
    SELECT c.*, m.* 
    FROM committees c 
    LEFT JOIN members m ON c.id = m.committee_id
    WHERE c.jurisdiction_id = $1
`, jurisdictionID)
```

**Use Indexes**:
- Index foreign keys
- Index frequently queried columns (phone, jurisdiction_id)
- Composite indexes for multi-column queries

**Pagination**:
```go
// Always paginate large result sets
func ListActivities(limit, offset int) ([]Activity, error) {
    return db.Query("SELECT * FROM activities ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
}
```

---

## Git Workflow

### Commit Messages

**Format**: `<type>: <description>`

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `refactor`: Code refactoring
- `test`: Add/update tests
- `chore`: Build/tooling changes

**Examples**:
```
feat: add Municipality level to jurisdiction hierarchy
fix: SQL injection vulnerability in complaint search
docs: update DATABASE_SCHEMA.md with join_requests table
refactor: extract permission check into middleware
```

**Body** (optional): Explain why, not what
```
feat: add anonymous complaint submission

Allows public users to submit complaints without login.
Implements IP hashing for anonymity preservation.
Addresses requirement from PHASE_F_COMPLAINT.md.
```

---

### Branch Naming

- `main`: Production-ready code
- `feature/add-municipality`: New features
- `fix/sql-injection-search`: Bug fixes
- `hotfix/security-patch`: Urgent production fixes

---

## Code Formatting

### Auto-Formatting (Required)

- **Go**: `gofmt` or `goimports`
- **Python**: `black`
- **JavaScript/TypeScript**: `prettier`

**Pre-commit Hook**:
```bash
#!/bin/sh
# .git/hooks/pre-commit

# Format Go code
go fmt ./...

# Format Python
black .

# Format TypeScript
npx prettier --write "src/**/*.{ts,tsx}"

git add -u
```

---

## Documentation

### Code Comments

**When to Comment**:
- Complex algorithms
- Non-obvious business logic
- Security decisions
- Performance optimizations
- Workarounds for external issues

**When NOT to Comment**:
```go
// Bad: Obvious comment
// Set user ID to the ID from request
userID := req.UserID

// Good: Explains why
// Use older ID format for backward compatibility with mobile app v1.x
userID := convertToLegacyID(req.UserID)
```

### Function Documentation

**Go**:
```go
// CreateCommittee creates a new committee in the specified jurisdiction.
// It validates that no active committee exists in the same jurisdiction
// and enforces hierarchical permissions (only parent can create child committee).
//
// Returns ErrDuplicateCommittee if active committee exists.
// Returns ErrUnauthorized if user lacks permission.
func CreateCommittee(ctx context.Context, req CreateCommitteeRequest) (*Committee, error) {
    // ...
}
```

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
