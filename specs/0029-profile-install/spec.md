---
status: Draft
depth: Standard
---

# Profiles That Agents Load

## Story

As someone who installs an OpenKata profile, I want my agent to
load it as an agent role. Today the install only writes a file
that no agent reads.

## Background

A profile, such as `spec-planner`, is one Markdown file. It has a
`description` and `tags` in its frontmatter, and a body with the
role's constraints, companion skills and scope.

The install writes it to `.agents/profiles/<name>/`. Before this
spec, the install also offered `.claude/profiles/<name>/`. No agent
loads profiles from either folder:

- Skills have a shared folder, `.agents/skills/`, that most agents
  read. Profiles have no shared folder.
- Claude Code has no profiles folder. Its closest feature is the
  subagent, a Markdown file in `.claude/agents/<name>.md`. Its
  frontmatter has a `name`, a `description` and optional `tools`
  and `model`, and its body holds the role's instructions. A
  profile has almost the same shape.
- Other agents have their own role features, such as custom agents
  or modes, each with its own file format and folder.

## Requirements

### Investigation

- List the agents that have a role or subagent feature, with its
  folder, its file format and the fields it needs. Start with the
  agents named on the Getting Started page.
- Decide which of those agents an install supports.
- Decide whether a profile stays one generic file that the install
  turns into each agent's format, or whether OpenKata publishes one
  file per agent.
- Decide how the profile's "Scope" section, which limits what the
  role may read and write, maps to each agent's tool permissions.

### Changes (to be fixed by the investigation)

- After an install, Claude Code loads the profile as a subagent
  from `.claude/agents/<name>.md`.
- Each other supported agent loads the profile through its own
  role feature.
- The `install_profile` tool description says where to write the
  profile for each supported agent.
- The install guide and the Getting Started page show the same.

## Out of Scope

- Changing what existing profiles say, beyond the frontmatter an
  agent needs.
- Skills and rules. They already have a shared folder or load
  through `AGENTS.md`.

## Rejected Approaches

- `.claude/profiles/<name>/`. Claude Code does not read it.

## Open Questions

- Which agents to support first, besides Claude Code.
- One generic profile file converted at install time, or one file
  per agent.
- Whether `.agents/profiles/` stays as the shared source, with
  links into each agent's folder.
- How the "Scope" section maps to tool permissions.
