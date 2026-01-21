# DEVELOPMENT_WORKFLOW.md

## Purpose

Defines **development workflow and collaboration processes** for BJDMS development team.

**Goal**: Efficient, collaborative development with minimal conflicts and high code quality.

---

## Git Branching Strategy

### Branch Types

```
main (production-ready)
  ├── develop (integration branch)
  │   ├── feature/add-municipality
  │   ├── feature/anonymous-complaints
  │   └── bugfix/phone-validation
  └── hotfix/security-patch
```

**Branches**:
- `main`: Production code, always deployable
- `develop`: Integration branch (optional, for larger teams)
- `feature/*`: New features
- `bugfix/*`: Bug fixes
- `hotfix/*`: Urgent production fixes

---

### Workflow

#### 1. Create Feature Branch

```bash
# Update main
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/add-municipality-level
```

---

#### 2. Develop & Commit

**Commit Often**:
```bash
git add .
git commit -m "feat: add Municipality to jurisdiction_levels table"
```

**Commit Message Format**: `<type>: <description>`

**Types**: feat, fix, docs, refactor, test, chore

---

#### 3. Push & Create Pull Request

```bash
git push origin feature/add-municipality-level
```

**On GitHub/GitLab**:
- Create Pull Request to `main`
- Fill PR template
- Request reviewers

---

#### 4. Code Review

**Reviewer Checks**:
- [ ] Code follows standards (CODE_STANDARDS.md)
- [ ] Tests included and passing
- [ ] No security vulnerabilities
- [ ] Documentation updated
- [ ] Backward compatible

**Approval**: Minimum 1 approval required

---

#### 5. Merge

**Squash and Merge** (keeps history clean):
```bash
# Automated via GitHub "Squash and merge" button
```

**Delete Feature Branch** after merge.

---

## Daily Development Flow

### Morning Routine

1. **Pull latest changes**:
   ```bash
   git checkout main
   git pull origin main
   ```

2. **Check for conflicts** in open branches:
   ```bash
   git checkout feature/my-feature
   git rebase main
   ```

3. **Review notifications**: PRs, issues, comments

---

### During Development

1. **Work on ONE feature at a time**
2. **Commit frequently** (every logical unit of work)
3. **Write tests** alongside code
4. **Run tests locally** (on remote server):
   ```bash
   ssh root@grayhawks.com -p 25920
   cd /opt/bjdms
   go test ./...
   ```

---

### Before Pushing

**Pre-Push Checklist**:
- [ ] Code compiles/runs
- [ ] All tests pass
- [ ] Linter clean (`golangci-lint run ./...`)
- [ ] No debug logs or commented code
- [ ] Commit messages clear

---

## Code Review Process

### Submitting PR

**PR Template**:
```markdown
## What Changed
Brief description of changes

## Why
Reason for change, issue it addresses

## How to Test
1. Step-by-step testing instructions
2. Expected results

## Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes (or migration plan documented)
- [ ] Security reviewed
```

---

### Reviewing PR

**Response Time**: Within 24 hours

**Review Focus**:
1. **Functionality**: Does it work as intended?
2. **Code Quality**: Readable, maintainable?
3. **Security**: Input validated, no vulnerabilities?
4. **Performance**: No obvious performance issues?
5. **Tests**: Adequate coverage?

**Feedback Style**:
- **Constructive**: "Consider using X instead of Y because..."
- **Specific**: "Line 45: potential SQL injection, use parameterized query"
- **Approving**: Acknowledge good work

---

### Handling Feedback

**Author's Responsibility**:
- Address all comments
- Mark resolved when fixed
- Ask questions if unclear
- Re-request review when ready

---

## Pair Programming

**When to Pair**:
- Complex features (e.g., Convener → Full Committee transition)
- Security-critical code (auth, finance)
- Knowledge transfer (senior + junior)

**Remote Pairing**:
- Screen share via Zoom/Google Meet
- Driver (writes code) + Navigator (reviews, suggests)
- Switch roles every 30 minutes

---

## Issue Tracking

### GitHub Issues

**Labels**:
- `bug`: Something broken
- `feature`: New functionality
- `enhancement`: Improvement to existing feature
- `security`: Security vulnerability
- `documentation`: Documentation update

**Issue Template**:
```markdown
**Description**
Clear description of issue/feature

**Steps to Reproduce** (for bugs)
1. Go to...
2. Click...
3. See error

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Priority**
P0 (Critical) / P1 (High) / P2 (Medium) / P3 (Low)
```

---

## Documentation Updates

**When Code Changes, Update Docs**:
- New feature → Update phase docs (PHASE_X.md)
- Schema change → Update DATABASE_SCHEMA.md
- API change → Update API_CONTRACT.md
- Deployment change → Update DEPLOYMENT_GUIDE.md

**In Same PR**: Documentation changes should be in same PR as code changes.

---

## Communication

### Daily Standup (Optional for small teams)

**15 minutes, daily**:
- What I did yesterday
- What I'm doing today
- Any blockers

**Remote**: Async via Slack/Discord message

---

### Weekly Planning

**Review**:
- Completed work
- Upcoming priorities
- Blockers/risks

---

## Onboarding New Developers

### Day 1

1. **Access Setup**:
   - GitHub/GitLab access
   - SSH key for grayhawks.com
   - Slack/Discord

2. **Read Documentation**:
   - SYSTEM_ARCHITECTURE.md
   - CODE_STANDARDS.md
   - This file (DEVELOPMENT_WORKFLOW.md)

3. **Environment Setup**:
   - Clone repository
   - SSH to grayhawks.com
   - Run test suite

---

### Week 1

**Starter Tasks**:
- Fix a "good first issue" bug
- Review PRs from team
- Read Phase C-H documentation

---

## Conflict Resolution

### Merge Conflicts

**Prevention**:
- Pull `main` frequently
- Keep feature branches short-lived (<1 week)
- Communicate with team about overlapping work

**Resolution**:
```bash
git checkout feature/my-feature
git pull origin main
# Resolve conflicts in files
git add .
git commit -m "fix: resolve merge conflicts"
git push origin feature/my-feature
```

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
