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
