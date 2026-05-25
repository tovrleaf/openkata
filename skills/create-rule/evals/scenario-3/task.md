# Django ORM Query Convention Rule

## Problem Description

A backend engineering team maintains a Django e-commerce application with an increasingly complex data layer. Over time, developers have accumulated different habits around database access: some write raw SQL for perceived performance gains, others use the ORM inconsistently and trigger N+1 query problems by neglecting `select_related` and `prefetch_related`. There's no agreed standard on when raw SQL is acceptable versus when the ORM should be used.

The engineering manager wants a concise, enforceable rule covering how database queries should be written — one that every developer and AI agent can follow immediately without having to interpret intent. The existing source files in `inputs/app/` show the kind of patterns currently in use across the project.

The rule should be practical and focused on enforcement rather than education. If you need to capture extended reasoning, examples, or background on ORM methods, put that material where it won't bloat the main rule file.

## Output Specification

Create a Django ORM query convention rule. Place it in a lowercase-hyphenated directory under `.agents/rules/` containing a `RULE.md`. 

Also produce a `validation-report.md` at the top level that records: which files you examined, whether each file's patterns conform to or conflict with the conventions you wrote, and any edge cases you observed.
