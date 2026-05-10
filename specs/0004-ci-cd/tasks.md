---
spec: 0004-ci-cd
---

# Tasks

### 1. Create build workflow
- **Status**: Pending
- **Goal**: GitHub Actions workflow that runs on PR: install
  Go, install templ, generate templates, build binary
- **Boundary**: `.github/workflows/build.yaml`
- **Depends**: None
- **Done when**: Workflow file exists and would compile the
  project on PR

### 2. Create deploy workflow
- **Status**: Pending
- **Goal**: GitHub Actions workflow that runs on merge to
  main: build binary, assume OIDC role, push to Lambda
- **Boundary**: `.github/workflows/deploy.yaml`
- **Depends**: None
- **Done when**: Workflow file exists with OIDC auth and
  deploy steps

### 3. Document OIDC setup instructions
- **Status**: Pending
- **Goal**: Step-by-step instructions for creating the AWS
  OIDC provider and IAM role in the console, plus GitHub
  repo settings for branch protection
- **Boundary**: `infra/README.md`
- **Depends**: 2
- **Done when**: README has clear manual steps for OIDC
  provider, IAM role, and branch protection

## Progress Log
