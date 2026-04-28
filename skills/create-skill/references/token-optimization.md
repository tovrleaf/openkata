# Token Optimization

Use this after the skill works. Optimize for lower context cost
without reducing execution quality.

## Keep in SKILL.md

- The trigger-bearing frontmatter
- The core workflow
- Critical decision rules
- Short examples that anchor the workflow
- Direct links to bundled references

## Move Out of SKILL.md

- Long domain primers
- Exhaustive edge-case catalogs
- Variant-specific instructions
- Large examples
- Detailed command references
- Documentation discoverable from the repo at runtime

## Compression Rules

- Delete repeated ideas before rewriting sentences
- Prefer short checklists over explanatory paragraphs
- Replace generic advice with workflow-specific rules
- Keep examples only if they teach something not already
  obvious from the instructions
- Avoid motivational or narrative text
- Prefer one sharp sentence over two soft ones

## Smell Tests

The main file is probably too large if:

- Multiple sections repeat the same workflow in different words
- Examples are longer than the instructions they illustrate
- Reference material dominates the core procedure
- The skill explains common concepts instead of workflow-specific
  guidance

## Final Pass

Ask:

1. What text can be deleted with no loss of behavior?
2. What text belongs in `references/`?
3. What assumptions should be stated once instead of repeated?
4. Is the description still strong enough to trigger correctly
   after trimming?
