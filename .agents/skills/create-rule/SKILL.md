---
name: create-rule
version: 1.0.0
description: >
  Creates rule files (RULE.md) that define always-on constraints for
  agent sessions. Use when the user wants to create a new rule, add a
  convention, enforce a standard, define a dojo kun, or establish a
  project-wide constraint.
---

# Create Rule

Create rules (dojo kun) for this project. Rules are always-on
constraints applied to every agent session — unlike skills, they
don't need triggers.

## Workflow

1. **Gather intent** — Ask the user what the rule should enforce.
   Collect:
   - A short name (lowercase, hyphenated, e.g.,
     `markdown-consistency`)
   - What conventions or constraints to enforce
   - Any references or style guides to draw from

2. **Determine placement** — Ask whether this rule is:
   - **Local** → `.agents/rules/<name>/` — internal, not
     distributed
   - **Distributable** → `rules/<name>/` — shared via the dojo

3. **Check for conflicts** — List both `.agents/rules/` and
   `rules/` to ensure the name isn't already taken.

4. **Generate RULE.md** — Create the RULE.md in the target
   directory:

   ```markdown
   # Rule Name

   Brief description of what this rule enforces.

   ## Conventions

   - Convention one
   - Convention two

   ## References

   - Links to style guides or standards
   ```

   Follow the style of existing rules — see
   `markdown-consistency` for reference.

5. **Create CHANGELOG.md** — Start at v1.0.0 with an initial
   `### Added` entry.

6. **Symlink if distributable** — For rules in `rules/`, ask
   the user if they want it symlinked into `.agents/rules/`. If
   yes:
   ```bash
   ln -s ../../rules/<name> .agents/rules/<name>
   ```

7. **Confirm** — Show the user the created rule and ask if
   adjustments are needed.

## Conventions

- Local rules go in `.agents/rules/<name>/`
- Distributable rules go in `rules/<name>/` with a symlink in
  `.agents/rules/`
- Every rule gets a CHANGELOG.md
- Keep rules focused — one rule, one concern
