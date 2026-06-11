# Rationale

git-naming defines branch naming, commit message format,
and AI attribution trailer conventions.

## Why every AI-assisted commit needs an Assisted-by trailer

Attribution in git history makes AI contributions
auditable. Teams can filter, review, or analyze
AI-assisted work separately. Without attribution,
it's invisible in the history.

## Why commit bodies are only required for non-trivial changes

Forcing a body on "fix typo" commits creates empty
busywork. Requiring it for trade-offs and decisions
ensures bodies exist where they add value. The rule
is: if someone would ask "why?", include a body.

## Why squash merges use the PR title, not concatenated commits

Concatenated commit messages produce unreadable
merge commits. The PR title is already a reviewed
summary. Individual commits are preserved in branch
history for anyone who needs detail.
