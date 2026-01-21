# AUDIT_COMPLIANCE.md

## Purpose

This document defines **audit requirements and compliance standards** for Bangladesh Jatiotabadi Jubodal Management System (BJDMS).

**Goal**: Ensure complete transparency, accountability, and audit-readiness for all organizational activities.

---

## Audit Principles

1. **Complete Audit Trail** - Every action logged with who, what, when, where
2. **Immutability** - Audit logs cannot be edited or deleted
3. **Real-Time Logging** - Actions logged immediately, not batched
4. **Accessibility** - Audit logs queryable by authorized users
5. **Retention** - Audit logs retained for 5 years minimum
6. **Transparency** - Regular audit reports published

---

## What Gets Audited

### Authentication Events

| Event | Logged Data | Retention |
|-------|-------------|-----------|
| Login success | user_id, IP, timestamp, user_agent | 5 years |
| Login failure | phone, IP, timestamp, reason | 5 years |
| Account lockout | user_id, IP, timestamp | 5 years |
| Password reset | user_id, IP, timestamp | 5 years |
| Token refresh | user_id, timestamp | 5 years |
| Logout | user_id, timestamp | 5 years |

---

### Committee Management Events

| Event | Logged Data | Retention |
|-------|-------------|-----------|
| Committee created | creator_id, committee_id, jurisdiction, type, timestamp | Indefinite |
| Committee dissolved | dissolver_id, committee_id, reason, timestamp | Indefinite |
| Member added | adder_id, committee_id, member_id, position, timestamp | Indefinite |
| Member removed | remover_id, committee_id, member_id, reason, timestamp | Indefinite |
| Position changed | changer_id, member_id, old_position, new_position, timestamp | Indefinite |
| Convener → Full transition | initiator_id, old_committee_id, new_committee_id, timestamp | Indefinite |

**Rationale for Indefinite**: Organizational history, leadership accountability.

---

### Financial Events

| Event | Logged Data | Retention |
|-------|-------------|-----------|
| Donation recorded | recorder_id, donor_id, amount, purpose, timestamp | 7 years (legal) |
| Expense created | requester_id, amount, category, timestamp | 7 years |
| Expense approved | approver_id, expense_id, timestamp | 7 years |
| Expense rejected | approver_id, expense_id, reason, timestamp | 7 years |
| Fund transfer | from_account, to_account, amount, approver_id, timestamp | 7 years |
| Budget created | creator_id, jurisdiction, fiscal_year, timestamp | 7 years |

---

### Activity & Event Management

| Event | Logged Data | Retention |
|-------|-------------|-----------|
| Activity submitted | user_id, activity_id, jurisdiction, timestamp | 2 years |
| Activity approved | approver_id, activity_id, timestamp | 2 years |
| Activity rejected | approver_id, activity_id, reason, timestamp | 2 years |
| Event created | creator_id, event_id, jurisdiction, timestamp | 2 years |
| Event attendance | user_id, event_id, checked_in_at | 2 years |

---

### Complaint System Events

| Event | Logged Data | Retention |
|-------|-------------|-----------|
| Complaint submitted (named) | submitter_id, target_id, timestamp | Indefinite |
| Complaint submitted (anonymous) | complaint_id, IP_hash, timestamp | Indefinite (metadata 90 days) |
| Complaint status changed | changer_id, complaint_id, old_status, new_status, timestamp | Indefinite |
| Complaint evidence viewed | viewer_id, evidence_id, timestamp | Indefinite |
| Complaint escalated | escalator_id, complaint_id, from_jurisdiction, to_jurisdiction, timestamp | Indefinite |

**Special Handling**: Anonymous complaint submitter identity NOT logged (privacy protection).

---

### Permission & Access Control

| Event | Logged Data | Retention |
|-------|-------------|-----------|
| Permission granted | granter_id, user_id, permission, timestamp | 5 years |
| Permission revoked | revoker_id, user_id, permission, timestamp | 5 years |
| Unauthorized access attempt | user_id, resource, IP, timestamp | 5 years |
| Role assigned | assigner_id, user_id, role, jurisdiction, timestamp | 5 years |

---

### Data Access

| Event | Logged Data | Retention |
|-------|-------------|-----------|
| Sensitive data export | user_id, data_type, record_count, timestamp | 5 years |
| Anonymous complaint metadata accessed | user_id, complaint_id, timestamp | 5 years |
| Financial report generated | user_id, report_type, jurisdiction, timestamp | 5 years |

---

## Audit Log Schema

**Table**: `audit_logs` (defined in DATABASE_SCHEMA.md)

```sql
CREATE TABLE audit_logs (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),  -- nullable for system actions
  action VARCHAR(100) NOT NULL,  -- 'login_success', 'committee_created', etc.
  entity VARCHAR(50) NOT NULL,  -- 'users', 'committees', 'complaints', etc.
  entity_id UUID,  -- ID of affected entity
  old_value JSONB,  -- optional: previous state
  new_value JSONB,  -- optional: new state
  ip_address VARCHAR(45),  -- IPv4 or IPv6
  user_agent TEXT,  -- Browser/device info
  metadata JSONB,  -- Additional contextual data
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for query performance
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_action ON audit_logs(action);
CREATE INDEX idx_audit_entity ON audit_logs(entity, entity_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at DESC);

-- Immutability: Prevent UPDATE/DELETE (application-level enforcement + DB trigger)
CREATE OR REPLACE FUNCTION prevent_audit_modification()
RETURNS TRIGGER AS $$
BEGIN
  RAISE EXCEPTION 'Audit logs are immutable';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_immutable_trigger
BEFORE UPDATE OR DELETE ON audit_logs
FOR EACH ROW EXECUTE FUNCTION prevent_audit_modification();
```

