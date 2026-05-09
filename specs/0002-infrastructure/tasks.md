---
spec: 0002-infrastructure
---

# Tasks

### 1. Adapt web server for Lambda
- **Status**: Pending
- **Goal**: Detect Lambda environment and use the Lambda
  adapter; embed static assets with go:embed
- **Boundary**: `cmd/openkata-web/`
- **Depends**: None
- **Done when**: Binary builds for linux/arm64, serves
  locally unchanged, and handles Lambda events when
  AWS_LAMBDA_FUNCTION_NAME is set

### 2. Initialize CDK project
- **Status**: Pending
- **Goal**: Scaffold CDK Go project in `infra/`
- **Boundary**: `infra/`
- **Depends**: None
- **Done when**: `cd infra && cdk synth` produces a
  CloudFormation template

### 3. Define Lambda stack
- **Status**: Pending
- **Goal**: CDK stack with Lambda function (arm64,
  provided.al2023) and Function URL
- **Boundary**: `infra/`
- **Depends**: 1, 2
- **Done when**: `cdk synth` includes Lambda function and
  Function URL resources

### 4. Create deploy workflow
- **Status**: Pending
- **Goal**: `make deploy` builds the binary, packages it,
  and runs `cdk deploy`
- **Boundary**: `Makefile`, `scripts/deploy.sh`
- **Depends**: 3
- **Done when**: `make deploy` runs end-to-end (with
  configured AWS credentials)

## Progress Log
