# JOIN_REQUEST_WORKFLOW.md

## Document Purpose

This document defines the **membership join request and approval workflow** for Bangladesh Jatiotabadi Jubodal.

**Research-Based**: Aligned with real-world Schedule 1, Form 'Ka' submission process, Tk 10 membership fee, and hierarchical approval structure.

---

## Overview

**User Intent**: "Interested people may submit join requests, but approval is strictly hierarchical"

**Workflow**:
```
Public/Interested Person → Submit Application → Local Committee Reviews → Approval → User Account Created → Position Assigned → Dashboard Access Granted
```

---

## Scope

### This Document Controls

✅ **Join Request Submission**
- Public form (no login required)
- Required fields (name, phone, NID, address, desired jurisdiction)  
- Captcha & rate limiting
- Application fee payment tracking (optional)

✅ **Approval Workflow**
- Who can approve (local committee leaders)
- Hierarchical approval rules
- Approval criteria
- Rejection with reason

✅ **User Account Creation**
- Auto-create user account upon approval
- Initial state: `verified_at = NULL`
- Credentials generation (SMS/email)

✅ **Position Assignment & Verification**
- Committee leader assigns position
- `verified_at` timestamp set → Dashboard access granted

---

### This Document Does NOT Control

❌ **Database Schema** (Owned by DATABASE_SCHEMA.md)
- `join_requests` table structure (already defined)

❌ **User Authentication** (Owned by PHASE_C_AUTH.md)
- Login mechanism
- Password policies

❌ **Committee Management** (Owned by PHASE_D_COMMITTEE.md)
- Committee structure
- Position rules

---

## Join Request Submission

### Public Form Fields

**Endpoint**: `POST /api/v1/public/join-requests` (No authentication)

**Required Fields**:
```json
{
  "full_name": "আব্দুল করিম",
  "full_name_en": "Abdul Karim",
  "phone": "+8801712345678",
  "email": "karim@example.com",  // optional
  "nid": "1234567890123",  // National ID (13 digits)
  "date_of_birth": "1995-05-15",
  "address": "123 Main Street, Ward 5, Joypurhat",
  "address_bn": "১২৩ মেইন স্ট্রিট, ওয়ার্ড ৫, জয়পুরহাট",
  "jurisdiction_id": "uuid-of-joypurhat-ward-5",
  "why_join": "I want to contribute to the nationalist movement...",  // optional
  "application_fee_paid": false,  // Tk 10 (optional for digital system)
  "captcha_token": "recaptcha_token"
}
```

**Validation**:
- Phone: Must be +880 format (Bangladesh), unique (no existing user with same phone)
- NID: 10 or 13 digits, unique
- Age: Must be 18+ years old
- Jurisdiction: Must be ward or union level (lowest levels for new members)

**Rate Limiting**: 3 submissions per IP per 24 hours

**Response**:
```json
{
  "request_id": "uuid",
  "reference_number": "JR-2026-0001234",
  "status": "pending",
  "message": "Your join request has been submitted. Reference: JR-2026-0001234. You will be notified via SMS when reviewed."
}
```

---

### Reference Number Format

`JR-{YEAR}-{SEQUENTIAL_7_DIGITS}`

Example: `JR-2026-0001234`

Allows applicants to check status via public portal.

---

## Approval Workflow

### Who Can Review/Approve

**Approval Authority by Jurisdiction Level**:

| Applicant's Requested Jurisdiction | Reviewed By |
|------------------------------------|-------------|
| Ward (rural/urban) | Union/Municipality Committee Leaders |
| Union | Upazila Committee Leaders |
| Municipality Ward | Municipality Committee Leaders |

**Authorized Approvers**:
- President
- General Secretary
- Organizational Secretary

**Approval Threshold**: Any 1 of the above can approve

---

### Review Process

#### 1. Notification to Reviewers

When join request submitted → Notify committee leaders of that jurisdiction via:
- Dashboard notification
- SMS (if configured)
- Email

**Notification**:
```
New join request from Abdul Karim for Ward 5, Joypurhat.
Reference: JR-2026-0001234
Review now: [Dashboard Link]
```

---

