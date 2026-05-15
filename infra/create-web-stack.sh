#!/usr/bin/env bash
# Create OpenKata infrastructure from scratch.
# Run once. After this, use `make deploy` for code updates.
set -euo pipefail

REGION="${AWS_REGION:-eu-north-1}"
FUNCTION_NAME="openkata-web"
ROLE_NAME="openkata-web-role"

echo "=== Step 1: Create IAM role ==="
aws iam create-role \
  --role-name "${ROLE_NAME}" \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [{
      "Effect": "Allow",
      "Principal": {"Service": "lambda.amazonaws.com"},
      "Action": "sts:AssumeRole"
    }]
  }' \
  --no-cli-pager

aws iam attach-role-policy \
  --role-name "${ROLE_NAME}" \
  --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

echo "Waiting for role to propagate..."
sleep 10

echo "=== Step 2: Build Lambda binary ==="
"$(go env GOPATH)/bin/templ" generate ./cmd/openkata-web/templates/
GOOS=linux GOARCH=arm64 go build -o /tmp/bootstrap ./cmd/openkata-web/
(cd /tmp && zip -j deploy.zip bootstrap)

echo "=== Step 3: Create Lambda function ==="
ACCOUNT_ID="$(aws sts get-caller-identity --query Account --output text)"
aws lambda create-function \
  --function-name "${FUNCTION_NAME}" \
  --runtime provided.al2023 \
  --architectures arm64 \
  --handler bootstrap \
  --role "arn:aws:iam::${ACCOUNT_ID}:role/${ROLE_NAME}" \
  --zip-file fileb:///tmp/deploy.zip \
  --memory-size 128 \
  --timeout 10 \
  --environment "Variables={OPENKATA_BUCKET=openkata-artifacts,OPENKATA_TABLE=openkata-downloads}" \
  --region "${REGION}" \
  --no-cli-pager

rm /tmp/bootstrap /tmp/deploy.zip

echo "=== Step 4: Create Function URL ==="
aws lambda create-function-url-config \
  --function-name "${FUNCTION_NAME}" \
  --auth-type NONE \
  --region "${REGION}" \
  --no-cli-pager

aws lambda add-permission \
  --function-name "${FUNCTION_NAME}" \
  --statement-id FunctionURLPublic \
  --action lambda:InvokeFunctionUrl \
  --principal "*" \
  --function-url-auth-type NONE \
  --region "${REGION}" \
  --no-cli-pager

aws lambda add-permission \
  --function-name "${FUNCTION_NAME}" \
  --statement-id FunctionURLInvoke \
  --action lambda:InvokeFunction \
  --principal "*" \
  --region "${REGION}" \
  --no-cli-pager

echo ""
echo "=== Done ==="
echo "Function URL:"
aws lambda get-function-url-config \
  --function-name "${FUNCTION_NAME}" \
  --region "${REGION}" \
  --query FunctionUrl \
  --output text \
  --no-cli-pager
