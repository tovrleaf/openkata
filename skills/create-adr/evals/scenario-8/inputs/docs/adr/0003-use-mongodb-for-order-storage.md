---
status: ACCEPTED
date: 2024-07-20
authors: [charlie@example.com]
---

# 0003. Use MongoDB for Order Storage

## Context

OrderFlow needed a database for storing order records. At launch, order structure varied significantly by partner — different partners had different required fields, custom metadata, and varying nested line-item schemas. A flexible document store seemed like the right fit for this heterogeneous data. The team had MongoDB experience and evaluated it alongside PostgreSQL and DynamoDB.

## Decision Drivers

- Order schemas varied widely by partner at launch; rigidity of SQL was seen as a risk
- Team had existing MongoDB expertise and an established ops runbook for it
- Time pressure favoured an approach the team already knew

## Decision

We will use MongoDB as the primary store for order records, with one collection per environment (orders_prod, orders_staging).

## Alternatives Considered

### PostgreSQL with JSONB

- **Pros:** ACID transactions, strong consistency, JSONB supports flexible schemas
- **Cons:** Team unfamiliar with JSONB query patterns; migration tooling (Flyway) was already set up for a different service
- **Rejected because:** The learning curve combined with launch pressure led to choosing the more familiar option.

### DynamoDB

- **Pros:** Managed, scales automatically, low operational overhead
- **Cons:** Limited query flexibility; secondary indexes have strict design constraints; vendor lock-in
- **Rejected because:** Ad-hoc query requirements for the operations dashboard ruled out DynamoDB's constrained query model.

## Consequences

### Positive

- Flexible schema accommodates partner-specific fields without migrations
- Team can iterate quickly using familiar tooling

### Negative

- No multi-document ACID transactions without explicit session management
- Schema flexibility makes it harder to enforce data invariants at the database level

### Neutral

- Index management requires manual attention as query patterns evolve