#### 2. Review Dashboard

**Endpoint**: `GET /api/v1/join-requests?status=pending&jurisdiction={id}`

**Display**:
```
Join Requests Pending Approval (Ward 5, Joypurhat)

Reference: JR-2026-0001234
Name: আব্দুল করিম (Abdul Karim)
Phone: +8801712345678
Age: 29
Address: 123 Main Street, Ward 5, Joypurhat
Submitted: 2026-01-21 10:30 AM

[Approve] [Reject] [View Details]
```

**Details View**:
- Full application details
- NID verification status (manual check)
- Duplicate check (same phone/NID)
- "Why join" motivation statement

---

#### 3. Approval Decision

**Approve**:
```
PATCH /api/v1/join-requests/{id}
{
  "status": "approved",
  "reviewed_by": "reviewer_user_id"
}
```

**Actions on Approval**:
1. Create user account:
   ```sql
   INSERT INTO users (
     full_name, phone, email, password_hash, 
     is_active, verified_at, created_at
   ) VALUES (
     'আব্দুল করিম', '+8801712345678', 'karim@example.com',
     bcrypt_hash(random_password),  -- Initial random password
     TRUE, NULL,  -- verified_at is NULL initially
     NOW()
   );
   ```

2. Send credentials via SMS:
   ```
   Welcome to Bangladesh Jubo Dal!
   Your account: +8801712345678
   Initial password: Temp@1234
   Login: https://bjdms.arint.win
   Please change your password immediately.
   ```

3. Update join_request:
   ```sql
   UPDATE join_requests
   SET status = 'approved',
       reviewed_by = '{reviewer_id}',
       reviewed_at = NOW(),
       created_user_id = '{new_user_id}'
   WHERE id = '{request_id}';
   ```

4. Log event:
   ```
   action: join_request_approved
   entity: join_requests
   entity_id: {request_id}
   user_id: {reviewer_id}
   ```

**Notification to Applicant** (SMS):
```
Your join request (JR-2026-0001234) has been approved!
Login: +8801712345678
Password: Temp@1234
Change password: https://bjdms.arint.win/login
```

---

**Reject**:
```
PATCH /api/v1/join-requests/{id}
{
  "status": "rejected",
  "rejection_reason": "Applicant does not reside in requested jurisdiction",
  "reviewed_by": "reviewer_user_id"
}
```

**Notification to Applicant** (SMS):
```
Your join request (JR-2026-0001234) could not be approved.
Reason: Applicant does not reside in requested jurisdiction.
For inquiries, contact your local Jubo Dal office.
```

---

## Position Assignment & Dashboard Access

### State After Approval

**New User Account Created**:
- `is_active = TRUE`
- `verified_at = NULL` ← **No dashboard access yet**  
- Can login but redirected to "Awaiting Position Assignment" page

---

### Position Assignment by Committee

**Workflow**:
1. Committee leader logs in
2. Goes to "Members" → "Pending Verification"
3. Sees list of approved join requests (users with `verified_at = NULL`)
4. Assigns position to new member

**Endpoint**: `POST /api/v1/committees/{committee_id}/members`
```json
{
  "user_id": "{new_user_id}",
  "position_id": "{position_id}",  // e.g., "Member" position
  "started_at": "2026-01-21"
}
```

**Actions**:
1. Create `committee_members` record
2. Update user: `verified_at = NOW()`  
3. Trigger permission refresh
4. Send notification to new member

**Notification** (SMS):
```
Congratulations! You have been assigned as Member of Ward 5 Committee, Joypurhat.
You now have full dashboard access.
Login: https://bjdms.arint.win
```

---

### Dashboard Access Gate

**Login Check**:
```javascript
if (user.is_active && user.verified_at !== null) {
  // Grant dashboard access
  redirect_to('/dashboard');
} else if (user.is_active && user.verified_at === null) {
  // Awaiting position assignment
  redirect_to('/pending-verification');
} else {
  // Inactive account
  redirect_to('/account-suspended');
}
```

