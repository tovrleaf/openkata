---
status: ACCEPTED
date: 2024-09-10
authors: [engineering@shipfast.io]
---

# 0001. Use PostgreSQL for Primary Database

## Context

ShipFast needs a reliable relational database for storing order, customer, and inventory data. The team evaluated several options given the need for ACID transactions and complex relational queries across entities.

## Decision Drivers

- Must support multi-table transactions (orders reference customers, inventory, and shipments)
- Team has strong SQL expertise
- Need for row-level locking during concurrent order processing
- Operational familiarity: team already runs PostgreSQL in staging

## Decision

We will use PostgreSQL as the primary relational database for all transactional data.

## Alternatives Considered

### MySQL

- **Pros:** Widely used, good tooling, slightly lower memory footprint
- **Cons:** Historically weaker support for advanced features (full-text search, JSON operators)
- **Rejected because:** PostgreSQL's richer feature set better matches our query patterns

### SQLite

- **Pros:** Zero configuration, embedded
- **Cons:** Not suitable for concurrent writes in a multi-instance deployment
- **Rejected because:** Cannot handle our multi-server production deployment

## Consequences

### Positive

- Strong transactional guarantees for order processing
- Rich query capabilities including JSON columns for flexible metadata

### Negative

- Requires operational expertise to tune for high-throughput writes

### Neutral

- Migrations managed via standard SQL migration tooling

## References

- PostgreSQL documentation: https://www.postgresql.org/docs/
