#!/usr/bin/env bash
# Create CloudFront Function for www → apex redirect
# and attach it to the distribution.
set -euo pipefail

DIST_ID="${1:?Usage: $0 <distribution-id>}"
FUNCTION_NAME="openkata-www-redirect"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Creating CloudFront Function..."
ETAG=$(aws cloudfront create-function \
  --name "${FUNCTION_NAME}" \
  --function-config '{"Comment":"Redirect www to apex","Runtime":"cloudfront-js-2.0"}' \
  --function-code "fileb://${SCRIPT_DIR}/www-redirect.js" \
  --query "ETag" \
  --output text \
  --no-cli-pager)

echo "Publishing function..."
ETAG=$(aws cloudfront publish-function \
  --name "${FUNCTION_NAME}" \
  --if-match "${ETAG}" \
  --query "ETag" \
  --output text \
  --no-cli-pager)

FUNCTION_ARN=$(aws cloudfront describe-function \
  --name "${FUNCTION_NAME}" \
  --query "FunctionSummary.FunctionMetadata.FunctionARN" \
  --output text \
  --no-cli-pager)

echo "Attaching to distribution ${DIST_ID}..."

# Get current distribution config
aws cloudfront get-distribution-config \
  --id "${DIST_ID}" \
  --no-cli-pager > /tmp/dist-config.json

DIST_ETAG=$(grep -o '"ETag": "[^"]*"' /tmp/dist-config.json | head -1 | cut -d'"' -f4)

# Add function association to default cache behavior
cat /tmp/dist-config.json \
  | python3 -c "
import json, sys
data = json.load(sys.stdin)
config = data['DistributionConfig']
config['DefaultCacheBehavior']['FunctionAssociations'] = {
    'Quantity': 1,
    'Items': [{
        'FunctionARN': '${FUNCTION_ARN}',
        'EventType': 'viewer-request'
    }]
}
json.dump(config, sys.stdout)
" > /tmp/dist-update.json

aws cloudfront update-distribution \
  --id "${DIST_ID}" \
  --distribution-config "file:///tmp/dist-update.json" \
  --if-match "${DIST_ETAG}" \
  --no-cli-pager > /dev/null

rm /tmp/dist-config.json /tmp/dist-update.json

echo ""
echo "=== Done ==="
echo "www.openkata.dev now redirects to openkata.dev"
echo "Allow a few minutes for the distribution to deploy."
