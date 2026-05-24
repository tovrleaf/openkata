---
status: ACCEPTED
date: 2024-01-15
authors:
  - sarah.chen
  - marcos.oliveira
---

# 0001. Use REST for API Communication

## Context

The platform was starting from scratch in early 2024. The team needed a standard way for the frontend clients (web and mobile) and internal services to communicate. We evaluated several API paradigms and needed to pick one to standardize on.

## Decision Drivers

- Team was already experienced with REST and HTTP conventions
- Excellent tooling support across all languages in our stack
- Simple to reason about and debug with standard HTTP clients
- No additional runtime dependencies required

## Decision

We will use REST over HTTP/JSON for all API communication — both client-to-service and service-to-service calls.

## Alternatives Considered

### GraphQL

- **Pros:** Flexible querying, reduces over-fetching, single endpoint
- **Cons:** Added complexity, learning curve, requires specialized tooling
- **Rejected because:** Team had no GraphQL experience and the product was simple enough at the time that REST would suffice.

### gRPC

- **Pros:** High performance, strongly typed contracts, streaming support
- **Cons:** Poor browser support, requires Protobuf schema management
- **Rejected because:** Browser client support was a hard requirement; gRPC-Web adds complexity we didn't need.

## Consequences

### Positive

- Fast to implement; all engineers knew REST
- Easy to test with curl, Postman, and standard HTTP tooling

### Negative

- Over-fetching and under-fetching may become issues as the product grows
- Multiple round-trips required for complex views

### Neutral

- OpenAPI spec required for contract documentation

## References

- [REST API Design Best Practices](https://restfulapi.net/)
