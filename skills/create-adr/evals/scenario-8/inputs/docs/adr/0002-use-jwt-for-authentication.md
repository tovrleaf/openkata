---
status: ACCEPTED
date: 2024-05-08
authors: [bob@example.com, alice@example.com]
---

# 0002. Use JWT for Authentication

## Context

OrderFlow's API needed an authentication mechanism for both partner integrations and internal service-to-service calls. The initial prototype used HTTP Basic Auth over TLS for simplicity, but this became untenable as the number of callers grew and key rotation became difficult.

## Decision Drivers

- Stateless tokens allow horizontal scaling without shared session storage
- Partners need short-lived tokens they can rotate without contacting support
- Internal services must authenticate with minimal latency overhead

## Decision

We will use signed JWT (RS256) tokens issued by our auth service for all API authentication. Partner tokens expire in 1 hour; internal service tokens expire in 15 minutes.

## Alternatives Considered

### API Keys (static long-lived secrets)

- **Pros:** Simple to implement and explain to partners
- **Cons:** Long-lived; rotation is a manual support operation; no built-in expiry
- **Rejected because:** Key compromise requires immediate manual intervention; rotation at scale is operationally expensive.

### OAuth 2.0 with Opaque Tokens

- **Pros:** Industry standard; broad library support
- **Cons:** Requires token introspection endpoint; adds latency on every request
- **Rejected because:** Token introspection adds ~30ms per request; unacceptable for high-frequency internal calls.

## Consequences

### Positive

- Tokens are self-contained — services can verify signatures locally without a round-trip
- Built-in expiry reduces blast radius of compromised tokens

### Negative

- Token revocation before expiry requires a denylist — additional infrastructure complexity
- RS256 key rotation requires coordinating public key distribution across all services

### Neutral

- Partners must implement token refresh logic; client SDK will abstract this
