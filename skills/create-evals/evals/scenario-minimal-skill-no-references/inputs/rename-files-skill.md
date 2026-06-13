---
name: rename-files
description: >
  Batch-renames files using glob patterns and rename
  expressions. Use when the user says "rename files",
  "bulk rename", or has files that need consistent naming.
metadata:
  version: "1.0.0"
  tags: "category:utilities"
---

# Rename Files

Batch-rename files matching a pattern.

## Workflow

1. **Detect files** — Find files matching the user's glob
   pattern. Show the count and list.

2. **Preview** — Show a table of current → new filenames.
   Do not apply until confirmed.

3. **Apply** — Rename files. Report successes and any
   failures (permission errors, name conflicts).

## Conventions

- Preserve file extensions unless explicitly asked to change
- Use `mv` not `cp + rm`
- Never overwrite existing files without confirmation

## Boundaries

- DOES rename files in the current directory tree
- Does NOT modify file contents
- Does NOT create directories

## Common Failures

- **Overwriting files** — two source files mapping to the
  same target name causes data loss. Always check first.
- **Missing preview** — applying renames without showing
  the plan first.
