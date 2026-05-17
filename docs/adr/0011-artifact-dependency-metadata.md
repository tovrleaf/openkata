---
status: PROPOSED
date: 2026-05-17
authors: [niko.kivela]
---

# 0011. Artifact dependency metadata in frontmatter

## Context

Skills and rules in this project can complement each other.
For example, the `commit-conventions` skill enforces a workflow
that pairs well with the `git-naming` rule (which carries the
format). Neither requires the other, but users benefit from
knowing about the companion.

Currently there is no way to express these relationships in
metadata. Users discover companions only by reading docs. The
MCP server and web UI cannot suggest related artifacts during
installation.

## Decision Drivers

- Artifacts must remain standalone — no hard coupling
- Relationships should be machine-readable for tooling
  (MCP install suggestions, web UI "related" sections)
- Must not break the Agent Skills spec (unknown fields are
  ignored by other platforms)
- Should distinguish "works better with" from "broken without"

## Questions to Resolve

- What field names? `recommends` / `requires`? `companions`?
  `enhances`?
- Should the relationship be typed? (e.g., "rule", "skill",
  "profile")
- Bidirectional or one-way? If `commit-conventions` recommends
  `git-naming`, should `git-naming` also recommend back?
- How to handle cross-repo references? Just names, or
  namespaced identifiers?
- Should `requires` trigger auto-install, or just block with
  an error?
- How does this interact with versioning? Should recommendations
  pin a version range?

## Options Under Consideration

### Option A: recommends / requires

```yaml
metadata:
  tags: "category:conventions"
  recommends:
    - git-naming
  requires:
    - markdown-style
```

- Simple, familiar from package managers
- Clear semantics: recommends = optional, requires = mandatory
- No type information — relies on naming conventions

### Option B: companions with type

```yaml
metadata:
  companions:
    - name: git-naming
      type: rule
      relationship: enhances
```

- Richer, supports tooling better
- More verbose, harder to write by hand
- Relationship types need a controlled vocabulary

### Option C: freeform references section

```yaml
metadata:
  related:
    - "rule:git-naming — carries the format this skill enforces"
```

- Human-readable, flexible
- Harder to parse programmatically
- No clear install semantics

## Decision

*Pending discussion.*

## Consequences

*To be determined after decision.*
