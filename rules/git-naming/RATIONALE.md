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

## Why PRs are never squash-merged

Squashing replaces many described commits with one, losing
the per-step reasoning, spec references, and trailers.
It also hides branch ancestry: `git branch --merged` and
`git branch -d` cannot tell that a squashed branch was
merged, so cleanup needs extra checks. Merge commits and
rebases keep every commit and its message intact.
