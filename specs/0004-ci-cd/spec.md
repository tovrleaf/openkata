---
status: Draft
depth: Standard
---

# CI/CD with GitHub Actions

## Story

As a maintainer, I want automated build checks on PRs and
deploy on merge so that the site stays healthy and deploys
without manual steps.

## Requirements

- GitHub Actions workflow: build on PR (compile + templ)
- GitHub Actions workflow: deploy on merge to main
- AWS OIDC identity provider for GitHub Actions (no stored
  credentials)
- AWS IAM role for GitHub Actions with iam-ci-policy.json
- Main branch protected: require PR with review
- Agent does not run AWS commands; provides instructions
  for manual setup where needed

## Out of Scope

- Tests (no test suite yet)
- Linting
- Skill review in CI
- Cache invalidation after deploy

## Open Questions

- None
