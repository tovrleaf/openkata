# Git-Based Changelog Skill

## Problem Description

A developer relations team is building an internal assistant that helps engineers prepare release notes and changelogs. Currently, engineers spend 20–30 minutes per release manually scanning commits, filtering noise, and summarizing user-facing changes. The team wants to package this workflow as a reusable agent skill so any engineer can trigger it consistently.

The team has agreed on the following requirements for the skill: it should examine the git history between two references (e.g., tags, branches, or commit SHAs), filter out chore/internal commits, and produce a structured changelog grouped by category (features, fixes, breaking changes). The skill should also be applicable in situations where an engineer just says "what changed since the last release?" or "summarize the recent commits for me."

## Output Specification

Create a skill for this changelog workflow. Place all output under a `changelog-skill/` directory:

- `changelog-skill/SKILL.md` — the complete skill definition
- Any supporting files the skill package needs (place them in appropriately named subdirectories)

Produce a brief `design-notes.md` in the working directory (not inside the skill folder) that explains: what the skill's name is, what trigger phrases are in the description, and the rationale for the boundary decisions you made.
