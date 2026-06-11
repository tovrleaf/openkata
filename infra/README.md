# Infrastructure

AWS infrastructure for Open Kata (website + MCP server).

## Prerequisites

- AWS CLI configured (`aws configure`)
- IAM user with [iam-admin-policy.json](iam-admin-policy.json)

## IAM Policies

| File | Purpose | Apply to |
|------|---------|----------|
| [iam-admin-policy.json](iam-admin-policy.json) | Create/destroy infra, deploy | Your IAM user |
| [iam-ci-policy.json](iam-ci-policy.json) | Deploy code, publish skills to S3 | `openkata-ci` role |
| [iam-ci-trust-policy.json](iam-ci-trust-policy.json) | OIDC trust for GitHub Actions | `openkata-ci` role (trust relationships) |
| [iam-mcp-role-policy.json](iam-mcp-role-policy.json) | S3 read, DynamoDB access | `openkata-mcp-role` (Lambda execution) |
| [iam-web-role-mcp-policy.json](iam-web-role-mcp-policy.json) | S3 read, DynamoDB access | `openkata-web-role` (Lambda execution) |

## Create infrastructure (one-time)

Run in order:

```bash
# 1. Web server (Lambda + Function URL)
./infra/create-web-stack.sh

# 2. MCP server (Lambda + S3 + DynamoDB + Function URL)
./infra/create-mcp-stack.sh

# 3. Add MCP permissions to web Lambda role
#    Console: IAM → Roles → openkata-web-role → Add inline policy
#    Paste contents of iam-web-role-mcp-policy.json
#    Name: openkata-web-mcp-access

# 4. Update CI role
#    Console: IAM → Roles → openkata-ci → Update inline policy
#    Replace with contents of iam-ci-policy.json
```

## Deploy code

```bash
make deploy          # web server
make deploy-mcp      # MCP server (TODO)
```

## CI/CD (GitHub Actions)

### OIDC Setup (one-time, in AWS console)

1. IAM → Identity providers → Add provider
   - Type: OpenID Connect
   - URL: `https://token.actions.githubusercontent.com`
   - Audience: `sts.amazonaws.com`

2. IAM → Roles → Create role
   - Trusted entity: Custom trust policy
   - Paste contents of [iam-ci-trust-policy.json](iam-ci-trust-policy.json)
   - Role name: `openkata-ci`

3. Attach inline policy to the role:
   - Paste contents of [iam-ci-policy.json](iam-ci-policy.json)
   - Name: `openkata-deploy`
