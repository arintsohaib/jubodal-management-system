# PHASE_G_SEARCH.md

## Phase Identification

- **Phase**: G
- **Name**: Search & Transparency
- **Status**: 🔴 Planned
- **Dependencies**:
  - Phase A (DATABASE_SCHEMA.md) - All searchable tables
  - Phase B (API_CONTRACT.md) - Search API endpoints
  - Phase C (PHASE_C_AUTH.md) - Permission-based search filtering
  - DOCKER_INFRASTRUCTURE.md - OpenSearch container

---

## Purpose

This document defines the **global search and transparency system** for Bangladesh Jatiotabadi Jubodal.

**Core Functions**:
- Search members by name, position, area (Bangla & English)
- Public directory of committee leaders
- Search activities, events, complaints (permission-filtered)
- Audit log search for transparency
- Advanced filtering and faceted search

---

## Scope Boundaries

### This Phase Controls

✅ **Search Infrastructure**
- OpenSearch cluster configuration
- Index mappings and analyzers
- Bangla language support (ইংরেজি + বাংলা)
- Synonym management

✅ **Indexing Pipeline**
- Real-time data synchronization from PostgreSQL
- Debezium CDC (Change Data Capture) or manual triggers
- Index refresh strategy
- Bulk indexing for historical data

✅ **Search API**
- Query DSL implementation
- Autocomplete/typeahead
- Fuzzy search for misspellings
- Faceted search (filter by jurisdiction, position, committee type)
- Result ranking algorithm

✅ **Public Directory**
- Publicly searchable committee leaders
- Privacy controls (phone number visibility)
- Organizational chart view

✅ **Audit Search**
- Search audit logs by user, action, entity, date range
- Super Admin only access
- Export capabilities

---

### This Phase Does NOT Control

❌ **Data Structures** (Owned by DATABASE_SCHEMA.md)
- Source tables being indexed

❌ **Permissions** (Owned by PHASE_C_AUTH.md)
- Who can search what
- Permission resolution

❌ **Infrastructure** (Owned by DOCKER_INFRASTRUCTURE.md)
- OpenSearch Docker image
- Network configuration

---

## Implementation Checklist

### OpenSearch Setup
- [ ] OpenSearch Docker container configured
- [ ] Bangla analyzer plugin installed
- [ ] Index templates created
- [ ] Synonym lists loaded (Bangla/English)

### Indexing
- [ ] Users index
- [ ] Committees index
- [ ] Activities index
- [ ] Events index
- [ ] Complaints index (permission-filtered)
- [ ] Audit logs index

### Search API
- [ ] Global search endpoint
- [ ] Autocomplete endpoint
- [ ] Advanced search with filters
- [ ] Export search results (CSV)

### Public Directory
- [ ] Public committee leader search
- [ ] Organizational chart API
- [ ] Privacy setting enforcement

---

## Technical Details

### OpenSearch Cluster Configuration

**Docker Compose** (to be added to DOCKER_INFRASTRUCTURE.md):
```yaml
opensearch:
  image: opensearchproject/opensearch:3.4.0
  environment:
    - discovery.type=single-node
    - OPENSEARCH_JAVA_OPTS=-Xms512m -Xmx512m
    - plugins.security.disabled=true  # Use app-level auth
  volumes:
    - opensearch_data:/usr/share/opensearch/data
  networks:
    - bjdms_network
```

**Scaling**: Single node for development, 3-node cluster for production

---

### Index Mappings

#### 1. Users Index

