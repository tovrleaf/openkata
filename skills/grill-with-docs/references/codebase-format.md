# CODEBASE.md Format

Location: `docs/context/CODEBASE.md`

## Purpose

Maps the project's contexts, where they live, and how
they relate. The "you are here" file for the repo.

## Structure

```md
# Codebase

## Contexts

- [Ordering](./src/ordering/) — receives and tracks
  customer orders
- [Billing](./src/billing/) — generates invoices and
  processes payments
- [Fulfillment](./src/fulfillment/) — manages warehouse
  picking and shipping

## Relationships

- **Ordering → Fulfillment**: Ordering emits `OrderPlaced`
  events; Fulfillment consumes them to start picking
- **Fulfillment → Billing**: Fulfillment emits
  `ShipmentDispatched`; Billing generates invoices
- **Ordering ↔ Billing**: Shared types for `CustomerId`
  and `Money`
```

## Rules

- List each context with its directory and a one-line
  purpose
- Relationships show direction and mechanism (events,
  shared types, API calls)
- Keep it structural — no implementation details
- Update when new contexts emerge or relationships change

## When to create

Create lazily when a grilling session reveals multiple
bounded contexts in the codebase. Single-context projects
don't need this file.
