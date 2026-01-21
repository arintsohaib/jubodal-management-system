# Phase A — Database Schema (PostgreSQL)

## Purpose

This document defines the **single source of truth for all data structures** used in the Bangladesh Jubo Dal Management System (BJDMS).

Any AI agent or developer must:

* Read this file before writing backend code
* Never change implemented tables without migration
* Keep backward compatibility

---

## Design Principles

1. PostgreSQL 16+
2. UUID as primary keys
3. Soft delete by default (`deleted_at`)
4. Audit timestamps everywhere
5. Hierarchy-first design

---

## Core Reference Tables

### 1. `jurisdiction_levels`

Defines hierarchy levels matching Bangladesh administrative structure.

| column | type        | notes                                                    |
| ------ | ----------- | -------------------------------------------------------- |
| id     | smallint PK | fixed values                                             |
| name   | varchar     | Central, Division, District, Upazila, Municipality, Union, Ward |
| type   | varchar     | national, regional, district, sub-district, local        |

**Hierarchy**:
- Central (National)
- Division (Regional - 8 divisions, optional oversight)
- District (64 zilas)
- Upazila (492-500 sub-districts, **rural**)
- Municipality (329-330 pourashavas, **urban**)
- Union (4,550+ union parishads, **rural**)
- Ward (Lowest unit, both rural and urban)

---

### 2. `jurisdictions`

Represents actual geographic/political units in Bangladesh.

| column            | type        | notes                                     |
| ----------------- | ----------- | ----------------------------------------- |
| id                | uuid PK     |                                           |
| level_id          | smallint FK | jurisdiction_levels                       |
| parent_id         | uuid FK     | self-reference                            |
| name              | varchar     | Joypurhat, Dhaka, etc                     |
| name_bn           | varchar     | Bengali name (জয়পুরহাট)                  |
| code              | varchar     | optional code                             |
| is_urban          | boolean     | true for Municipality/City, false for rural|
| population        | integer     | optional population data                  |
| created_at        | timestamp   |                                           |
| updated_at        | timestamp   |                                           |

**Notes**:
- `is_urban`: Distinguishes Municipality/City Corporation (urban) from Upazila/Union (rural)
- `name_bn`: Bengali name for localization
- `parent_id`: NULL for Central, points to parent jurisdiction for all others

---

## User & Identity

### 3. `users`

| column        | type      | notes    |
| ------------- | --------- | -------- |
| id            | uuid PK   |          |
| full_name     | varchar   |          |
| phone         | varchar   | unique   |
| email         | varchar   | nullable |
| password_hash | text      |          |
| is_active     | boolean   |          |
| verified_at   | timestamp |          |
| created_at    | timestamp |          |
| updated_at    | timestamp |          |
| deleted_at    | timestamp |          |

---

## Committee Structure

### 4. `committee_types`

| id | name               |
| -- | ------------------ |
| 1  | Full Committee     |
| 2  | Convener Committee |

---

### 5. `positions`

Defines all possible positions in Bangladesh Jatiotabadi Jubodal committees.

| id  | name                          | rank | committee_type  | notes                                    |
| --- | ----------------------------- | ---- | --------------- | ---------------------------------------- |
| 1   | President                     | 1    | Full            | Top leader                               |
| 2   | Senior Vice President         | 2    | Full            | Senior leadership                        |
| 3   | Vice President                | 3    | Full            | Multiple vice presidents possible        |
| 4   | General Secretary             | 10   | Full            | Chief administrator                      |
| 5   | Joint General Secretary       | 11   | Full            | Assists General Secretary                |
| 6   | Organizational Secretary      | 15   | Full            | Internal organization                    |
| 7   | Joint Secretary               | 16   | Full            | Various departments                      |
| 8   | Assistant Organizing Secretary| 17   | Full            | Assists organizational work              |
| 9   | Publicity Secretary           | 18   | Full            | Media and communications                 |
| 10  | Office Secretary              | 19   | Full            | Administrative functions                 |
| 11  | Treasurer                     | 20   | Full            | Financial management                     |
| 12  | Deputy Treasurer              | 21   | Full            | Assists treasurer                        |
| 13  | Member                        | 50   | Full            | General committee member                 |
| 14  | Convener                      | 1    | Convener        | Equivalent to President (interim)        |
| 15  | Member Secretary              | 10   | Convener        | Equivalent to General Secretary (interim)|
| 16  | Joint Convener                | 2    | Convener        | Assists convener                         |

**Position Equivalence Rules** (for Convener → Full Committee transition):
- Convener (rank 1) = President (rank 1)
- Member Secretary (rank 10) = General Secretary (rank 10)
- Joint Convener may become Vice President or Senior Vice President

