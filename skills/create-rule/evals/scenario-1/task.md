# Python Documentation Convention Rule

## Problem Description

A Python team is growing rapidly and onboarding new developers every quarter. Looking through the codebase, you'll find that docstrings are written in a mix of styles: some functions use Google-style `Args:` / `Returns:` sections, others use NumPy-style with `----------` underlines, others use Sphinx-style `:param:` / `:type:` / `:rtype:` tags, and some functions use bare one-liner strings or no docstrings at all. A few utility functions use inline comments in place of docstrings.

The tech lead wants a single rule that every developer (and AI agent) can follow without ambiguity. The team has no existing docstring linter configured — they want the rule to define the standard, and they may add tooling later to enforce it.

The Python source files to use as reference are located under `inputs/utils/`.

## Output Specification

Create a rule defining the docstring convention. Place the rule in a directory named after the convention (use a lowercase hyphenated name) containing `RULE.md`.

Also produce a file named `validation-report.md` that records which files you examined and summarizes how each file's existing docstrings compare to the conventions you wrote.
