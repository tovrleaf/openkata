---
name: openkata-rule-conventions
version: 1.0.0
description: >
  Applies OpenKata project conventions when creating or updating
  rules in this repository. Adds changelog management, symlink
  placement, and acknowledgments on top of the generic create-rule
  workflow. Use alongside create-rule when working in the OpenKata
  repo.
---

# OpenKata Rule Conventions

Project-specific additions to the `create-rule` workflow. This
skill activates alongside the generic `create-rule` and adds
OpenKata conventions.

## Additional Steps

After the generic `create-rule` workflow, also do:

1. **Determine placement** — Ask whether this rule is:
   - **Local** → `.agents/rules/<name>/`
   - **Distributable** → `rules/<name>/`

2. **Create CHANGELOG.md** — Every rule gets a changelog
   starting at v1.0.0 with an initial `### Added` entry.
   Follow the markdown-consistency rule for formatting.

3. **Symlink if distributable** — For rules in `rules/`, ask
   the user if they want it symlinked into `.agents/rules/`:
   ```bash
   ln -s ../../rules/<name> .agents/rules/<name>
   ```

4. **Acknowledge sources** — If the rule draws on external
   style guides or standards, create
   `references/ACKNOWLEDGMENTS.md` listing each source with a
   link, license, what was adapted, and the version it was
   adopted in.

## Conventions

- Local rules go in `.agents/rules/<name>/`
- Distributable rules go in `rules/<name>/` with a symlink
  in `.agents/rules/`
- Every rule gets a CHANGELOG.md
- Changelogs document rule-facing changes only
- Keep rules focused — one rule, one concern
