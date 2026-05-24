# Standardize API Error Responses

## Problem Description

A TypeScript Express service has grown organically across several developers. Looking at the route handlers, you notice that error responses returned to clients use inconsistent shapes: some use `{ message }`, others use `{ msg }`, `{ error }`, `{ err }`, or `{ status, message }`. Success responses are similarly inconsistent — some wrap data in `{ data }`, others in `{ result }`, others return the object directly.

The team wants a single, project-wide rule so that all current and future handlers return responses in a predictable shape. This matters for frontend consumers who parse these responses and for monitoring tools that inspect error payloads.

The existing TypeScript handler files are under `inputs/src/handlers/`. There is also a linting configuration at `inputs/.eslintrc.json` that is already in use.

## Output Specification

Create a rule for API response conventions. Place the rule in a directory named `api-response-shape/` containing `RULE.md`.

Also produce a file named `validation-report.md` that records which existing files you examined and what you found.
