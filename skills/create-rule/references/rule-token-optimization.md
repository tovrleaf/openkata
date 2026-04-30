# Rule Token Optimization

Rules load every session. Every line costs tokens on every
interaction, even when the rule is irrelevant to the current
task. Optimize aggressively.

## Keep in RULE.md

- Enforceable conventions (the core content)
- Section headings that group by concern
- Exceptions that change behavior

## Move to references/

- Rationale for why a convention exists
- Extended examples
- Links to external style guides
- Historical context

## Compression Rules

- State conventions, never explain them
- One convention per bullet — no compound rules
- Delete any line an agent could infer from context
- Prefer "Use X" over "You should use X"
- Remove motivational or narrative text

## Smell Tests

The rule is probably too large if:

- It exceeds 100 lines
- Multiple bullets say the same thing differently
- It explains concepts the agent already knows
- It includes examples for self-evident conventions
- Sections overlap in scope (merge them)

## Final Pass

Ask:

1. Can an agent follow every convention literally?
2. What text can be deleted with no loss of enforceability?
3. What text is rationale disguised as convention?
4. Is any convention already covered by another rule?
