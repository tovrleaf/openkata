---
status: ACCEPTED
date: 2024-03-10
authors:
  - priya.nair
---

# 0001. Use Vite for Build Tooling

## Context

The dashboard was originally bootstrapped with Create React App (CRA). As the codebase grew, build times became a developer productivity bottleneck — full rebuilds were taking over 2 minutes. The team needed a faster alternative.

## Decision Drivers

- Developer experience: fast hot module replacement (HMR)
- Build speed: reduce CI pipeline time
- Ecosystem compatibility with React 18

## Decision

We will use Vite as the build tool and dev server for the analytics dashboard, replacing Create React App.

## Alternatives Considered

### Webpack 5 (custom config)

- **Pros:** Mature ecosystem, highly configurable
- **Cons:** Complex configuration, slower than Vite for dev HMR
- **Rejected because:** Configuration overhead and slower HMR outweighed the benefits

### Turbopack

- **Pros:** Very fast (Rust-based), Next.js integration
- **Cons:** Experimental at the time, not stable for production builds
- **Rejected because:** Not production-ready when the decision was made

## Consequences

### Positive

- HMR reduced from ~5s to under 300ms for most changes
- Simplified build config

### Negative

- Some CRA-specific tooling had to be replaced

### Neutral

- Team needed to learn Vite config conventions

## References

- [Vite documentation](https://vitejs.dev/)
