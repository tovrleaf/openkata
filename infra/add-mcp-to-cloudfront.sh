#!/usr/bin/env bash
# Add MCP Lambda as a second origin to the CloudFront distribution.
# Routes /mcp* requests to the MCP Lambda Function URL.
# Run once. Idempotent — safe to re-run.
set -euo pipefail

REGION="${AWS_REGION:-eu-north-1}"
DOMAIN="openkata.dev"

# Get MCP Lambda Function URL origin domain
MCP_URL=$(aws lambda get-function-url-config \
  --function-name openkata-mcp \
  --region "${REGION}" \
  --query FunctionUrl \
  --output text \
  --no-cli-pager)
MCP_ORIGIN=$(echo "${MCP_URL}" | sed 's|https://||;s|/$||')

# Get the CloudFront distribution ID for openkata.dev
DIST_ID=$(aws cloudfront list-distributions \
  --query "DistributionList.Items[?Aliases.Items[?contains(@,'${DOMAIN}')]].Id | [0]" \
  --output text \
  --no-cli-pager)

if [[ -z "${DIST_ID}" || "${DIST_ID}" == "None" ]]; then
  echo "Error: No CloudFront distribution found for ${DOMAIN}"
  exit 1
fi

echo "Distribution: ${DIST_ID}"
echo "MCP origin: ${MCP_ORIGIN}"

# Get current config
CONFIG=$(aws cloudfront get-distribution-config \
  --id "${DIST_ID}" \
  --no-cli-pager)

ETAG=$(echo "${CONFIG}" | jq -r '.ETag')
DIST_CONFIG=$(echo "${CONFIG}" | jq '.DistributionConfig')

# Check if mcp origin already exists
if echo "${DIST_CONFIG}" | jq -e '.Origins.Items[] | select(.Id == "mcp")' >/dev/null 2>&1; then
  echo "MCP origin already exists, skipping."
  exit 0
fi

# Add MCP origin
DIST_CONFIG=$(echo "${DIST_CONFIG}" | jq --arg domain "${MCP_ORIGIN}" '
  .Origins.Items += [{
    "Id": "mcp",
    "DomainName": $domain,
    "CustomOriginConfig": {
      "HTTPPort": 80,
      "HTTPSPort": 443,
      "OriginProtocolPolicy": "https-only",
      "OriginSslProtocols": { "Quantity": 1, "Items": ["TLSv1.2"] },
      "OriginReadTimeout": 30,
      "OriginKeepaliveTimeout": 5
    },
    "OriginPath": "",
    "CustomHeaders": { "Quantity": 0 },
    "ConnectionAttempts": 3,
    "ConnectionTimeout": 10,
    "OriginShield": { "Enabled": false }
  }] | .Origins.Quantity += 1')

# Add /mcp* cache behavior (no caching, forward all)
DIST_CONFIG=$(echo "${DIST_CONFIG}" | jq '
  .CacheBehaviors.Items = ((.CacheBehaviors.Items // []) + [{
    "PathPattern": "/mcp*",
    "TargetOriginId": "mcp",
    "ViewerProtocolPolicy": "redirect-to-https",
    "AllowedMethods": {
      "Quantity": 7,
      "Items": ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"],
      "CachedMethods": { "Quantity": 2, "Items": ["GET", "HEAD"] }
    },
    "CachePolicyId": "4135ea2d-6df8-44a3-9df3-4b5a84be39ad",
    "OriginRequestPolicyId": "b689b0a8-53d0-40ab-baf2-68738e2966ac",
    "Compress": true,
    "SmoothStreaming": false,
    "FieldLevelEncryptionId": "",
    "LambdaFunctionAssociations": { "Quantity": 0 },
    "FunctionAssociations": { "Quantity": 0 }
  }]) | .CacheBehaviors.Quantity = (.CacheBehaviors.Items | length)')

echo "Updating distribution..."
aws cloudfront update-distribution \
  --id "${DIST_ID}" \
  --if-match "${ETAG}" \
  --distribution-config "${DIST_CONFIG}" \
  --no-cli-pager > /dev/null

echo "Done. /mcp* now routes to MCP Lambda."
echo "Allow 5-10 minutes for CloudFront to deploy."
