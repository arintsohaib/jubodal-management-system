# BACKUP_RECOVERY.md

## Purpose

Defines **backup strategy and disaster recovery procedures** for BJDMS to ensure data protection and business continuity.

**RPO (Recovery Point Objective)**: 24 hours (max data loss acceptable)  
**RTO (Recovery Time Objective)**: 4 hours (max downtime acceptable)

---

## Backup Strategy

### What to Backup

1. **PostgreSQL Database** (Critical)
   - All tables
   - Full dump including schema + data
   - Frequency: Daily

2. **File Storage** (S3/MinIO) (Important)
   - Activity proofs
   - Complaint evidence
   - Donation receipts
   - Profile photos
   - Frequency: Weekly, incremental daily

3. **Configuration Files** (Important)
   - `/opt/bjdms/.env`
   - `/opt/bjdms/docker-compose.yml`
   - NGINX configs
   - Frequency: On change (manual)

4. **Redis Data** (Optional)
   - Sessions (can be regenerated)
   - Frequency: None (ephemeral data)

---

## Database Backup

### Daily Automated Backup

**Cron Job** (on server):
```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /opt/bjdms/scripts/backup_database.sh
```

**Backup Script** (`/opt/bjdms/scripts/backup_database.sh`):
```bash
#!/bin/bash
set -e

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/bjdms/backups/database"
BACKUP_FILE="bjdms_db_${DATE}.sql.gz"

# Create backup directory
mkdir -p $BACKUP_DIR

# Dump database
docker-compose exec -T postgres pg_dump -U bjdms_user bjdms_db | gzip > $BACKUP_DIR/$BACKUP_FILE

# Verify backup
if [ -f "$BACKUP_DIR/$BACKUP_FILE" ]; then
    echo "Backup successful: $BACKUP_FILE"
    
    # Delete backups older than 30 days
    find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete
    
    # Optional: Upload to offsite storage (S3, Google Drive)
    # aws s3 cp $BACKUP_DIR/$BACKUP_FILE s3://bjdms-backups/database/
else
    echo "Backup failed!" | mail -s "BJDMS Backup Failed" admin@bjd.org.bd
    exit 1
fi
```

**Make executable**:
```bash
chmod +x /opt/bjdms/scripts/backup_database.sh
```

---

### Weekly Full Backup

**Cron Job**:
```bash
# Every Sunday at 3 AM
0 3 * * 0 /opt/bjdms/scripts/backup_full.sh
```

**Script**: Same as daily, but upload to offsite storage.

---

## File Storage Backup

### Using rsync (Simple)

**Weekly Backup** (Sundays 4 AM):
```bash
0 4 * * 0 /opt/bjdms/scripts/backup_files.sh
```

**Script**:
```bash
#!/bin/bash
DATE=$(date +%Y%m%d)
BACKUP_DIR="/opt/bjdms/backups/files"
SOURCE_DIR="/opt/bjdms/volumes/minio"

mkdir -p $BACKUP_DIR

# Incremental backup
rsync -av --delete $SOURCE_DIR $BACKUP_DIR/minio_$DATE

# Keep last 4 weekly backups only
find $BACKUP_DIR -maxdepth 1 -type d -name "minio_*" -mtime +28 -exec rm -rf {} \;
```

---

## Offsite Backup

**Options**:

1. **S3-Compatible Storage** (AWS, Wasabi, Backblaze B2)
   - Automated sync from backup directory
   - Encrypted upload
   - Lifecycle rules (delete after 1 year)

2. **External Drive** (Manual)
   - Monthly copy to external drive
   - Stored offsite (different location)

3. **Cloud Storage** (Google Drive, Dropbox)
   - For config files and weekly database dumps
   - Encrypted before upload

**Recommended**: S3-compatible for automation + external drive for offline safety.

---

## Backup Verification

### Monthly Backup Test

**Process**:
1. Download latest backup
2. Restore to test environment (`/opt/bjdms-test`)
3. Verify data integrity (spot-check critical tables)
4. Document test results

**Checklist**:
- [ ] Database restores without errors
- [ ] User count matches expected
- [ ] Financial records intact
- [ ] Files accessible

