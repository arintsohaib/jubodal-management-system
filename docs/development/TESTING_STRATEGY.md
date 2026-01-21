# TESTING_STRATEGY.md

## Purpose

Defines **testing approach and requirements** for BJDMS to ensure code quality, prevent regressions, and enable confident deployments.

**Goal**: 70% overall coverage, 90% for critical modules, all tests automated and running in CI.

---

## Testing Pyramid

```
      /\
     /  \    E2E Tests (10%)
    /____\   - Full user flows
   /      \  Integration Tests (30%)
  /        \ - API + Database + Services
 /__________\ Unit Tests (60%)
              - Pure logic, isolated functions
```

---

## Test Types

### 1. Unit Tests

**What**: Test individual functions/methods in isolation

**Tools**:
- Go: `testing` package
- Python: `pytest`
- JavaScript: `Jest`

**Coverage Target**: 80%

**Example** (Go):
```go
func TestValidatePhone_ValidBangladeshNumber_ReturnsNil(t *testing.T) {
    err := ValidatePhone("+8801712345678")
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}

func TestValidatePhone_InvalidFormat_ReturnsError(t *testing.T) {
    err := ValidatePhone("1234567890")
    if err == nil {
        t.Error("Expected error for invalid format, got nil")
    }
}
```

---

### 2. Integration Tests

**What**: Test multiple components together (API + Database, Service + Repository)

**Tools**:
- Go: `testing` + `testcontainers` (Docker-based test DB)
- Python: `pytest` + `pytest-docker`

**Coverage Target**: 60%

**Example** (Go):
```go
func TestCommitteeService_Create_SavesToDatabase(t *testing.T) {
    // Setup: Start test database container
    ctx := context.Background()
    container := startTestPostgres(t)
    defer container.Terminate(ctx)
    
    db := connectToDB(container)
    repo := NewCommitteeRepository(db)
    service := NewCommitteeService(repo)
    
    // Execute
    committee, err := service.CreateCommittee(ctx, CreateCommitteeRequest{
        JurisdictionID: uuid.New(),
        CommitteeType:  "Full",
    })
    
    // Verify
    assert.NoError(t, err)
    assert.NotNil(t, committee)
    
    // Verify database state
    saved, _ := repo.GetByID(ctx, committee.ID)
    assert.Equal(t, "Full", saved.CommitteeType)
}
```

---

### 3. End-to-End Tests

**What**: Test complete user flows (login → create committee → add member)

**Tools**:
- **Playwright** or **Cypress** (browser automation)
- **API tests**: Postman/Newman, REST Client

**Coverage Target**: Critical user paths only (~10 flows)

**Example** (Playwright):
```javascript
test('committee leader can create and dissolve committee', async ({ page }) => {
  // Login
  await page.goto('https://grayhawks.com/login');
  await page.fill('[name=phone]', '+8801712345678');
  await page.fill('[name=password]', 'TestPass123');
  await page.click('button[type=submit]');
  
  // Create committee
  await page.click('text=Create Committee');
  await page.fill('[name=jurisdiction]', 'Ward 5, Joypurhat');
  await page.selectOption('[name=type]', 'Full');
  await page.click('button:has-text("Create")');
  
  // Verify success
  await expect(page.locator('text=Committee created successfully')).toBeVisible();
  
  // Dissolve committee
  await page.click('text=Dissolve Committee');
  await page.click('button:has-text("Confirm")');
  
  // Verify dissolved
  await expect(page.locator('text=Committee dissolved')).toBeVisible();
});
```

---

## Critical Test Scenarios

### Authentication & Authorization

- [ ] Login with valid credentials → Success
- [ ] Login with invalid password → Fail
- [ ] Account lockout after 5 failed attempts
- [ ] Refresh token → New access token
- [ ] Expired token → 401 Unauthorized
- [ ] Unauthorized access → 403 Forbidden
- [ ] Permission check: User without permission → Denied

---

### Committee Management

- [ ] Create Full Committee → Success
- [ ] Create duplicate committee in same jurisdiction → Error
- [ ] Convener → Full transition → Members carried over
- [ ] Position uniqueness: Assign 2nd President → Error
- [ ] Dissolve committee → All members ended
- [ ] Committee size limit: Add 152nd member to district → Error

---

### Join Requests

- [ ] Submit join request (public, no login) → Pending
- [ ] Duplicate phone number → Error
- [ ] Underage applicant (<18) → Error
- [ ] Approve join request → User created
- [ ] Reject join request → Notification sent
- [ ] Assign position → Dashboard access granted (verified_at set)

---

### Complaint System

- [ ] Submit anonymous complaint → submitter_id = NULL
- [ ] Submit named complaint → logged
- [ ] Evidence upload → File stored, virus scanned
- [ ] Route to district leaders → Assigned
- [ ] Escalate complaint → Higher jurisdiction
- [ ] Anonymous complaint metadata → Only Super Admin can access

---

### Financial System

- [ ] Record donation → Transaction created
- [ ] Receipt generated → PDF downloadable
- [ ] Expense approval (>Tk 10,000) → Requires 2 approvers
- [ ] Fund transfer → Both balances updated atomically
- [ ] Transaction immutability → UPDATE/DELETE rejected

---

## Test Data Management

