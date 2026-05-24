# ADR-014: Migrate from Session-Based Auth to JWT

## Status

Proposed

## Context

Our current session-based authentication is creating operational pain as we expand to multiple regions. Session state is stored in a centralised store, which has become a bottleneck and a single point of failure. The team wants to move to a stateless authentication model.

## Decision

We will migrate from session-based authentication to JSON Web Tokens (JWT). Access tokens will be short-lived. Refresh tokens will be used to obtain new access tokens without requiring re-login.

Token revocation will be handled via a blocklist maintained in a cache layer for the duration of the access token lifetime.

## Consequences

- Services no longer depend on a centralised session store
- Token revocation requires maintaining the blocklist
- Clients will need to handle token refresh logic

## Alternatives Considered

- Sticky sessions: rejected because they break horizontal scaling
- Mutual TLS: considered too complex for our use case

## Open Questions

- What should the access token lifetime be?
- Should refresh tokens be rotated on each use?
