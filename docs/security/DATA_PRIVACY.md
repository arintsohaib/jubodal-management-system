# DATA_PRIVACY.md

## Purpose

This document defines **data privacy policies and compliance requirements** for Bangladesh Jatiotabadi Jubodal Management System (BJDMS).

**Context**: Bangladesh does not have comprehensive GDPR-equivalent legislation (as of 2026), but this document establishes privacy best practices aligned with international standards.

---

## Privacy Principles

1. **Lawfulness, Fairness, Transparency** - Clear communication about data usage
2. **Purpose Limitation** - Data used only for stated purposes
3. **Data Minimization** - Collect only necessary data
4. **Accuracy** - Keep data accurate and up-to-date
5. **Storage Limitation** - Delete data when no longer needed
6. **Integrity & Confidentiality** - Protect data from unauthorized access
7. **Accountability** - BJDMS responsible for compliance

---

## Personal Data Classification

### Highly Sensitive Data

**Definition**: Data that, if compromised, poses significant risk to individual safety or privacy.

**Examples**:
- Anonymous complaint submitter identity (IP address, metadata)
- National ID (NID) numbers
- Password hashes
- Financial donation records (if anonymous)

**Protection Measures**:
- Encrypted at rest (optional: PGP encryption)
- Access restricted to Super Admin only
- Audit logging for all access
- Deleted after retention period

---

### Sensitive Personal Data

**Definition**: Personal data that requires heightened protection.

**Examples**:
- Phone numbers
- Email addresses
- Home addresses
- Committee positions
- Financial transactions (non-anonymous)

**Protection Measures**:
- HTTPS for transmission
- Role-based access control
- Masked in public views (phone: +880171***5678)
- User can opt out of public directory

---

### Public Data

**Definition**: Data intended for organizational transparency.

**Examples**:
- Committee leader names (if opted in to public directory)
- Committee structure (positions, jurisdiction)
- Public events
- Approved activities (if marked public)

**Protection Measures**:
- User consent required for public visibility
- Can be searched without authentication

---

## Data Collection

### What We Collect

| Data Type | Purpose | Legal Basis | Retention |
|-----------|---------|-------------|-----------|
| Name, phone, NID | User identification, membership | Organizational membership | Until account deletion + 2 years |
| Email | Communication, notifications | User consent | Until account deletion |
| Address | Jurisdiction assignment | Organizational operations | Until account deletion + 2 years |
| Activity logs | Organizational accountability | Legitimate organizational interest | 2 years |
| Financial data | Donation tracking, budgeting | Legal/financial compliance | 7 years (audit requirement) |
| Complaint data | Organizational integrity | Whistleblower protection | Indefinite (unless requested deletion) |
| Audit logs | Security, accountability | Legal compliance | 5 years |

---

### How We Collect

1. **Join Request Form**: Name, phone, NID, address, date of birth
2. **User Profile**: Email, photo (optional), biography (optional)
3. **Activity Submission**: Activity details, proof uploads
4. **Complaint Submission**: Complaint text, evidence files, submitter info (if named)
5. **Donation Recording**: Donor name, amount, payment method
6. **Automated**: Login timestamps, IP addresses (hashed for anonymous complaints), user agent

---

## User Rights

### 1. Right to Access

**What**: User can request all personal data held about them.

**How to Exercise**:
- Log in to dashboard → Settings → "Download My Data"
- Generates JSON export of all user data
- Delivered within 30 days

**Includes**:
- Profile information
- Committee memberships
- Activities submitted
- Complaints filed (named only, not anonymous)
- Audit log of user's actions

---

### 2. Right to Rectification

**What**: User can correct inaccurate data.

**How to Exercise**:
- Edit profile directly in dashboard
- Request correction via committee leader (for NID, phone changes)
- Verified within 14 days

**Limitations**: Some fields (NID) require manual verification to prevent fraud.

---

### 3. Right to Erasure ("Right to be Forgotten")

**What**: User can request account deletion.

**How to Exercise**:
- Dashboard → Settings → "Delete My Account"
- Confirmation required (SMS OTP)
- Processed within 30 days

