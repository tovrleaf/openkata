---
status: ACCEPTED
date: 2023-08-15
authors: [platform-team]
---

# 0001. Use Redis for Session Storage

## Context

The PayStream API initially stored service sessions in-process memory, which prevented horizontal scaling — deploying multiple API instances caused sessions to be lost when requests hit different pods. A shared session store was needed to support the planned move to a multi-replica Kubernetes deployment.

## Decision Drivers

- Must support concurrent access from multiple API replicas
- Low latency reads for every authenticated request
- Ops team already runs Redis for the rate-limiter cache
- Minimal operational overhead to introduce

## Decision

We will use Redis as the shared session store via the `connect-redis` adapter, replacing in-process memory storage.

## Alternatives Considered

### In-process memory with sticky sessions

- **Pros:** Zero infrastructure changes; simplest possible setup.
- **Cons:** Requires sticky-session routing in the load balancer; sessions lost on pod restart; incompatible with auto-scaling.
- **Rejected because:** Sticky sessions are fragile under rolling deploys and prevent true horizontal scaling.

### PostgreSQL session table

- **Pros:** Durable; already in the stack.
- **Cons:** Higher read latency than Redis; requires a session GC job; adds load to the primary database.
- **Rejected because:** Session lookups on every request would add unnecessary database load and latency.

## Consequences

### Positive

- API can scale horizontally without session affinity requirements.
- Sessions survive pod restarts.

### Negative

- Redis becomes a required dependency; its availability now affects API uptime.
- Session data must be serializable (no circular references or functions).

### Neutral

- Session TTL management moves to Redis TTL configuration rather than application logic.
