# SECURITY_POLICY.md

## Purpose

This document defines the **security policies, threat model, and security controls** for Bangladesh Jatiotabadi Jubodal Management System (BJDMS).

**Scope**: Application-level security, infrastructure security guidelines, and security operations.

---

## Security Principles

1. **Defense in Depth** - Multiple layers of security controls
2. **Least Privilege** - Minimum necessary permissions
3. **Zero Trust** - Verify every request, never assume
4. **Security by Default** - Secure settings out of the box
5. **Audit Everything** - All security events logged
6. **Privacy First** - Protect user data, especially anonymous complaints

---

## Threat Model

### Assets to Protect

| Asset | Sensitivity | Impact if Compromised |
|-------|-------------|----------------------|
| User credentials (passwords) | **Critical** | Account takeover, unauthorized access |
| Anonymous complaint submitter identity | **Critical** | Whistleblower exposure, safety risk |
| Committee member personal data (phone, NID) | **High** | Privacy violation, identity theft |
| Financial transaction records | **High** | Financial fraud, corruption |
| Audit logs | **High** | Evidence tampering, accountability loss |
| Donation records | **Medium** | Donor privacy violation |
| Activity/event data | **Low** | Limited impact (mostly public info) |

---

### Threat Actors

1. **External Attackers** (Hackers, Script Kiddies)
   - **Motivation**: Data theft, system disruption, ransom
   - **Capabilities**: SQL injection, XSS, brute force, DDoS
   - **Mitigation**: WAF, input validation, rate limiting, HTTPS

2. **Malicious Insiders** (Compromised Committee Members)
   - **Motivation**: Political sabotage, data leakage
   - **Capabilities**: Legitimate access, privilege abuse
   - **Mitigation**: RBAC, audit logging, separation of duties

3. **Political Adversaries**
   - **Motivation**: Disrupt operations, discredit organization
   - **Capabilities**: DDoS, social engineering, information warfare
   - **Mitigation**: DDoS protection (Cloudflare), backup systems, incident response plan

4. **Curious Public**
   - **Motivation**: Access private information, anonymous complaint data
   - **Capabilities**: Public API probing, brute force
   - **Mitigation**: Rate limiting, captcha, public/private API separation

---

### Attack Scenarios & Mitigations

#### Scenario 1: Anonymous Complaint De-anonymization

**Attack**: Attacker tries to identify anonymous complaint submitter via timing analysis, IP correlation, or database breach.

**Mitigations**:
- IP addresses hashed (SHA-256) immediately on receipt
- Metadata stored in separate restricted table (`anonymous_complaint_metadata`)
- Metadata deleted after 90 days
- Access logging for all evidence views
- No direct link between complaint and submitter in main database

**Detection**: Monitor unusual access patterns to `anonymous_complaint_metadata`, multiple failed authorization attempts

---

#### Scenario 2: Credential Stuffing / Brute Force

**Attack**: Attacker uses leaked credentials from other breaches or brute forces passwords.

**Mitigations**:
- Rate limiting: 10 login attempts per 15 minutes per IP
- Account lockout: 5 failed attempts → 30-minute lockout
- Password complexity requirements (8+ chars, uppercase, lowercase, digit)
- Bcrypt with cost factor 12 (slow hashing)
- Optional: 2FA/MFA for sensitive roles (Super Admin, Treasurer)

**Detection**: Monitor failed login attempts, lockout events, login from unusual locations

---

#### Scenario 3: SQL Injection

**Attack**: Attacker injects SQL code via form inputs to extract database data.

**Mitigations**:
- **Parameterized queries only** (prepared statements)
- ORM usage (prevents most SQL injection)
- Input validation on all API endpoints
- Least privilege database user (app user cannot DROP tables)
- WAF (Web Application Firewall) rules

**Detection**: WAF alerts on SQL injection patterns, unexpected database errors in logs

---

#### Scenario 4: Cross-Site Scripting (XSS)

**Attack**: Attacker injects JavaScript into user input fields (e.g., complaint description, activity title) that executes in other users' browsers.