**What Gets Deleted**:
- User profile (name, phone, email, address)
- Login credentials
- Optional: Activity submissions (user's choice)

**What Is Retained** (Legal/Organizational Necessity):
- Financial transactions (7 years for audit)
- Audit logs (5 years for accountability)
- Committee membership history (organizational record)
- Anonymous complaints (no PII anyway)

**Soft Delete**: Account marked `deleted_at`, data retained for 2 years, then purged.

---

### 4. Right to Data Portability

**What**: User can export data in machine-readable format.

**How to Exercise**: Same as "Right to Access" (JSON export).

**Format**: JSON (includes all user data with schema documentation).

---

### 5. Right to Object

**What**: User can object to data processing for specific purposes.

**How to Exercise**:
- Opt out of public directory: Profile → Privacy Settings → "Hide from public directory"
- Opt out of non-essential communications: Notification Preferences

**Limitations**: Cannot opt out of:
- Organizational accountability (activity logs)
- Financial record-keeping (legal requirement)
- Audit logging (security requirement)

---

### 6. Right to Restrict Processing

**What**: User can request temporary halt of data processing.

**How to Exercise**: Contact BJDMS data protection officer.

**Example**: User disputes data accuracy → Processing restricted until verified.

---

## Special Data Categories

### Anonymous Complaint Submitter Data

**Heightened Protection**:
- IP address immediately hashed (SHA-256) upon submission
- Metadata stored in separate restricted table
- Only Super Admin can access
- Metadata deleted after 90 days
- No correlation possible with other user data

**User Rights Exception**:
- No "right to access" for anonymous complaints (would de-anonymize)
- No "right to erasure" (data already minimized/hashed)

---

### Children's Data

**Policy**: BJDMS does not knowingly collect data from individuals under 18.

**Age Verification**: Date of birth required, must be 18+ to join.

**If Discovered**: Immediate account deletion, data purged.

---

## Data Sharing

### Internal Sharing

**Within BJDMS**:
- Committee leaders can view members in their jurisdiction (RBAC enforced)
- Central leaders can view all data (organizational necessity)
- Treasurer can view financial data (role requirement)

**No Sharing Without Need-to-Know**: General members cannot see other members' phone numbers, addresses (unless public directory opt-in).

---

### External Sharing

**BJDMS does NOT share data with third parties, except**:

1. **Legal Compliance**: If required by Bangladesh law (court order, law enforcement request)
   - Legal basis documented
   - User notified (unless legally prohibited)

2. **Service Providers** (Minimal, only if used):
   - SMS gateway (phone numbers for notifications)
   - Email provider (email addresses for transactional emails)
   - **Data Processing Agreement** required
   - Service providers bound by confidentiality

3. **Public Disclosure**: Data marked as public (user consent obtained)

**No Marketing/Commercial Use**: User data NEVER sold or used for commercial purposes.

---

## Data Retention

| Data Type | Retention Period | Reason | After Retention |
|-----------|------------------|--------|-----------------|
| User profiles (active) | Until account deletion | Operational | Soft delete → 2 years → Purge |
| Activity logs | 2 years | Accountability | Archived to cold storage |
| Financial records | 7 years | Legal/audit requirement | Archived, never deleted |
| Audit logs | 5 years | Security/accountability | Archived |
| Join requests (rejected) | 1 year | Operational | Deleted |
| Anonymous complaint metadata | 90 days | Abuse detection | Deleted |
| Complaint data (named) | Indefinite | Organizational integrity | Unless user requests deletion |

---

## Data Security Measures

(See SECURITY_POLICY.md for full details)

**Summary**:
- TLS 1.3 encryption in transit
- Bcrypt password hashing (cost 12)
- Optional PGP encryption for highly sensitive fields
- Access control via RBAC/ABAC
- Audit logging for all data access
- Regular security audits

---

## Privacy by Design

### Features Built for Privacy

1. **Anonymous Complaint System**: No PII collected for anonymous complaints
2. **Public Directory Opt-In**: Users choose if their info is publicly searchable
3. **Phone Number Masking**: Public views show +880171***5678
4. **Soft Delete**: Accounts marked deleted, data removed after grace period
5. **Data Minimization**: Only collect necessary fields
6. **Granular Permissions**: Access to data based on role and jurisdiction

---

## User Consent

### Explicit Consent Required For:

- Joining BJDMS (consent to membership and data processing)
- Public directory inclusion (opt-in checkbox)
- Non-essential notifications (can opt out anytime)
- Photo upload (optional, user-controlled)

### Implied Consent:

- Activity logging (necessary for organizational accountability)
- Audit logging (necessary for security)
- Financial record-keeping (legal requirement)

**Consent Withdrawal**: User can withdraw consent anytime (may affect service availability).

---

## Cookies & Tracking

**Cookies Used**:

| Cookie Name | Purpose | Duration | Type |
|-------------|---------|----------|------|
| `access_token` | Authentication | 15 minutes | Essential |
| `refresh_token` | Token refresh | 90 days | Essential |
| `session_id` | Session tracking | 7 days | Essential |

**No Third-Party Tracking**: No Google Analytics, Facebook Pixel, or other tracking scripts (unless explicitly documented and user-consented).

**Do Not Track**: BJDMS respects browser DNT (Do Not Track) signals.

---

## Data Breach Response

(See SECURITY_POLICY.md for incident response plan)

**User Notification**:
- If personal data breach occurs → Notify affected users within 72 hours
- Notification via SMS + email
- Details: What data affected, when occurred, what actions taken, user remediation steps

**Authority Notification**: If legally required in Bangladesh (not currently mandated as of 2026).

---

## Privacy Governance

### Data Protection Officer (DPO)

**Role**: Oversee data protection compliance, handle user requests, privacy audits.

**Contact**: privacy@bjd.org.bd

**Responsibilities**:
- Review data collection practices
- Handle user rights requests (access, erasure, rectification)
- Conduct privacy impact assessments
- Train committee leaders on privacy requirements

---

### Privacy Impact Assessment (PIA)

**Required For**:
- New features collecting personal data
- Changes to data retention policies
- Integration with third-party services

**Process**:
1. Identify data collected
2. Assess necessity and proportionality
3. Evaluate risks to user privacy
4. Define mitigation measures
5. Document and approve

---

## Compliance with Bangladesh Laws

**Relevant Legislation**:
- Bangladesh ICT Act, 2006 (amended 2013)
- Bangladesh Penal Code (data protection sections)
- Bangladesh Constitution (Right to Privacy - implied)

**BJDMS Compliance**:
- Data stored within Bangladesh (server located at arint.win)
- Cooperation with law enforcement (legal requests honored)
- No unlawful data processing

---

## International Data Transfers

**Current Policy**: All data stored in Bangladesh (arint.win server).

**If Future International Transfer**:
- Contract with data processor (Data Processing Agreement)
- Ensure adequate protection in destination country
- User notification and consent

---

## Children's Privacy (COPPA/GDPR-equivalent)

**Age Verification**: Must be 18+ to join BJDMS.

**No Children's Data**: If discovered that user is <18:
1. Account immediately suspended
2. Data deleted within 7 days
3. Parent/guardian notified (if contact info available)

---

## User Education

**Privacy Policy Page**: Publicly accessible, plain language explanation of data practices.

**Transparency Reports**: Annual report on:
- Number of user data requests (access, erasure)
- Law enforcement requests received
- Data breaches (if any)

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
