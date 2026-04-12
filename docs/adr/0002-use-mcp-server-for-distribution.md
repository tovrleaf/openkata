---
status: PROPOSED
date: 2026-04-12
authors: [niko.kivela]
---

# 0002. Use MCP server for distributing skills and agent artifacts

## Context

OpenKata needs a distribution mechanism for reusable AI agent artifacts.
Today, using a skill means cloning the repository and copying folders by
hand. Updates to skills don't flow to existing projects — every project
gets a snapshot that immediately goes stale.

The problem extends beyond skills. Across projects, the same agent
configurations, role prompts, and templates (like the ADR template) are
copy-pasted between repositories. When any of these artifacts improve,
existing projects don't benefit.

The distribution mechanism must also support a private/enterprise use case:
distributing internal skills and agent configurations across teams with
varying tech stacks, directory structures, monorepo layouts, and maturity
levels — all behind an internal developer portal.

Multiple coding agents exist (Kiro, Claude Code, Codex, Cursor, and more
emerging), each with different conventions for where skills and
configurations live. The distribution mechanism needs to adapt to the
target environment rather than assuming a single structure.

## Decision Drivers

- Updates to artifacts must propagate to existing projects, not just
  new installs
- Must distribute more than just skills: agent configs, role prompts,
  templates
- No external runtime dependencies — every coding agent already supports
  MCP, so no additional tooling required
- Must support both public (OpenKata) and private (enterprise) distribution
- Must adapt to different project structures, languages, monorepos, and
  agent conventions
- Should align with the direction the ecosystem is heading

## Decision

We will build an MCP server as the primary distribution mechanism for
OpenKata artifacts. The server will act as an environment-aware installer
that copies files into the user's project, adapting placement to the
target agent and project structure.

Users end up with hard copies of files in their project. The MCP server
is needed for installation and updates, not for runtime operation — projects
work independently after artifacts are installed.

The server will start simple (serve files, detect agent type) and layer
on intelligence incrementally (project structure detection, monorepo
support, update propagation).

## Alternatives Considered

### skills.sh (npx skills)

- **Pros:** Already exists with 90K+ installs, supports 20+ agents,
  handles versioning, zero development effort, established ecosystem
- **Cons:** Requires Node.js, only distributes Agent Skills (not configs,
  prompts, or templates), controlled by Vercel with no influence over
  roadmap, fixed installation pattern with no environment adaptation,
  no private/enterprise distribution support
- **Rejected because:** Doesn't cover the full artifact scope (configs,
  prompts, templates), can't adapt to different project structures, and
  doesn't support private distribution behind an enterprise developer
  portal. Still usable by users who prefer it — this decision doesn't
  block skills.sh usage.

### Git submodules

- **Pros:** Built into Git, no external tooling, supports versioning
- **Cons:** Adds significant Git complexity, poor developer experience,
  difficult to adapt to different project structures, submodule updates
  are error-prone
- **Rejected because:** Too much Git complexity for the benefit. Poor
  fit for teams with varying maturity levels.

### Manual clone and copy

- **Pros:** Zero infrastructure, works everywhere, no dependencies
- **Cons:** No update propagation, no versioning, no adaptation to
  project structure, doesn't scale beyond a handful of projects
- **Rejected because:** This is the current approach and it's the
  problem we're solving.

## Consequences

### Positive

- Single mechanism distributes all artifact types (skills, configs,
  prompts, templates)
- No Node.js or other external dependency — agents already speak MCP
- Environment-aware installation adapts to different agents, project
  structures, and languages
- Supports both public and private distribution with MCP's built-in
  auth (OAuth 2.1)
- Aligned with ecosystem direction — Anthropic has indicated skills
  will be distributed via MCP
- Hard copies mean projects work offline and independently

### Negative

- Requires building and maintaining a server — this is infrastructure,
  not just content
- Must track directory conventions for each supported agent as new
  agents emerge
- More complex than skills.sh for the simple public case
- V1 will have less ecosystem reach than skills.sh's established
  leaderboard and discovery

### Neutral

- skills.sh remains usable alongside this for users who prefer it
- The MCP server scope will grow as new artifact types are identified
- Enterprise/private distribution details will be a separate ADR

## Non-goals

- Defining the exact adaptation logic per agent or project structure —
  that's implementation detail for later
- Building the enterprise developer portal integration — separate ADR
- Replacing skills.sh — users who prefer it can continue using it
- Runtime skill serving (live/hot updates without file copies) — users
  get hard copies

## References

- [skills.sh](https://skills.sh) — existing skills ecosystem by Vercel Labs
- [MCP specification](https://modelcontextprotocol.io) — Model Context Protocol
- [Agent Skills specification](https://agentskills.io/specification)
- ADR 0001 — Use Agent Skills format (covers packaging, this ADR covers distribution)
