---
name: create-profile
description: >
  Creates agent profiles that define roles with scoped
  permissions and constraints. Investigates the repo for existing
  profiles and conventions, writes a focused profile markdown
  file, and validates scope boundaries. Use when the user wants
  to create an agent profile, define a specialized role, scope an
  agent to specific files or domains, set up a new team member
  persona, restrict what an agent can touch, or wants consistent
  behavior from a domain-specific agent.
metadata:
  version: "1.1.0"
  tags: "category:scaffolding"
---

# Create Profile

Create agent profiles — role definitions that scope an agent to
a specific domain with explicit constraints and permissions.
Profiles are standalone markdown files that define agent roles.

## Workflow

1. **Gather intent** — Understand what the profile should scope.
   Collect:
   - A short name (lowercase, one word preferred)
   - What domain or responsibility the agent covers
   - What it should and should not touch

2. **Investigate the repo** — Before asking questions, search
   for facts that reduce ambiguity:
   - Existing profiles (find where they live in this repo;
     ask the user if ambiguous)
   - Project structure to understand domain boundaries
   - Existing rules that the profile should reference
   - File patterns that define the scope

   Don't ask the user what you can look up yourself.

3. **Clarify** — Ask only questions that materially affect the
   profile:
   - Which files/directories are in scope
   - Which are explicitly out of scope
   - What rules the profile must follow
   - Whether the profile has design authority or just executes

4. **Write the profile** — Structure:

   ```markdown
   # <Name> Agent

   One-line role description — domain and boundary.

   ## Constraints

   - Hard rules this agent must follow
   - Which rules or conventions the profile must follow

   ## Scope

   What files and actions are in bounds. What is explicitly
   out of bounds. When to stop and hand off.

   ## Design Intent

   How this agent should approach decisions within its
   domain. What aesthetic or architectural principles guide
   its choices.
   ```

   Writing principles:
   - **Scope is a boundary, not a task list.** State what's
     in and out, not every possible action.
   - **Constraints are non-negotiable.** If it's a preference,
     it belongs in Design Intent.
   - **Include a stop condition.** State when the agent should
     stop and hand off to a human or another profile.
   - **Keep it under 40 lines.** A profile is a lens, not a
     manual.
   - **Reference rules, don't repeat them.** If a rule
     covers it, point to it.

5. **Validate** — Check the profile against the repo:
   - Verify the scope matches actual file structure
   - Confirm no overlap with existing profiles
   - Check that referenced rules exist

6. **Confirm** — Show the user the created profile and ask if
   adjustments are needed.

## Example Scenario

User: "I want an agent that only handles database migrations
and schema changes."

1. Investigates repo — finds `migrations/`, Flyway config,
   existing `frontend` profile for reference
2. Asks: "Should it also handle seed data, or just schema?"
3. Creates `database.md` profile scoping to `migrations/`
   and schema files, excluding application code
4. Validates — directory exists, no overlap with frontend
   profile

## Example Output

```markdown
# Database Agent

Manages database migrations and schema changes.

## Constraints

- Follow the git-naming rule for branches and commits
- Never modify application code
- Run migrations in a transaction when supported

## Scope

Modify only:
- `migrations/`
- `schema/`

Do not touch: application code, API handlers, frontend.
If a migration requires application changes, describe
what's needed and stop.

## Design Intent

Prefer reversible migrations. Every up migration has a
corresponding down. Avoid data-destructive operations
without explicit confirmation.
```

## Common Failures

- **Scope too broad** — "handles backend" is half the repo.
  Narrow to specific directories or file types.
- **Scope too narrow** — if the agent can't do useful work
  within its boundaries, the profile is over-constrained.
- **Missing stop condition** — without "stop and hand off"
  guidance, the agent will attempt out-of-scope work rather
  than flagging it.
- **Duplicates a rule** — if the constraint applies to all
  agents, it's a rule, not a profile constraint.

## Boundaries

**Does:**
- Create profile markdown files in the project's profile directory
- Investigate repo structure to define scope
- Reference existing rules

**Does not:**
- Create skills or rules (use create-skill, create-rule)
- Modify existing profiles without confirmation
- Create profiles for non-agent roles