**Pending Verification Page**:
```
Your account is awaiting committee approval.

Status: Account Created ✅
Next Step: Position Assignment (Pending)

You will receive an SMS notification once your position is assigned.
Estimated time: 1-3 days.

For inquiries, contact:
Ward 5 Committee, Joypurhat
Phone: +880XXXXXXXXXX
```

---

## Business Logic Rules

### Rule 1: Jurisdiction Level Restriction

New join requests can ONLY apply to:
- Ward level (rural or urban)
- Union level

**Rationale**: New members start at grassroots. District/Central positions filled through internal promotion.

**Validation**:
```
IF requested_jurisdiction.level NOT IN ('ward', 'union'):
  REJECT with error "Join requests are only accepted for ward/union level"
```

---

### Rule 2: Age Requirement

Applicant must be 18+ years old at time of submission.

**Validation**:
```
age = CURRENT_DATE - date_of_birth
IF age < 18 years:
  REJECT with error "Applicant must be at least 18 years old"
```

---

### Rule 3: Duplicate Prevention

**Phone Number**: Unique check  
**NID**: Unique check

**Validation**:
```
IF EXISTS user WITH same phone:
  REJECT "An account with this phone number already exists"

IF EXISTS user WITH same NID:
  REJECT "An account with this NID already exists"
```

**Edge Case**: Same person applies twice (e.g., forgot first application)
- Show existing application status if pending
- Allow resubmission if previous was rejected

---

### Rule 4: Auto-Cleanup

**Stale Requests**: Pending for > 90 days → Auto-reject

```
UPDATE join_requests
SET status = 'rejected',
    rejection_reason = 'Application expired after 90 days of inactivity'
WHERE status = 'pending'
  AND submitted_at < NOW() - INTERVAL '90 days';
```

**Notification**: SMS to applicant informing them to reapply

---

## Additional Features

### Status Check Portal

**Public Endpoint**: `GET /api/v1/public/join-requests/status?reference={JR-2026-0001234}`

**Response**:
```json
{
  "reference_number": "JR-2026-0001234",
  "status": "pending",  // or approved, rejected
  "submitted_at": "2026-01-21T10:30:00Z",
  "message": "Your application is under review by Ward 5 Committee."
}
```

**If approved**:
```json
{
  "reference_number": "JR-2026-0001234",
  "status": "approved",
  "approved_at": "2026-01-23T14:15:00Z",
  "message": "Your application has been approved. Check your phone for login credentials."
}
```

**If rejected**:
```json
{
  "reference_number": "JR-2026-0001234",
  "status": "rejected",
  "rejection_reason": "Applicant does not reside in requested jurisdiction",
  "message": "Your application could not be approved. Contact your local office for details."
}
```

---

### Membership Fee Tracking (Optional)

**If Tk 10 fee is collected**:

1. Applicant pays at local office (cash)
2. Local treasurer marks `application_fee_paid = TRUE` + `payment_reference`
3. Receipt issued (similar to donation receipt)

**Digital Enhancement** (Future):
- Online payment via bKash/Nagad
- Auto-approval upon payment confirmation
- Digital receipt

---

## Integration Points

### With Phase C (Auth)
- User account creation
- Initial password generation
- `verified_at` controls dashboard access

### With Phase D (Committee)
- Position assignment workflow
- Committee members table
- Permission updates

### With Notification System
- SMS notifications for status updates
- Dashboard notifications for reviewers

---

## Testing Requirements

### Submission Tests
- Submit valid join request
- Reject underage applicant (<18)
- Reject duplicate phone number
- Captcha validation
- Rate limiting enforcement

### Approval Tests
- Approve join request → User created
- Reject join request → Notification sent
- Auto-cleanup stale requests (>90 days)

### Position Assignment Tests
- Assign position → `verified_at` set
- Dashboard access granted after assignment
- Permission refresh triggered

---

## Future Enhancements

### Online Payment Integration

- bKash/Nagad payment during submission
- Auto-approve upon payment
- Digital membership card generation

### Biometric Verification

- NID verification via government API (Bangladesh)
- Photo upload for ID verification
- Facial recognition (future)

### Bulk Onboarding

- Committee uploads CSV of new members (e.g., after offline registration drive)
- Bulk user creation
- Bulk SMS credential distribution

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
