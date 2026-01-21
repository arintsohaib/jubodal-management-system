# API_CONTRACT.md

## Purpose

This document defines the **single source of truth** for all backend APIs.
No AI agent, developer, or service may invent endpoints outside this file.

If this file stays consistent, backend, frontend, and mobile clients will never conflict.

---

## Global Rules

* Base URL (dev): [https://arint.win](https://arint.win)
* All APIs are versioned: `/api/v1/`
* Authentication: JWT (Bearer Token)
* All responses are JSON
* All timestamps are UTC ISO-8601

---

## Authentication APIs

### POST /api/v1/auth/login

Authenticate a user.

Request:

```json
{
  "phone": "+8801XXXXXXXXX",
  "password": "string"
}
```

Response:

```json
{
  "access_token": "jwt",
  "refresh_token": "jwt",
  "expires_in": 3600
}
```

---

### POST /api/v1/auth/refresh

Refresh access token.

---

## User APIs

### GET /api/v1/users/me

Returns authenticated user profile.

---

## Committee APIs

### GET /api/v1/committees

List committees by level & region.

Query Params:

* level: central|district|upazila|union|ward
* region_code

---

### POST /api/v1/committees

Create committee (authorized roles only).

---

### GET /api/v1/committees/{id}

Get committee details.

---

## Committee Member APIs

### POST /api/v1/committees/{id}/members

Add member to committee.

---

### PATCH /api/v1/committee-members/{id}

Update member role or status.

---

## Complaint APIs

### POST /api/v1/complaints

Submit a complaint.

---

### GET /api/v1/complaints

List complaints (role-based visibility).

---

## Audit APIs

### GET /api/v1/audit-logs

System-wide audit logs (super admin only).

---

## Error Response Format

```json
{
  "error": true,
  "code": "STRING_CODE",
  "message": "Human readable message"
}
```

---

## Final Note

If this document is broken, the system will break.  
If this document remains intact, the system will remain stable — no matter how large it grows.

⚠️ This file is the single source of truth.  
⚠️ No major change is allowed without updating this document first.
