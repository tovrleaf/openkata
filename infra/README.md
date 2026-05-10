# Infrastructure

CloudFormation template for the Open Kata website.

## Prerequisites

- AWS CLI configured (`aws configure`)
- IAM role with the policy below

## IAM Policies

Two policies for different contexts:

- [iam-admin-policy.json](iam-admin-policy.json) — for your
  machine. Creates/destroys infrastructure and deploys code.
- [iam-ci-policy.json](iam-ci-policy.json) — for GitHub
  Actions. Only deploys code to an existing Lambda.

## Create infrastructure (one-time)

```bash
./infra/create-stack.sh
```

## Deploy code

```bash
make deploy
```

## Get the URL

```bash
aws cloudformation describe-stacks \
  --stack-name openkata-web \
  --region eu-north-1 \
  --query "Stacks[0].Outputs[?OutputKey=='FunctionUrl'].OutputValue" \
  --output text
```

## Destroy

```bash
aws cloudformation delete-stack \
  --stack-name openkata-web \
  --region eu-north-1
```

## CI/CD (GitHub Actions)

### OIDC Setup (one-time, in AWS console)

1. IAM → Identity providers → Add provider
   - Type: OpenID Connect
   - URL: `https://token.actions.githubusercontent.com`
   - Audience: `sts.amazonaws.com`

2. IAM → Roles → Create role
   - Trusted entity: Web identity
   - Provider: `token.actions.githubusercontent.com`
   - Organization: `tovrleaf`
   - Repository: `openkata`
   - Branch: `main`
   - Role name: `openkata-ci`

3. Attach inline policy to the role:
   - Paste contents of [iam-ci-policy.json](iam-ci-policy.json)
   - Name: `openkata-deploy`

### Branch Protection (in GitHub)

1. Settings → Branches → Add rule
   - Branch name pattern: `main`
   - Require pull request before merging
   - Require approvals: 1
   - Do not allow bypassing
