# Tasks: Infrastructure for Lambda Deployment

## Tasks

### 1. Adapt web server for Lambda
- **Status**: Done
- **Goal**: Detect Lambda environment and use the Lambda
  adapter; embed static assets with go:embed
- **Boundary**: `cmd/openkata-web/`
- **Depends**: None
- **Done when**: Binary builds for linux/arm64, serves
  locally unchanged, and handles Lambda events when
  AWS_LAMBDA_FUNCTION_NAME is set

### 2. Create CloudFormation template
- **Status**: Done
- **Goal**: Template with Lambda function (arm64,
  provided.al2023), Function URL, IAM role, and S3 bucket
  for deploy artifacts
- **Boundary**: `infra/template.yaml`
- **Depends**: 1
- **Done when**: `aws cloudformation validate-template`
  passes

### 3. Create deploy workflow
- **Status**: Done
- **Goal**: `make deploy` builds the binary, zips it,
  uploads to S3, and runs `aws cloudformation deploy`
- **Boundary**: `Makefile`, `scripts/deploy.sh`
- **Depends**: 2
- **Done when**: `make deploy` runs end-to-end (with
  configured AWS credentials)

## Progress Log