```json
{
  "settings": {
    "index": {
      "number_of_shards": 3,
      "number_of_replicas": 1
    },
    "analysis": {
      "analyzer": {
        "bangla_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "bengali_normalization", "bengali_stem"]
        },
        "autocomplete_analyzer": {
          "type": "custom",
          "tokenizer": "edge_ngram_tokenizer",
          "filter": ["lowercase"]
        }
      },
      "tokenizer": {
        "edge_ngram_tokenizer": {
          "type": "edge_ngram",
          "min_gram": 2,
          "max_gram": 10,
          "token_chars": ["letter", "digit"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {"type": "keyword"},
      "full_name": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "autocomplete": {"type": "text", "analyzer": "autocomplete_analyzer"},
          "keyword": {"type": "keyword"}
        }
      },
      "full_name_bn": {
        "type": "text",
        "analyzer": "bangla_analyzer",
        "fields": {
          "autocomplete": {"type": "text", "analyzer": "autocomplete_analyzer"}
        }
      },
      "phone": {"type": "keyword"},
      "jurisdiction_id": {"type": "keyword"},
      "jurisdiction_name": {"type": "text"},
      "jurisdiction_name_bn": {"type": "text", "analyzer": "bangla_analyzer"},
      "jurisdiction_level": {"type": "keyword"},
      "position_name": {"type": "keyword"},
      "position_name_bn": {"type": "text", "analyzer": "bangla_analyzer"},
      "committee_type": {"type": "keyword"},
      "is_active": {"type": "boolean"},
      "verified_at": {"type": "date"},
      "is_public": {"type": "boolean"},
      "created_at": {"type": "date"}
    }
  }
}
```

**Key Features**:
- Bangla analyzer for proper Bengali text search
- Autocomplete on names (both English & Bangla)
- `is_public` flag for public directory visibility

---

#### 2. Committees Index

```json
{
  "mappings": {
    "properties": {
      "id": {"type": "keyword"},
      "committee_type": {"type": "keyword"},
      "jurisdiction_id": {"type": "keyword"},
      "jurisdiction_name": {"type": "text"},
      "jurisdiction_name_bn": {"type": "text", "analyzer": "bangla_analyzer"},
      "jurisdiction_level": {"type": "keyword"},
      "is_active": {"type": "boolean"},
      "formed_at": {"type": "date"},
      "dissolved_at": {"type": "date"},
      "member_count": {"type": "integer"},
      "president_name": {"type": "text"},
      "general_secretary_name": {"type": "text"}
    }
  }
}
```

---

#### 3. Activities Index

```json
{
  "mappings": {
    "properties": {
      "id": {"type": "keyword"},
      "title": {"type": "text"},
      "description": {"type": "text"},
      "category": {"type": "keyword"},
      "status": {"type": "keyword"},
      "user_id": {"type": "keyword"},
      "user_name": {"type": "text"},
      "jurisdiction_id": {"type": "keyword"},
      "jurisdiction_name": {"type": "text"},
      "jurisdiction_level": {"type": "keyword"},
      "is_public": {"type": "boolean"},
      "created_at": {"type": "date"},
      "approved_at": {"type": "date"}
    }
  }
}
```

---

### Indexing Strategy

#### Real-Time Indexing (Preferred)

**Option 1: Database Triggers + Queue**

```sql
-- PostgreSQL trigger example
CREATE OR REPLACE FUNCTION notify_user_change()
RETURNS TRIGGER AS $$
BEGIN
  PERFORM pg_notify('user_changed', row_to_json(NEW)::text);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_changed_trigger
AFTER INSERT OR UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION notify_user_change();
```

**Application listens to notifications** → Indexes to OpenSearch

**Option 2: Debezium CDC (Change Data Capture)**
- Captures PostgreSQL WAL (Write-Ahead Log)
- Streams changes to Kafka
- Kafka → OpenSearch Sink Connector
- **Pros**: Zero application code, guaranteed consistency
- **Cons**: Additional infrastructure (Kafka)

**Recommended for BJDMS**: Option 1 (simpler, fewer dependencies)

---

#### Bulk Indexing (Historical Data)