---

## Audit Query Examples

### Who logged in from unusual location?

```sql
SELECT user_id, ip_address, created_at
FROM audit_logs
WHERE action = 'login_success'
  AND ip_address NOT IN (SELECT DISTINCT ip_address FROM audit_logs WHERE user_id = target_user_id LIMIT 10)
ORDER BY created_at DESC;
```

---

### All actions by a specific user

```sql
SELECT action, entity, entity_id, created_at
FROM audit_logs
WHERE user_id = 'target-user-uuid'
ORDER BY created_at DESC;
```

---

### Who approved this expense?

```sql
SELECT user_id, created_at
FROM audit_logs
WHERE action = 'expense_approved'
  AND entity_id = 'expense-uuid';
```

---

### All committee creations in last 30 days

```sql
SELECT user_id, entity_id, metadata, created_at
FROM audit_logs
WHERE action = 'committee_created'
  AND created_at > NOW() - INTERVAL '30 days'
ORDER BY created_at DESC;
```

---

## Audit Reporting

### Monthly Audit Report

**Generated By**: Super Admin, Central Leaders

**Contents**:
- Total users added
- Committees created/dissolved
- Financial transactions summary (donations, expenses)
- Activities submitted/approved
- Complaints received/resolved
- Security events (failed logins, lockouts)

**Distribution**: Central Committee, optionally public (anonymized)

---

### Annual Transparency Report

**Public-Facing Document**:
- Total membership growth
- Number of active committees by level
- Financial summary (total donations, expenses by category)
- Activities by district
- Complaint system usage (total complaints, resolution rate)

**Purpose**: Organizational accountability, public trust

---

## Compliance Audits

### Internal Audit (Quarterly)

**Conducted By**: BJDMS Audit Committee or designated officer

**Scope**:
- Review sample transactions for proper approval
- Verify audit log completeness (no gaps)
- Check unauthorized access attempts
- Validate data retention compliance

**Deliverable**: Audit findings report, corrective action recommendations

---

### External Audit (Annual)

**Conducted By**: Independent auditor (if required by law or donors)

**Scope**:
- Financial audit (7-year retention verified)
- Data protection compliance
- Security controls assessment
- Audit log integrity check

**Deliverable**: Audit opinion, compliance certificate

---

## Regulatory Compliance

### Bangladesh Legal Requirements

**ICT Act Compliance**:
- Activity logs demonstrate lawful use
- No evidence of prohibited content
- Cooperation with law enforcement (audit trail of legal requests)

**Financial Transparency** (if tax-exempt status):
- 7-year financial record retention
- Donation receipts issued and auditable
- Annual financial statement filed

---

### Donor Compliance (if applicable)

**If receiving international funding**:
- Demonstrate transparent fund usage (audit trail)
- Financial reports aligned with grant requirements
- Compliance with donor country laws (e.g., US FCPA, UK Bribery Act)

---

## Audit Log Access Control

**Who Can View Audit Logs**:

| Role | Access Level |
|------|--------------|
| Super Admin | All audit logs |
| Central Leaders | Central + all child jurisdictions |
| District Leaders | District + child jurisdictions (Upazila, Union, Ward) |
| Committee Treasurers | Financial audit logs for their jurisdiction |
| Regular Members | Own actions only |

**API Endpoint**: `GET /api/v1/audit-logs?user={id}&action={type}&from={date}&to={date}`

**Permission Required**: `audit.view`

---

## Audit Log Integrity

### Tamper Prevention

1. **Database Trigger**: Prevents UPDATE/DELETE on `audit_logs` table
2. **Checksum Verification** (Advanced):
   - Each audit log entry hashed (SHA-256)
   - Hash chain linking each entry to previous
   - Tampering breaks chain, detected on verification

**Example**:
```
Entry 1: hash(entry_1_data)
Entry 2: hash(entry_2_data + hash_of_entry_1)
Entry 3: hash(entry_3_data + hash_of_entry_2)
```

**Verification**: Recompute hash chain, compare with stored hashes → any mismatch = tampering detected

---

### Backup & Archival

- **Daily Backup**: Audit logs backed up to offsite storage
- **Long-Term Archival**: Logs older than 1 year moved to cold storage (S3 Glacier equivalent)
- **Retention**: 5 years minimum, financial logs 7 years
- **Restoration**: Can restore archived logs within 24 hours if needed for investigation

---

## Audit Alerts

**Real-Time Alerts** (for suspicious activity):

| Condition | Alert Recipient | Action |
|-----------|-----------------|--------|
| >10 failed logins from same IP in 1 hour | Ops team | Block IP, investigate |
| Super Admin permission granted | All Super Admins | Verify legitimacy |
| Anonymous complaint metadata accessed | DPO + Super Admins | Review justification |
| Large data export (>1000 records) | Ops team | Review necessity |
| Audit log query on sensitive data | Super Admin | Log the query itself |

---

## Non-Compliance Remediation

**If Audit Gaps Detected**:
1. Identify root cause (bug, misconfiguration, malicious deletion attempt)
2. Restore from backup if possible
3. Document incident
4. Fix underlying issue
5. Notify affected users (if required)

**If Unauthorized Modification Attempt Detected**:
1. Immediate account suspension
2. Full investigation
3. Legal action if malicious
4. System-wide audit

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
