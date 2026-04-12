---
status: ACCEPTED
date: 2026-04-12
authors: [niko.kivela]
---

# 0001. Use Agent Skills format for distributing reusable AI agent capabilities

## Context

The openkata project needs a standard format for packaging and distributing
reusable capabilities for AI coding agents. Several options exist: custom
markdown conventions, MCP server tools, or the Agent Skills format
(agentskills.io) — an open specification originally developed by Anthropic
and adopted by VS Code/Copilot, Claude Code, Gemini CLI, and others.

We need a format that is simple to author, portable across agents, and
doesn't require infrastructure to distribute.

## Decision Drivers

- Must work across multiple AI agents without vendor lock-in
- No build step or infrastructure required for distribution
- Low authoring friction — contributors shouldn't need to learn a complex format
- Compatible with existing tooling conventions (adr-tools numbering, markdown)

## Decision

We will use the [Agent Skills format](https://agentskills.io/specification) as defined at agentskills.io/specification.
Each skill is a folder containing a SKILL.md file with YAML frontmatter and
markdown instructions, optionally accompanied by scripts/, references/, and
assets/ directories.

Numbering and format conventions are compatible with
[adr-tools](https://github.com/npryce/adr-tools).

## Alternatives Considered

### Custom markdown conventions

- **Pros:** Full control over format, no external dependency
- **Cons:** No ecosystem support, agents wouldn't discover skills automatically,
  every project would reinvent the structure
- **Rejected because:** No portability — skills would only work with custom tooling

### MCP server tools

- **Pros:** Rich runtime capabilities, structured tool interfaces
- **Cons:** Requires a running server, more complex to author and distribute,
  heavier infrastructure dependency
- **Rejected because:** Too much overhead for packaging static knowledge and
  workflows that are fundamentally markdown-based

## Consequences

### Positive

- Skills are portable across any agent that supports the Agent Skills spec
- No build step or infrastructure required — skills are plain files
- Distribution is simple: users copy a folder into their project

### Negative

- We're coupled to the Agent Skills spec — if it changes significantly, skills
  may need updating
- Contributors must learn the SKILL.md format (though it's minimal: two required
  frontmatter fields + markdown body)

### Neutral

- The spec is still evolving, so we may need to adapt as it matures

## References

- [Agent Skills specification](https://agentskills.io/specification)
- [adr-tools](https://github.com/npryce/adr-tools) — numbering convention
