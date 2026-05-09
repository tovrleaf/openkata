#!/usr/bin/env bash
set -euo pipefail

FUNCTION_NAME="openkata-web"
REGION="${AWS_REGION:-eu-north-1}"

echo "Generating templates..."
"$(go env GOPATH)/bin/templ" generate ./cmd/openkata-web/templates/

BUILD_DIR="$(mktemp -d)"
trap 'rm -rf "${BUILD_DIR}"' EXIT

echo "Building for Lambda (linux/arm64)..."
GOOS=linux GOARCH=arm64 go build -o "${BUILD_DIR}/bootstrap" ./cmd/openkata-web/

echo "Packaging..."
zip -j "${BUILD_DIR}/deploy.zip" "${BUILD_DIR}/bootstrap"

echo "Deploying to ${FUNCTION_NAME} in ${REGION}..."
aws lambda update-function-code \
  --function-name "${FUNCTION_NAME}" \
  --zip-file "fileb://${BUILD_DIR}/deploy.zip" \
  --region "${REGION}" \
  --no-cli-pager
echo "Done."
