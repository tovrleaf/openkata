# Example: Standard Depth Spec

A minimal example showing the expected output format for
a Standard depth feature.

## specs/0001-user-search/spec.md

```yaml
---
status: Draft
depth: Standard
---
```

```markdown
# User Search

## Story

As a user, I want to search for other users by name
so I can find and connect with people I know.

## Requirements

- Search input with debounced API calls (300ms)
- Results show name, avatar, and role
- Empty state: "No users found"
- Min 2 characters to trigger search
- Results limited to 20 entries

## Out of Scope

- Pagination of results
- Search filters (by role, department)
- Fuzzy matching

## Open Questions

- None
```

## specs/0001-user-search/tasks.md

```markdown
# Tasks: User Search

## Tasks

### 1. Add search API endpoint
- **Status**: Pending
- **Goal**: Create GET /api/users/search?q= endpoint
- **Boundary**: `src/routes/users.ts`
- **Depends**: None
- **Done when**: Endpoint returns matching users as JSON,
  limited to 20 results, min 2 char query enforced
- **Verify**: `npm test -- --grep "user search"`

### 2. Add search UI component
- **Status**: Pending
- **Goal**: Search input with debounced calls and
  result list display
- **Boundary**: `src/components/UserSearch.tsx`
- **Depends**: 1
- **Done when**: Input triggers search after 300ms,
  results render with name/avatar/role, empty state shows
- **Verify**: `npm test -- --grep "UserSearch"`

## Progress Log
```
