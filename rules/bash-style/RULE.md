---
name: bash-style
description: >
  Bash scripting conventions based on the Google Shell Style Guide.
  Applied to all .sh files and bash scripts.
---

# Bash Style

Write bash scripts that are safe, readable, and consistent.

## Shell Options

Start every script with:

```bash
#!/usr/bin/env bash
set -euo pipefail
```

## Formatting

- Indent 2 spaces, no tabs.
- Max 80 characters per line.
- `; then` and `; do` on the same line as `if`/`for`/`while`.
- Split pipelines one per line when they don't fit:

  ```bash
  command1 \
    | command2 \
    | command3
  ```

## Naming

- Functions and variables: `lower_snake_case`
- Constants and exports: `UPPER_SNAKE_CASE`, declared with
  `readonly` or `declare -r`
- Source filenames: lowercase, underscores if needed

## Variables

- Quote all variables: `"${var}"` not `$var`.
- Use `"${var}"` with braces for all non-special variables.
- Use `local` for function variables. Separate declaration
  from command substitution assignment:

  ```bash
  local my_var
  my_var="$(some_command)"
  ```

## Functions

- Use `()` syntax, braces on the same line:

  ```bash
  my_func() {
    …
  }
  ```

- Add a header comment for non-trivial functions.
- Put all functions near the top, after constants.
- Scripts with functions must have a `main` function called
  at the bottom: `main "$@"`

## Conditionals and Tests

- Use `[[ … ]]` not `[ … ]` or `test`.
- Use `==` for string equality, not `=`.
- Use `(( … ))` for arithmetic, not `let` or `expr`.
- Use `-z`/`-n` for empty/non-empty string tests.

## Safety

- Always check return values. Use `if ! command` or check
  `$?`.
- Use arrays for argument lists, not space-separated strings.
- Use `$(command)` not backticks.
- Use `./*` not `*` for wildcard expansion.
- Never use `eval`.
- Errors go to stderr: `echo "error" >&2`

## Comments

- Every file starts with a `#` description after the shebang.
- Use `# TODO(name):` format for TODOs.

## Scope

- Keep scripts under 100 lines. If longer, consider
  rewriting in a structured language.
- Shell is for small utilities and wrappers, not complex
  logic.