**Mitigations**:
- **Output encoding** for all user-generated content
- Content Security Policy (CSP) headers
- DOMPurify or similar sanitization library
- Framework-level XSS protection (Next.js escapes by default)

**Detection**: CSP violation reports, unusual JavaScript execution patterns

---

#### Scenario 5: Privilege Escalation

**Attack**: Regular member tries to access Super Admin functions by manipulating API requests.

**Mitigations**:
- **Server-side authorization checks** on every request (never trust client)
- JWT claims verified on backend
- Role-based and jurisdiction-based filtering
- Audit log all permission denials

**Detection**: Failed authorization attempts logged, alerts on repeated attempts from same user

---

#### Scenario 6: Data Breach (Database Dump)

**Attack**: Attacker gains database access and dumps all data.

**Mitigations**:
- **Encryption at rest** for sensitive columns (passwords via bcrypt, optional: PGP for NID)
- Database firewall rules (no external access)
- Principle of least privilege for database users
- Regular security audits and penetration testing
- Database access logging

**Detection**: Unusual database queries, large data exports, access from unauthorized IPs

---

## Security Controls

### 1. Authentication & Authorization

✅ **Implemented** (See PHASE_C_AUTH.md):
- JWT-based authentication
- Bcrypt password hashing (cost 12)
- Refresh token rotation
- RBAC + ABAC permission model
- Account lockout after 5 failed attempts

🔄 **Recommended Enhancements**:
- Multi-Factor Authentication (MFA) for Super Admin, Treasurer, Central Leaders
- Passwordless authentication via SMS OTP (supplement to password)
- Session timeout after 15 minutes of inactivity
- Concurrent session limit (max 3 devices per user)

---

### 2. Network Security

✅ **HTTPS Only**:
- All traffic encrypted via TLS 1.3
- NGINX reverse proxy terminates SSL
- HTTP → HTTPS redirect enforced
- HSTS (HTTP Strict Transport Security) header enabled

✅ **Firewall Rules**:
- PostgreSQL port 5432: Internal network only
- Redis port 6379: Internal network only
- OpenSearch port 9200: Internal network only
- NGINX port 443: Public internet
- SSH port 25920: Restricted IPs only

🔄 **Recommended**:
- DDoS protection via Cloudflare (or similar CDN)
- Rate limiting at NGINX level (supplement to application-level)
- IP whitelisting for admin panel access

---

### 3. Input Validation

**All API Endpoints**:
- **Phone number**: Regex validation `^\+880\d{10}$`
- **Email**: RFC 5322 compliant validation
- **NID**: 10 or 13 digits only
- **File uploads**: MIME type validation, virus scanning (ClamAV)
- **Text fields**: Max length enforcement, special character escaping

**File Upload Security**:
- Max file size: 10MB (images), 100MB (videos)
- Allowed extensions: jpg, png, pdf, mp4, mov
- EXIF data stripped from images (anonymity)
- Files stored with random UUIDs (no original filenames in path)
- Virus scan before storage

---

### 4. API Security

**Rate Limiting**:
| Endpoint | Limit | Window |
|----------|-------|--------|
| `/api/v1/auth/login` | 10 requests | 15 minutes |
| `/api/v1/auth/refresh` | 20 requests | 15 minutes |
| `/api/v1/public/join-requests` | 3 requests | 24 hours |
| `/api/v1/public/complaints/submit` | 10 requests | 24 hours |
| All other endpoints | 100 requests | 1 minute (per user) |

**CORS (Cross-Origin Resource Sharing)**:
- Allowed origins: `https://bjdms.arint.win`, `https://www.bjdms.arint.win`
- Credentials allowed: Yes (for cookies/JWT)
- Allowed methods: GET, POST, PATCH, DELETE
- No `Access-Control-Allow-Origin: *` (too permissive)

**Security Headers** (NGINX configuration):
```nginx
add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "DENY" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Permissions-Policy "geolocation=(), microphone=(), camera=()" always;
add_header Content-Security-Policy "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:;" always;
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;
```

