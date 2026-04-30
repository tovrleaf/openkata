# Markdown Consistency

Apply these conventions to all markdown files.

## Document Structure

- Start with a single `#` heading as the document title
- Use ATX-style headings (`#`) only — never setext underlines (`===`, `---`)
- Increment heading levels by one — never skip levels
- No trailing punctuation in headings
- Surround headings with blank lines (one before, one after)

## Line Length

- Wrap prose at 72–80 characters
- Exceptions: links, tables, headings, and code blocks may exceed this

## Spacing

- One blank line before and after headings, lists, code blocks, and tables
- No multiple consecutive blank lines
- No trailing whitespace
- Files end with a single newline

## Emphasis

- `**bold**` for key terms and warnings — never `__bold__`
- `*italic*` for emphasis and titles — never `_italic_`
- `` `backticks` `` for code, file names, commands, and field names
- Don't bold entire sentences or paragraphs

## Lists

- Use `-` for unordered lists — never `*` or `+`
- Use `1.` for all ordered list items (lazy numbering) in long lists
- Use sequential numbering (`1.`, `2.`, `3.`) for short, stable lists
- One space after the list marker
- Indent nested lists by 2 spaces
- Surround lists with blank lines

## Code

- Use fenced code blocks (triple backticks) — never indented code blocks
- Always specify the language after opening backticks
- Use `text` as language when no specific language applies
- Use inline backticks for short code, file names, commands, and variables
- Don't prefix shell commands with `$` unless showing output alongside

## Links

- Use descriptive link text — never "click here" or "link"
- Use reference-style links when the URL would break line length or
  when the same URL appears multiple times
- Use relative paths for links within the same repository

## Tables

- Use tables for tabular data that benefits from scanning — not for
  layout or prose
- Align columns with pipes and hyphens for source readability
- Keep cells concise — use reference links for long URLs in tables
- Prefer lists over tables when data doesn't have uniform structure

## Images

- Always include descriptive alt text
- Use relative paths for local images

## General

- Prefer markdown over inline HTML
- No hard tabs — use spaces
- Wrap prose at 72–80 characters
- One sentence per line is allowed when the file is primarily
  diff-reviewed, but wrap at 72–80 characters by default
