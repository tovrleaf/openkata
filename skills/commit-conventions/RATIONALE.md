# Rationale

## Why check for existing conventions first

Repos may already have `.commitlintrc` or established
formats. Overwriting them causes friction with existing
CI. Check first, defer to what exists.

## Why ban git add .

Bulk staging hides unrelated changes in a commit.
Forcing file-by-file staging makes the agent conscious
of what each commit contains. This prevents accidental
secret commits and mixed-concern diffs.

## Why validate after committing

Post-commit validation (`git log -1`) creates a
self-correction loop. If the message is malformed,
the agent can amend immediately. Catching errors
after the fact is cheaper than preventing them with
complex pre-validation.

## Why delegate format rules to git-naming

The skill handles workflow (when/how to commit), the
rule handles format (what the message looks like).
Separation means updating format conventions doesn't
require touching the workflow skill.