```python
# Python script for initial bulk indexing
from opensearchpy import OpenSearch, helpers
import psycopg2

# Connect to PostgreSQL
conn = psycopg2.connect("postgresql://...")
cursor = conn.cursor()

# Connect to OpenSearch
client = OpenSearch([{'host': 'opensearch', 'port': 9200}])

# Fetch all users
cursor.execute("""
  SELECT u.id, u.full_name, u.phone, 
         j.name as jurisdiction_name, j.level_id,
         p.name as position_name, c.committee_type_id
  FROM users u
  LEFT JOIN committee_members cm ON u.id = cm.user_id AND cm.ended_at IS NULL
  LEFT JOIN committees c ON cm.committee_id = c.id
  LEFT JOIN positions p ON cm.position_id = p.id
  LEFT JOIN jurisdictions j ON c.jurisdiction_id = j.id
  WHERE u.is_active = true
""")

# Prepare bulk actions
actions = []
for row in cursor.fetchall():
    action = {
        "_index": "users",
        "_id": row[0],
        "_source": {
            "id": row[0],
            "full_name": row[1],
            "phone": row[2],
            "jurisdiction_name": row[3],
            "jurisdiction_level": row[4],
            "position_name": row[5],
            "committee_type": row[6],
            "is_public": True  # Default
        }
    }
    actions.append(action)

# Bulk index
helpers.bulk(client, actions)
```

**Frequency**: Run once for initial setup, then real-time updates take over

---

### Search API Implementation

#### Global Search Endpoint

**Endpoint**: `GET /api/v1/search`

**Query Parameters**:
```
q: Search query (supports Bangla & English)
type: users | committees | activities | events | all
jurisdiction: Filter by jurisdiction ID
level: central | division | district | upazila | municipality | union | ward
position: Filter by position name
limit: Results per page (default 20, max 100)
offset: Pagination offset
```

**Example Request**:
```
GET /api/v1/search?q=সভাপতি&type=users&level=district&limit=10
```

**OpenSearch Query**:
```json
{
  "query": {
    "bool": {
      "must": [
        {
          "multi_match": {
            "query": "সভাপতি",
            "fields": ["full_name^2", "full_name_bn^2", "position_name_bn"],
            "type": "best_fields",
            "fuzziness": "AUTO"
          }
        }
      ],
      "filter": [
        {"term": {"jurisdiction_level": "district"}},
        {"term": {"is_active": true}}
      ]
    }
  },
  "from": 0,
  "size": 10,
  "highlight": {
    "fields": {
      "full_name": {},
      "full_name_bn": {}
    }
  }
}
```

**Response**:
```json
{
  "total": 64,
  "results": [
    {
      "id": "uuid",
      "full_name": "John Doe",
      "full_name_bn": "জন ডো",
      "position": "President",
      "position_bn": "সভাপতি",
      "jurisdiction": "Dhaka District",
      "jurisdiction_bn": "ঢাকা জেলা",
      "committee_type": "Full Committee",
      "phone": "+880171234567",  // Only if user has permission
      "highlight": {
        "position_name_bn": ["<em>সভাপতি</em>"]
      }
    }
  ]
}
```

---

#### Autocomplete Endpoint

**Endpoint**: `GET /api/v1/search/autocomplete`

**Query Parameters**:
```
q: Partial query (minimum 2 characters)
type: users | committees
limit: Max suggestions (default 5)
```

**Example**:
```
GET /api/v1/search/autocomplete?q=আব&type=users&limit=5
```

**OpenSearch Query**:
```json
{
  "query": {
    "multi_match": {
      "query": "আব",
      "fields": ["full_name.autocomplete", "full_name_bn.autocomplete"],
      "type": "bool_prefix"
    }
  },
  "size": 5,
  "_source": ["full_name", "full_name_bn", "position_name", "jurisdiction_name"]
}
```

**Response**:
```json
{
  "suggestions": [
    {"name": "আবুল কাশেম", "position": "President", "jurisdiction": "Dhaka"},
    {"name": "আব্দুর রহমান", "position": "General Secretary", "jurisdiction": "Chittagong"}
  ]
}
```

---

### Public Directory

#### Public Committee Leader Search

**Endpoint**: `GET /api/v1/public/directory` (No authentication required)

**Features**:
- Search only users with `is_public = true`
- Phone numbers masked (show only if user opted in)
- Committee leaders (President, General Secretary, Organizational Secretary)

