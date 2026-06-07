# Rationale

## Why a mandatory intake gate

Jumping to writing without understanding requirements
produces generic skills that don't match the user's
actual workflow. Three questions minimum forces the
agent to understand before acting.

## Why descriptions optimize for activation

LLMs use the description field for routing — deciding
which skill to activate. A description that teaches
the workflow instead of describing triggers causes
mis-activation. The body teaches; the description
matches.

## Why ban version in frontmatter

The Agent Skills spec doesn't include version in
frontmatter. It belongs in metadata. Putting it in
the wrong place triggers spec validation warnings and
confuses tooling about the source of truth.

## Why validate with negative prompts

A skill that activates on everything is useless.
Negative prompts ("this should NOT trigger the skill")
test that the trigger surface has boundaries. Without
them, over-eager activation degrades the full system.
