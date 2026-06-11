# Rationale

create-rule builds always-on constraint files that
apply to every agent session in the project.

## Why rules must be enforceable by reading literally

If an agent must interpret intent, the rule is too
vague. Rules that require judgment produce inconsistent
results across sessions. Literal enforceability means
any agent applies it the same way.

## Why token cost awareness matters for rules

Rules load every session regardless of relevance.
Every line costs tokens. A rule that explains "why"
wastes budget on context the agent doesn't need to
follow the instruction. Rationale goes in RATIONALE.md.

## Why rules are validated against existing files

A new rule might contradict existing code patterns or
other rules. Checking 2–3 existing files before
committing catches conflicts early — before the rule
causes confusing agent behavior.

## Why rules state conventions, not explanations

Agents don't need to understand why. They need to know
what. "Use 2-space indent" is actionable. "Indent with
2 spaces because it improves readability" wastes tokens
on justification the agent can't use.
