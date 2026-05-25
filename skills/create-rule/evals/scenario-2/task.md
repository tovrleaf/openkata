# TypeScript Naming Convention Rule

## Problem Description

A TypeScript frontend team has been working on a large React application for two years. Code reviews increasingly surface disagreements about naming: some developers use `camelCase` for utility functions and `PascalCase` for React component files, while others apply no consistent rule. Module file naming is also inconsistent — some files use kebab-case (`user-profile.ts`), others use camelCase (`userProfile.ts`), and React component files sometimes use either convention.

The tech lead wants a single project-wide rule to settle these disputes once and for all. There's a `src/` folder with various TypeScript files, and the project has been building up rules and tooling config in `.agents/rules/` over time.

## Output Specification

Create a naming convention rule and place it under `.agents/rules/` in a new directory containing a `RULE.md`.

Also produce a file named `validation-report.md` at the top level documenting: which existing rules and configurations you found, which naming concerns you chose to include in the new rule and why, and any edge cases or exceptions you identified.
