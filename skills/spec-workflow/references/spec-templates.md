# Spec Templates

## brief.md (Quick depth only)

```markdown
# Brief: Title

Scope: one or two sentences describing what this change does.

Date: YYYY-MM-DD
```

## spec.md

```markdown
---
status: Draft
depth: Standard
---

# Feature Title

## Story
As a ..., I want ..., so that ...

## Requirements
- Acceptance criteria

## Out of Scope
- What this feature does not include

## Rejected Approaches
- Approach and why it was rejected

## Open Questions
- Unresolved items
```

## design.md

```markdown
# Design: Feature Title

## Architecture
- Approach, file paths, components affected

## Decisions
- Key trade-offs and rationale

## Dependencies
- External systems, libraries, other features
```

## tasks.md

```markdown
# Tasks: Feature Title

## Tasks

### 1. Task name
- **Status**: Pending
- **Goal**: What this achieves
- **Boundary**: Files/modules this task touches
- **Key files**: Existing code to reference for patterns
- **Depends**: Task numbers that must complete first
- **Done when**: Observable criteria
- **Verify**: Command to prove completion (optional)

### 2. Next task
- **Status**: Pending
- **Goal**: ...
- **Boundary**: ...
- **Key files**: ...
- **Depends**: ...
- **Done when**: ...
- **Verify**: ...

## Progress Log
```

## validation-report.md

```markdown
# Validation: Feature Title

## Requirements Check
- [x] Requirement — verified by ...
- [ ] Requirement — FAILED: ...

## Out of Scope Verified
- Confirmed: no work done on excluded items

## Issues Found
- Description of any problems found
```