---

### 5. Data Protection

**Encryption**:
- **In Transit**: TLS 1.3 for all HTTP traffic, SSH for server access
- **At Rest**: 
  - Passwords: Bcrypt (cost 12)
  - Optional: PGP encryption for sensitive fields (NID, complaint evidence)
  - Database volume encryption (LUKS or cloud provider encryption)

**Data Minimization**:
- Collect only necessary data
- Anonymous complaints: No PII beyond what's required
- Phone numbers masked in public directory
- Email optional for join requests

**Data Retention**:
- Audit logs: 5 years
- Financial records: 7 years (legal requirement)
- Anonymous complaint metadata: 90 days, then deleted
- Soft-deleted users: Retained for 2 years, then purged
- Activity proofs: 2 years, then archived to cold storage

---

### 6. Secure Development Practices

**Code Review**:
- All code changes require review before merge
- Security-focused review for authentication, authorization, financial modules

**Dependency Management**:
- Automated dependency scanning (npm audit, go mod verify)
- Update critical vulnerabilities within 7 days
- Monthly dependency updates for non-critical

**Static Analysis**:
- Go: `gosec` for security scanning
- Python: `bandit` for security issues
- JavaScript/TypeScript: ESLint with security rules

**Secret Management**:
- **No secrets in code** (environment variables only)
- `.env` file on server only (never committed to git)
- Rotate secrets quarterly (database passwords, JWT signing keys)

---

## Incident Response Plan

### Severity Levels

| Level | Description | Response Time |
|-------|-------------|---------------|
| **P0 - Critical** | Data breach, system down, active attack | Immediate (< 1 hour) |
| **P1 - High** | Privilege escalation, DoS, significant vulnerability | < 4 hours |
| **P2 - Medium** | Minor vulnerability, failed attack attempt | < 24 hours |
| **P3 - Low** | Security audit finding, non-critical issue | < 1 week |

---

### Incident Response Steps

1. **Detection**: Automated alerts (failed logins, unusual DB queries, WAF triggers)
2. **Containment**: Isolate affected systems, block attacker IPs, revoke compromised credentials
3. **Investigation**: Analyze logs, determine extent of breach, identify attack vector
4. **Eradication**: Patch vulnerability, remove malicious code, reset compromised accounts
5. **Recovery**: Restore services, verify integrity, monitor for re-attack
6. **Post-Incident**: Document incident, update security controls, notify affected users (if required)

---

### Breach Notification

**If Personal Data Breach Occurs**:
- Notify affected users within 72 hours (via SMS/email)
- Report to relevant authorities (if legally required in Bangladesh)
- Public statement if breach affects >1000 users
- Offer remediation (password reset, account monitoring)

---

## Security Monitoring

### Automated Alerts

1. **Failed Login Spikes**: >50 failed logins in 5 minutes → Alert ops team
2. **Privilege Escalation Attempts**: Any failed authorization for Super Admin actions → Alert
3. **Unusual Database Access**: Queries to `anonymous_complaint_metadata` by non-Super Admin → Alert
4. **Large Data Exports**: >1000 records exported in single query → Review
5. **File Upload Anomalies**: Virus detected, suspicious file types → Block + Alert

---

### Penetration Testing

**Frequency**: Annually (minimum)

**Scope**:
- Web application security (OWASP Top 10)
- API security
- Authentication bypass attempts
- Privilege escalation
- Database security

**Post-Test**: Address all High/Critical findings within 30 days

---

## Vulnerability Disclosure

**Responsible Disclosure Policy**:
- Email: security@bjd.org.bd
- PGP key provided for encrypted reports
- Response within 48 hours
- Fix timeline: 30 days for High/Critical, 90 days for Medium/Low
- Bounty program (optional): Tk 5,000 - 50,000 based on severity

---

## Compliance

- Bangladesh ICT Act compliance
- GDPR-equivalent data protection (for international users)
- Financial audit requirements (7-year retention)

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
