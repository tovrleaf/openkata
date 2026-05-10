---
title: CI/CD with GitHub Actions
status: Draft
depth: Standard
created: 2026-05-10
---

# CI/CD with GitHub Actions

## Goal

Automate build verification on PRs and deploy on merge to
main. Protect main branch from direct pushes.

## Requirements

1. GitHub Actions workflow: build on PR (compile + templ)
2. GitHub Actions workflow: deploy on merge to main
3. AWS OIDC identity provider for GitHub Actions (no stored
   credentials)
4. AWS IAM role for GitHub Actions with iam-ci-policy.json
5. Main branch protected: require PR with review
6. Agent does not run AWS commands; provides instructions
   for manual setup where needed

## Out of Scope

- Tests (no test suite yet)
- Linting
- Skill review in CI
- Cache invalidation after deploy

## Auth Model

GitHub Actions assumes an AWS role via OIDC. The role is
scoped to `lambda:UpdateFunctionCode` only (iam-ci-policy).
No long-lived credentials stored in GitHub.

## Acceptance Criteria

- PR triggers build workflow, fails if code doesn't compile
- Merge to main triggers deploy, site updates automatically
- Direct push to main is blocked
- No AWS access keys stored anywhere
