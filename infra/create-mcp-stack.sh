#!/usr/bin/env bash
# Create OpenKata MCP infrastructure from scratch.
# Idempotent — safe to re-run.
set -euo pipefail

REGION="${AWS_REGION:-eu-north-1}"
FUNCTION_NAME="openkata-mcp"
ROLE_NAME="openkata-mcp-role"
BUCKET_NAME="openkata-artifacts"
TABLE_NAME="openkata-downloads"

ACCOUNT_ID="$(aws sts get-caller-identity --query Account --output text)"

echo "=== Step 1: Create S3 bucket ==="
if aws s3api head-bucket --bucket "${BUCKET_NAME}" 2>/dev/null; then
  echo "Bucket already exists, skipping."
else
  aws s3api create-bucket \
    --bucket "${BUCKET_NAME}" \
    --region "${REGION}" \
    --create-bucket-configuration LocationConstraint="${REGION}" \
    --no-cli-pager
fi

aws s3api put-public-access-block \
  --bucket "${BUCKET_NAME}" \
  --public-access-block-configuration \
    BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true \
  --no-cli-pager

echo "=== Step 2: Create DynamoDB table ==="
if aws dynamodb describe-table --table-name "${TABLE_NAME}" --region "${REGION}" 2>/dev/null >/dev/null; then
  echo "Table already exists, skipping."
else
  aws dynamodb create-table \
    --table-name "${TABLE_NAME}" \
    --attribute-definitions AttributeName=artifact,AttributeType=S \
    --key-schema AttributeName=artifact,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region "${REGION}" \
    --no-cli-pager
fi

echo "=== Step 3: Create IAM role ==="
if aws iam get-role --role-name "${ROLE_NAME}" 2>/dev/null >/dev/null; then
  echo "Role already exists, updating policy."
else
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
fi

aws iam put-role-policy \
  --role-name "${ROLE_NAME}" \
  --policy-name openkata-mcp-access \
  --policy-document file://infra/iam-mcp-role-policy.json \
  --no-cli-pager

echo "=== Step 4: Create Lambda function ==="
GOOS=linux GOARCH=arm64 go build -o /tmp/bootstrap ./cmd/openkata-mcp/
(cd /tmp && zip -j deploy.zip bootstrap)

if aws lambda get-function --function-name "${FUNCTION_NAME}" --region "${REGION}" 2>/dev/null >/dev/null; then
  echo "Function already exists, updating code."
  aws lambda update-function-code \
    --function-name "${FUNCTION_NAME}" \
    --zip-file fileb:///tmp/deploy.zip \
    --region "${REGION}" \
    --no-cli-pager
else
  aws lambda create-function \
    --function-name "${FUNCTION_NAME}" \
    --runtime provided.al2023 \
    --architectures arm64 \
    --handler bootstrap \
    --role "arn:aws:iam::${ACCOUNT_ID}:role/${ROLE_NAME}" \
    --zip-file fileb:///tmp/deploy.zip \
    --memory-size 128 \
    --timeout 10 \
    --environment "Variables={OPENKATA_BUCKET=${BUCKET_NAME},OPENKATA_TABLE=${TABLE_NAME}}" \
    --region "${REGION}" \
    --no-cli-pager
fi

rm /tmp/bootstrap /tmp/deploy.zip

echo "=== Step 5: Create Function URL ==="
if aws lambda get-function-url-config --function-name "${FUNCTION_NAME}" --region "${REGION}" 2>/dev/null >/dev/null; then
  echo "Function URL already exists, skipping."
else
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
fi

echo ""
echo "=== Done ==="
echo "Function URL:"
aws lambda get-function-url-config \
  --function-name "${FUNCTION_NAME}" \
  --region "${REGION}" \
  --query FunctionUrl \
  --output text \
  --no-cli-pager