### Test Database

**Setup**:
- Use Docker container (PostgreSQL) for test DB
- Isolated from production and development DBs
- Reset before each test suite

**Migration**:
```bash
# Before tests, apply migrations to test DB
migrate -path migrations -database "postgresql://test_user:test_pass@localhost:5433/test_db" up
```

**Cleanup**:
```go
func teardown(t *testing.T, db *sql.DB) {
    db.Exec("TRUNCATE users, committees, committee_members CASCADE")
}
```

---

### Test Fixtures

**Use Factory Pattern**:
```go
func createTestUser(t *testing.T, db *sql.DB) *User {
    user := &User{
        FullName: "Test User",
        Phone:    "+880171" + randomDigits(7),
        PasswordHash: hashPassword("TestPass123"),
    }
    db.Exec("INSERT INTO users (...) VALUES (...)")
    return user
}

func createTestCommittee(t *testing.T, db *sql.DB, jurisdictionID uuid.UUID) *Committee {
    // ...
}
```

---

## Mocking

**When to Mock**:
- External APIs (SMS gateway, payment gateway)
- Slow operations (file uploads to S3)
- Non-deterministic functions (time.Now(), random generators)

**Go** (use interfaces + mock implementations):
```go
type SMSGateway interface {
    SendSMS(phone, message string) error
}

// Production
type TwilioSMSGateway struct { ... }

// Test
type MockSMSGateway struct {
    SentMessages []string
}

func (m *MockSMSGateway) SendSMS(phone, message string) error {
    m.SentMessages = append(m.SentMessages, message)
    return nil
}

// In test
func TestJoinRequestApproval_SendsSMS(t *testing.T) {
    sms := &MockSMSGateway{}
    service := NewJoinRequestService(repo, sms)
    
    service.ApproveRequest(ctx, requestID)
    
    assert.Equal(t, 1, len(sms.SentMessages))
    assert.Contains(t, sms.SentMessages[0], "approved")
}
```

---

## Test Execution

### Local (Remote Server Only)

**Remember**: No local execution allowed per user rules.

```bash
# On grayhawks.com
ssh root@grayhawks.com -p 25920
cd /opt/bjdms

# Run unit tests
go test ./... -v

# Run with coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

---

### CI/CD Pipeline

**On every commit** (GitHub Actions, GitLab CI):
```yaml
test:
  runs-on: ubuntu-latest
  services:
    postgres:
      image: postgres:15
      env:
        POSTGRES_PASSWORD: test
      options: >-
        --health-cmd pg_isready
        --health-interval 10s
        --health-timeout 5s
        --health-retries 5
  steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - name: Run migrations
      run: migrate -path migrations -database $DATABASE_URL up
    - name: Run tests
      run: go test ./... -v -cover
    - name: Check coverage
      run: |
        coverage=$(go test ./... -coverprofile=coverage.out | grep coverage | awk '{print $2}' | sed 's/%//')
        if (( $(echo "$coverage < 70" | bc -l) )); then
          echo "Coverage $coverage% is below 70%"
          exit 1
        fi
```

---

## Performance Testing

### Load Testing

**Tool**: `k6` or `Apache JBench`

**Scenario**: Simulate 1000 concurrent users

```javascript
import http from 'k6/http';
import { check } from 'k6';

export let options = {
  vus: 1000,  // Virtual users
  duration: '5m',
};

export default function() {
  let res = http.get('https://grayhawks.com/api/v1/activities');
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
}
```

**Targets**:
- p95 response time < 500ms
- Error rate < 1%
- System stable under 1000 concurrent users

---

## Security Testing

### Automated Security Scans

**Tools**:
- **SAST** (Static Application Security Testing): `gosec`, `bandit`, ESLint security plugins
- **Dependency Scanning**: `npm audit`, `go mod verify`, Snyk
- **Container Scanning**: Trivy, Clair

**On every PR**:
```bash
# Go security scan
gosec ./...

# Python security scan
bandit -r .

# npm audit
npm audit --audit-level=high
```

---

### Penetration Testing

**Manual Testing** (Quarterly):
- SQL injection attempts
- XSS injection
- CSRF attacks
- Authentication bypass
- Privilege escalation

**Automated** (Weekly):
- OWASP ZAP automated scan
- Nikto web server scan

---

## Test Documentation

**Each test should**:
- Have descriptive name (what it tests, expected outcome)
- Include comments for complex setup
- Assert specific, meaningful conditions

**Example**:
```go
// TestCreateCommittee_ExceedsMaxMembers_ReturnsError verifies that 
// attempting to add a member beyond the jurisdiction's committee size limit
// returns an error and does not modify the database.
func TestCreateCommittee_ExceedsMaxMembers_ReturnsError(t *testing.T) {
    // Setup: Create district committee with 151 members (max for district)
    committee := createTestCommittee(t, db, districtJurisdiction)
    for i := 0; i < 151; i++ {
        addTestMember(t, db, committee.ID)
    }
    
    // Execute: Attempt to add 152nd member
    err := service.AddMember(ctx, committee.ID, newMemberID)
    
    // Verify: Error returned, no new member in database
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "committee size limit exceeded")
    
    count := countMembers(t, db, committee.ID)
    assert.Equal(t, 151, count)
}
```

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