**Query**:
```json
{
  "query": {
    "bool": {
      "must": [
        {"term": {"is_public": true}},
        {"terms": {"position_name": ["President", "General Secretary", "Convener", "Member Secretary"]}}
      ]
    }
  }
}
```

**Privacy Controls**:
- Users can set `profile_visibility` in their profile:
  - `public`: Full details visible
  - `members_only`: Only authenticated users
  - `leaders_only`: Only committee leaders
  - `private`: No search visibility

---

## Business Logic Rules

### Search Result Filtering by Permission

**Authenticated Users**:
- Can search all members in their jurisdiction + child jurisdictions
- Can search all committees in their jurisdiction + child jurisdictions
- Can search activities/events in their jurisdiction

**Public (Unauthenticated)**:
- Can search directory (leaders only, with `is_public=true`)
- Cannot search activities, complaints, audit logs

**Super Admin**:
- Can search everything across all jurisdictions

**Implementation**:
```
ON search_request:
  base_query = build_search_query(user_query)
  
  IF user is authenticated:
    user_jurisdictions = get_user_accessible_jurisdictions(user_id)
    base_query.add_filter("jurisdiction_id IN user_jurisdictions")
  ELSE:
    base_query.add_filter("is_public = true")
    base_query.add_filter("position IN leadership_positions")
  
  results = opensearch.search(base_query)
  return sanitize_results(results, user_permissions)
```

---

### Bangla Language Support

**Challenges**:
- Bangla font variations (Unicode normalization)
- Stemming for Bangla words
- Synonym handling (সভাপতি = President)

**Solutions**:
- **ICU Analysis Plugin**: Unicode normalization
- **Bengali Stemmer**: Reduces words to root forms
- **Synonym File**: Bangla ↔ English mappings

**Synonym Examples**:
```
সভাপতি, President
সাধারণ সম্পাদক, General Secretary
জয়পুরহাট, Joypurhat
ঢাকা, Dhaka
```

---

### Search Result Ranking

**Ranking Factors**:
1. **Exact match** (highest priority)
2. **Position importance** (President > Member)
3. **Jurisdiction level** (Central > District > Ward)
4. **Recent activity** (recently active users rank higher)
5. **Verification status** (verified users > unverified)

**Boosting**:
```json
{
  "query": {
    "function_score": {
      "query": { "multi_match": {...} },
      "functions": [
        {"filter": {"term": {"position_name": "President"}}, "weight": 3},
        {"filter": {"term": {"jurisdiction_level": "central"}}, "weight": 2},
        {"filter": {"term": {"is_verified": true}}, "weight": 1.5}
      ],
      "score_mode": "multiply"
    }
  }
}
```

---

## Integration Points

### With PostgreSQL
- Triggers notify OpenSearch indexer on data changes
- Bulk indexing script for historical data
- Schema changes require index mapping updates

### With Phase C (Auth)
- Permission-based result filtering
- Jurisdiction-aware search
- Public vs authenticated search paths

### With DOCKER_INFRASTRUCTURE.md
- OpenSearch container configuration
- Network connectivity to API service
- Volume persistence for indices

---

## Testing Requirements

### Indexing Tests
- User created → Indexed in OpenSearch
- User updated → Index updated
- User deleted (soft delete) → Removed from public search
- Bulk indexing script completes without errors

### Search Tests
- Search Bangla name → Returns correct results
- Search English name → Returns correct results
- Fuzzy search (misspelling) → Returns intended result
- Autocomplete with 2 characters → Returns suggestions
- Permission filtering → Users see only allowed results

### Performance Tests
- Search query < 100ms response time
- Autocomplete < 50ms response time
- 1000 concurrent searches handled
- Index refresh lag < 1 second

---

## Future Extensibility

### Advanced Features

**Geospatial Search**:
- Search committees near a location
- Requires `venue_coordinates` in events
- OpenSearch geo_point mapping

**Aggregations & Analytics**:
- "How many activities by district?"
- Trending search queries
- Popular committee members

**AI-Powered Search**:
- Natural language queries: "Show me all presidents in Dhaka"
- Semantic search (vector embeddings)
- Query intent detection

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