---

## Disaster Recovery Scenarios

### Scenario 1: Database Corruption

**Symptoms**: Database won't start, corrupted tables

**Recovery Steps**:
1. Stop application:
   ```bash
   docker-compose down
   ```

2. Restore from backup:
   ```bash
   # Find latest backup
   ls -lh /opt/bjdms/backups/database | tail -5
   
   # Restore
   gunzip < /opt/bjdms/backups/database/bjdms_db_20260121_020000.sql.gz | \
   docker-compose exec -T postgres psql -U bjdms_user bjdms_db
   ```

3. Verify restoration:
   ```bash
   docker-compose exec postgres psql -U bjdms_user -d bjdms_db -c "SELECT COUNT(*) FROM users;"
   ```

4. Restart application:
   ```bash
   docker-compose up -d
   ```

**Recovery Time**: ~30 minutes

---

### Scenario 2: Complete Server Failure

**Symptoms**: Server hardware failure, cannot SSH

**Recovery Steps**:

1. **Provision New Server**:
   - Same specs (Debian 13, Docker)
   - Same domain (arint.win) pointing to new IP

2. **Install Dependencies**:
   ```bash
   apt update && apt install docker.io docker-compose git -y
   ```

3. **Restore Code**:
   ```bash
   mkdir /opt/bjdms
   cd /opt/bjdms
   git clone {repository_url} .
   ```

4. **Restore Configuration**:
   - Copy `.env` from offsite backup
   - Copy `docker-compose.yml`

5. **Restore Database**:
   - Download latest backup from offsite
   - Restore per Scenario 1

6. **Restore Files**:
   - Download S3/file backups
   - Place in `/opt/bjdms/volumes/minio`

7. **Start Services**:
   ```bash
   docker-compose up -d
   ```

8. **Verify**:
   - Health check passes
   - Test login
   - Test critical functions

**Recovery Time**: ~4 hours (meets RTO)

---

### Scenario 3: Accidental Data Deletion

**Symptoms**: User or admin accidentally deletes records

**Recovery Steps**:

1. **Identify Deleted Data**:
   - Check audit logs: When was it deleted?
   - What was deleted (table, IDs)?

2. **Find Point-in-Time Backup**:
   - Locate backup from BEFORE deletion

3. **Selective Restore**:
   ```bash
   # Restore entire database to temp location
    gunzip < backup.sql.gz | psql -h localhost -U bjdms_user temp_bjdms_db
   
   # Export specific deleted records
   psql -h localhost -U bjdms_user temp_bjdms_db \
     -c "COPY (SELECT * FROM users WHERE id IN ('uuid1', 'uuid2')) TO STDOUT" | \
     psql -h postgres -U bjdms_user bjdms_db \
     -c "COPY users FROM STDIN"
   ```

4. **Verify**:
   - Restored records visible
   - No duplicates created

**Recovery Time**: ~1-2 hours

---

### Scenario 4: Ransomware Attack

**Symptoms**: Files encrypted, ransom demand

**Immediate Actions**:
1. **Isolate Server**: Disconnect from network
2. **Do NOT pay ransom**
3. **Notify authorities**

**Recovery**:
1. Wipe server
2. Rebuild from scratch (Scenario 2)
3. Restore from OFFSITE backup (not on same server!)
4. Investigate infection vector
5. Patch vulnerability

**Prevention**:
- Offsite backups mandatory (ransomware can't encrypt offsite)
- Regular security updates
- Firewall rules

---

## Backup Retention Policy

| Backup Type | Frequency | Retention |
|-------------|-----------|-----------|
| Daily database | Daily | 30 days |
| Weekly database | Weekly | 12 weeks |
| Monthly database | Monthly | 12 months |
| Files (incremental) | Daily | 7 days |
| Files (full) | Weekly | 4 weeks |
| Configuration | On change | Indefinite (version controlled) |

---

## Backup Monitoring

**Alert if**:
- Backup script fails (exit code ≠ 0)
- Backup file size < 50% of previous backup (possible corruption)
- Backup age > 36 hours (missed backup)

**Notification**: Email + SMS to ops team

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
