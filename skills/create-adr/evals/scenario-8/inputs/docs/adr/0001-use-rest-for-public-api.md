---
status: ACCEPTED
date: 2024-03-12
authors: [alice@example.com]
---

# 0001. Use REST for the Public API

## Context

OrderFlow launched its public API using REST because it was the most widely understood API style among potential integration partners in 2024. The API surface was small — five endpoints — and the team had strong REST experience. GraphQL and gRPC were evaluated but considered premature given the early product stage.

## Decision Drivers

- REST is universally understood by external partner developers
- The team had no GraphQL or gRPC production experience at the time
- The API had fewer than 10 endpoints and straightforward query patterns
- Time-to-market pressure favoured the lowest learning-curve approach

## Decision

We will expose the public API as a REST/JSON API hosted at `/api/v1/`.

## Alternatives Considered

### GraphQL

- **Pros:** Flexible queries, single endpoint, strong typing
- **Cons:** Higher onboarding cost for partners, requires schema introspection tooling
- **Rejected because:** Partner integrators were REST-only shops; the flexibility of GraphQL added complexity the product didn't yet need.

### gRPC

- **Pros:** High performance, strongly typed contracts via Protobuf
- **Cons:** Binary protocol — harder for partners to debug; limited browser support
- **Rejected because:** Partners use browser-based admin tools; gRPC would require a REST gateway proxy anyway.

## Consequences

### Positive

- Every partner engineering team already knows REST; zero onboarding friction on the API style itself
- Curl-friendly — support and QA can test endpoints without special tooling

### Negative

- No built-in schema introspection or code-generation support
- Over-fetching and under-fetching requires versioned endpoint proliferation as the product grows

### Neutral

- API versioning strategy (URL vs header vs content negotiation) left to a future decision
