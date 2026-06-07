# Rationale

## Why set -euo pipefail is mandatory

Without it, scripts continue after failures silently.
Debugging a script that failed on line 12 but ran to
line 50 is painful. Fail-fast makes errors immediate
and locatable.

## Why ban eval

eval executes arbitrary strings as code. In scripts
that handle user input or variable paths, eval creates
injection vulnerabilities. The risk never justifies
the convenience.

## Why 100-line limit

Bash lacks structured error handling, type safety, and
testability. Past 100 lines, these gaps become bugs.
Rewriting in a structured language at that point is
cheaper than debugging bash edge cases.

## Why separate local declaration from assignment

`local x=$(cmd)` masks the return code of `cmd`
because `local` always returns 0. Separating them
(`local x; x=$(cmd)`) lets `set -e` catch failures.
Subtle but causes real production bugs.