**Notes**:
- Lower rank numbers indicate higher authority
- Ranks 1-20: Leadership positions
- Rank 50+: General members
- Union-level committees can have up to 71 members total

---

### 6. `committees`

| column            | type    | notes           |
| ----------------- | ------- | --------------- |
| id                | uuid PK |                 |
| committee_type_id | int FK  | committee_types |
| jurisdiction_id   | uuid FK |                 |
| is_active         | boolean |                 |
| formed_at         | date    |                 |
| dissolved_at      | date    | nullable        |

---

### 7. `committee_members`

| column       | type    | notes     |
| ------------ | ------- | --------- |
| id           | uuid PK |           |
| committee_id | uuid FK |           |
| user_id      | uuid FK |           |
| position_id  | int FK  | positions |
| started_at   | date    |           |
| ended_at     | date    | nullable  |

---

## Permission System

### 8. `roles`

| id | name            |
| -- | --------------- |
| 1  | Super Admin     |
| 2  | Central Leader  |
| 3  | District Leader |
| 4  | Unit Leader     |

---

### 9. `permissions`

| id | key             | description |
| -- | --------------- | ----------- |
| 1  | committee.read  |             |
| 2  | committee.write |             |
| 3  | complaint.read  |             |
| 4  | complaint.write |             |

---

### 10. `role_permissions`

| role_id | permission_id |

---

### 11. `user_roles`

| user_id | role_id | jurisdiction_id |

---

## Activity & Events

### 12. `activities`

| column          | type      |
| --------------- | --------- |
| id              | uuid PK   |
| user_id         | uuid FK   |
| jurisdiction_id | uuid FK   |
| title           | varchar   |
| description     | text      |
| created_at      | timestamp |

---

## Complaint System

### 13. `complaints`

| column          | type      | notes                         |
| --------------- | --------- | ----------------------------- |
| id              | uuid PK   |                               |
| submitted_by    | uuid FK   | nullable                      |
| target_user_id  | uuid FK   |                               |
| jurisdiction_id | uuid FK   |                               |
| is_anonymous    | boolean   |                               |
| description     | text      |                               |
| status          | varchar   | received/review/action/closed |
| created_at      | timestamp |                               |

---

---

### 14. `complaint_evidence`

| column       | type    | notes |
| ------------ | ------- | ----- |
| id           | uuid PK |       |
| complaint_id | uuid FK |       |
| file_path    | text    |       |
| uploaded_at  | timestamp |     |

---

## Join Request System

### 15. `join_requests`

Tracks membership application requests (research confirmed: Form 'Ka', Tk 10 fee).

| column            | type      | notes                                          |
| ----------------- | --------- | ---------------------------------------------- |
| id                | uuid PK   |                                                |
| full_name         | varchar   | Applicant's full name                          |
| full_name_bn      | varchar   | Name in Bengali                                |
| phone             | varchar   | Contact number (+880 format)                   |
| email             | varchar   | nullable                                       |
| nid               | varchar   | National ID number                             |
| date_of_birth     | date      |                                                |
| address           | text      | Current address                                |
| address_bn        | text      | Address in Bengali                             |
| jurisdiction_id   | uuid FK   | Requested jurisdiction (ward/union level)      |
| application_fee_paid | boolean | Tk 10 paid (optional for digital)           |
| payment_reference | varchar  | nullable                                       |
| status            | varchar   | pending/approved/rejected                      |
| submitted_at      | timestamp |                                                |
| reviewed_by       | uuid FK   | nullable, references users (approver)          |
| reviewed_at       | timestamp | nullable                                       |
| rejection_reason  | text      | nullable                                       |
| created_user_id   | uuid FK   | nullable, user created upon approval           |

**Workflow**:
1. Public submits join request (no login required)
2. Local committee (ward/union) reviews application
3. Approval creates user account with `verified_at = NULL`
4. Committee assigns position → `verified_at` set → dashboard access granted

**Notes**:
- Matches real-world Schedule 1, Form 'Ka' submission
- Phone must match +880 Bangladesh format
- Approval is hierarchical (reviewed by local committee leaders)

---

## Audit Logs

### 16. `audit_logs`

| column     | type      | notes |
| ---------- | --------- | ----- |
| id         | uuid PK   |       |
| user_id    | uuid FK   | nullable for system actions |
| action     | varchar   | login, create_committee, assign_position, etc |
| entity     | varchar   | users, committees, complaints, join_requests |
| entity_id  | uuid      |       |
| ip_address | varchar   | optional |
| user_agent | text      | optional browser/device info |
| created_at | timestamp |       |

---

## Rules

1. Never delete records permanently
2. All foreign keys enforced
3. Index FK columns
4. Migrations mandatory for changes

---

## Status

Phase A: 🟢 Implemented (Schema Defined)
