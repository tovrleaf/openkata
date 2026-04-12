---
status: PROPOSED
date: 2026-04-12
authors: [niko.kivela]
---

# 0003. Use Go as the implementation language

## Context

OpenKata needs an implementation language for its server-side components:
the MCP server for skill distribution (ADR 0002) and potentially a web
interface if the project gains traction. The language choice affects
development speed, runtime performance, deployment cost, and long-term
maintainability.

The project is currently maintained by a small team that relies on AI coding
agents to write and maintain the code.

## Decision Drivers

- Must produce fast, low-memory binaries for cheap AWS deployment
- Must be testable with a strong standard testing story
- Must compile to a single binary with no runtime dependencies
- Should support both MCP server (stdio) and HTTP server (web/API) use cases
  in one language
- Agent-friendly — AI coding agents must be able to write and maintain the code
  effectively

## Decision

We will use Go for all server-side components: the MCP server, API backend,
and web serving. Go's `net/http` standard library handles HTTP natively, and
a single binary can serve both the MCP protocol (stdio) and a web interface
(Lambda behind API Gateway).

The MCP server uses `github.com/mark3labs/mcp-go` (v0.47.1, 8.6K stars,
MIT licensed).

## Alternatives Considered

### Python

- **Pros:** Official Anthropic MCP SDK, fastest to prototype, largest
  AI/ML ecosystem
- **Cons:** Slower startup, higher memory usage, runtime dependency
  (Python interpreter), dynamic typing works against "sound, never break"
  goal
- **Rejected because:** Higher infrastructure costs at scale and weaker
  type safety. Good for prototyping but not the right fit for a server
  that should be fast and reliable.

### Rust

- **Pros:** Fastest execution, lowest memory, strongest type safety,
  smallest binaries
- **Cons:** Slower development iteration, steeper learning curve for
  contributors and agents, overkill for I/O-bound file-serving workloads
- **Rejected because:** The server is I/O-bound (reading and copying files),
  not CPU-bound. Rust's performance advantage is wasted here, and the
  development speed tradeoff isn't worth it.

## Consequences

### Positive

- Single binary deployment — no runtime dependencies
- Fast startup and minimal memory footprint, keeping AWS costs low
- One language covers MCP server, API, and web serving
- Strong standard library for HTTP, testing, and file I/O
- `go test` is built in with no external test framework needed
- Well-supported by AI coding agents

### Negative

- Go's error handling is verbose compared to Python or Rust
- Less expressive type system than Rust (no sum types, no exhaustive matching)

### Neutral

- Go's simplicity means agents produce consistent, readable code
- The mcp-go library is actively maintained but not an official Anthropic SDK

## References

- [mcp-go](https://github.com/mark3labs/mcp-go) — Go MCP SDK used by this project
- ADR 0002 — Use MCP server for distribution (this ADR covers the language for that server)
