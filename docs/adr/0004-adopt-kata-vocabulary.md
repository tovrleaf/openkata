---
status: ACCEPTED
date: 2026-04-12
authors: [niko.kivela]
---

# 0004. Adopt kata vocabulary for project naming and communication

## Context

The project name "Open Kata" draws from the martial arts concept of kata —
codified sequences of movements practiced until mastered. The framework
distributes multiple types of artifacts (skills, roles, prompts, templates,
ADRs) and needs consistent terminology across documentation, tooling, and
communication.

Without a shared vocabulary, documentation mixes generic terms ("server",
"skills", "framework") that don't convey the project's identity or help
users build a mental model of how the pieces fit together.

## Decision Drivers

- Terminology should reinforce the core idea: codified, repeatable practices
  mastered through use
- Must be intuitive for newcomers while giving the project a distinct identity
- Should map cleanly to existing technical concepts

## Decision

We will use the following vocabulary across all documentation, tooling, and
communication:

| Concept | Kata term | Description |
|---------|-----------|-------------|
| Skills | kata | Codified practices agents follow |
| Roles | sensei profiles | Agent role definitions with scoped permissions |
| Prompts | kata forms | Standardized templates for commits, PRs, reviews |
| ADRs | dojo records | Architecture decisions preserved for the school |
| MCP server | the dojo | Where kata are served and practiced |
| Framework | the ryu (school) | The complete system of practices |

## Alternatives Considered

### Use generic technical terms only

- **Pros:** No learning curve, immediately understood
- **Cons:** No project identity, indistinguishable from any other tool
- **Rejected because:** The kata metaphor is the project's identity and
  makes the concepts more memorable and cohesive

### Use a different metaphor (e.g., recipes, playbooks)

- **Pros:** Also intuitive, widely understood
- **Cons:** Overused in the industry, doesn't carry the "mastery through
  practice" connotation
- **Rejected because:** Kata specifically implies codified, repeatable
  practice — which is exactly what the project delivers

## Consequences

### Positive

- Consistent language across docs, tooling, and communication
- Strong project identity that reinforces the core concept
- "The dojo serves kata" is an intuitive mental model for distribution

### Negative

- New contributors need to learn the vocabulary mapping
- Non-English speakers may not know the Japanese terms

### Neutral

- Technical terms (MCP server, skills, ADRs) remain valid and understood —
  the kata vocabulary is an overlay, not a replacement
