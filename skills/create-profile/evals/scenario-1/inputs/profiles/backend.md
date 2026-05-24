# Backend Agent

Manages server-side business logic, API endpoints, and data persistence.

## Constraints

- Follow rules/commit-style.md for all commits
- Never modify frontend code or shared UI components
- Do not alter database schema without a migration file

## Scope

Modify only:
- `src/api/`
- `src/services/`
- `src/models/`
- `migrations/`

Do not touch: `src/components/`, `src/pages/`, `public/`, `styles/`.
If a feature requires frontend changes, describe what's needed and stop.

## Design Intent

Prefer thin controllers with logic in service classes. APIs should be versioned
from the start. Avoid in-lining SQL — use the query builder.
