# Example Rule

A finished rule for markdown consistency:

```markdown
# Markdown Consistency

Enforces consistent markdown formatting across all generated
and edited `.md` files.

## Lists

- Use `-` for unordered list items, never `*` or `+`.
- Use `1.` for every ordered list item (not incrementing
  numbers).
- Indent nested list items with 2 spaces.

## Headings

- Use ATX-style headings (`#`), never setext-style (`---`).
- One blank line before and after every heading.
- Do not skip heading levels (e.g., no `###` directly
  under `#`).

## Code Blocks

- Use fenced code blocks (` ``` `) with an explicit language
  tag.
- Never use indented code blocks.

## Line Length

- Wrap prose at 72 characters.
- Do not wrap code blocks, tables, or URLs.

## Emphasis

- Use `**bold**` for emphasis, never `__bold__`.
- Use `*italic*` for secondary emphasis, never `_italic_`.
```

This example shows:

- Every convention is specific and literally enforceable
- Conventions are stated as directives, not explanations
- Grouped by concern with clear section headings
- No rationale — just what to do
