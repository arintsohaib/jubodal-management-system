# Bangladesh Jubo Dal Management System (BJDMS)

## Core Design Principles (Must Never Break)

1. **Modular First** – প্রতিটি ফিচার আলাদা মডিউল
2. **Hierarchy Aware** – Central → Division → District → Upazila/Municipality → Union/Ward
3. **Permission Driven** – Role না থাকলে Action নেই
4. **Verification Required** – Committee ছাড়া কেউ core dashboard পাবে না
5. **Audit Everything** – সব action লগ হবে
6. **Backward Compatible** – নতুন ফিচার পুরনো কিছু ভাঙবে না
7. **Urban/Rural Aware** – Municipality (urban) and Upazila/Union (rural) paths supported

---

## Technology Stack (Stable 2026 Standard)

### Frontend

* Next.js 16 (App Router)
* React 19
* TypeScript
* Tailwind + Shadcn
* Real-time via WebSocket / Server Actions

### Backend

* Go 1.25 (Primary Core Services)
* Python 3.14 (FastAPI) – Analytics / AI / Reports

### Data Layer

* PostgreSQL 18 – Primary DB
* Redis 8.4 – Cache / Permission / Session
* OpenSearch 3.4 – Search Engine
* S3-compatible Storage – Files & Proofs

### Security

* OAuth2 + JWT + Refresh Token
* RBAC + ABAC
* Full Audit Log

⚠️ Stack change only via Central Architecture Decision

---

## System Status Tracking Convention

Every module must be marked as one:

* 🟢 Implemented
* 🟡 In Progress
* 🔴 Planned

---

## Phase-0: Foundation (Mandatory First)

### 0.1 Auth & Identity Service 🟢

* Login only (No public signup)
* Admin/Invite based onboarding
* Token rotation

### 0.2 Role & Permission Engine 🟢

* Central roles
* District roles
* Committee type aware (Full / Convener)

---

## Phase-1: Core Management System 🟢

### 1.1 User Profile Module 🟢

* Name
* Position
* Committee Type
* Jurisdiction Level
* Verification Status

### 1.2 Committee Management 🟢

* Full Committee
* Convener Committee
* Auto-equivalence logic:

  * Convener == President
  * Member Secretary == General Secretary

### 1.3 Dashboard System 🟢

* Role based widgets
* Jurisdiction filtered data

---

## Phase-2: Activity, Task & Event 🟢

### 2.1 Activity Log 🟢

* Daily activity submission
* Proof upload
* Approval workflow

### 2.2 Task Management 🟢

* Assign task
* Deadline
* Progress tracking

### 2.3 Event Management 🟢

* Event creation
* Attendance
* Reports

---

## Phase-3: Complaint & Whistleblower System 🟢

### 3.1 Complaint Intake 🟢

* Named complaint
* Anonymous complaint
* Evidence upload

### 3.2 Complaint Routing 🟢

* Default → District Leaders
* Optional → Central Leaders

### 3.3 Complaint Lifecycle 🟢

* Received
* Under Review
* Action Taken
* Closed

⚠️ Anonymous complaints must preserve metadata securely

---

## Phase-4: Search & Transparency 🟢

### 4.1 Global Search 🟢

* Name
* Position
* Area

### 4.2 Audit & Logs 🟢

* Who did what
* When
* From where

---

## Phase-5: Donation & Finance 🟢

### 5.1 Donation Module 🟢

* District-wise fund
* Purpose tagging
* Transparent ledger

### 5.2 Reports 🟢

* Monthly
* Yearly
* Audit ready

---

## Phase-6: AI & Advanced Layer 🟢

* Activity heatmap
* Performance analytics
* Anomaly detection
* AI assistant (read-only)

---

## Phase-7: Notifications & Hardening 🟢

* Real-time WebSocket notifications
* SMS Gateway Integration (Mediator)
* Production Infrastructure Hardening
* Scalability Optimization

---

## Rules for Adding New Features (VERY IMPORTANT)

1. New feature MUST belong to a Phase
2. Must not change existing DB schema without migration
3. Must update this file BEFORE coding
4. Must mark status correctly
5. Must document permissions impact

---

## AI Agent Instruction Block (READ THIS FIRST)

When an AI agent is instructed:

1. Read this file fully
2. Identify current phase
3. Check implemented modules
4. Add new feature without touching implemented ones
5. Update phase status

❌ Never refactor everything together
❌ Never bypass permission engine

---

## Ownership

System Owner: Grayhawk Sentinel Ltd.
Architecture Authority: Grayhawk Sentinel Ltd.
Official Website: https://grayhawks.com

---

If this document is broken, the system will break.  
If this document remains intact, the system will remain stable — no matter how large it grows.

⚠️ This file is the single source of truth.
⚠️ No major change is allowed without updating this document first.

