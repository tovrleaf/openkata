#!/usr/bin/env bash
# Cache /static/* at the CloudFront edge.
# Creates the openkata-static cache policy and a /static/* behavior on the
# web Lambda origin. The web server sets Cache-Control on static files, and a
# URL with the file's hash in ?v= is cached for a year.
# Run once, after a deploy that sets those headers. Idempotent — safe to
# re-run.
set -euo pipefail

DOMAIN="openkata.dev"
POLICY_NAME="openkata-static"

# Get the CloudFront distribution ID for openkata.dev
DIST_QUERY="DistributionList.Items"
DIST_QUERY+="[?Aliases.Items[?contains(@,'${DOMAIN}')]].Id | [0]"
DIST_ID=$(aws cloudfront list-distributions \
  --query "${DIST_QUERY}" \
  --output text \
  --no-cli-pager)

if [[ -z "${DIST_ID}" || "${DIST_ID}" == "None" ]]; then
  echo "Error: No CloudFront distribution found for ${DOMAIN}"
  exit 1
fi

echo "Distribution: ${DIST_ID}"

# Find or create the cache policy. Only the query string v is in the cache
# key: no headers, so the Function URL never sees the viewer's Host, and no
# cookies. TTLs come from the response's Cache-Control.
POLICY_QUERY="CachePolicyList.Items"
POLICY_QUERY+="[?CachePolicy.CachePolicyConfig.Name=='${POLICY_NAME}']"
POLICY_QUERY+=".CachePolicy.Id | [0]"
POLICY_ID=$(aws cloudfront list-cache-policies \
  --type custom \
  --query "${POLICY_QUERY}" \
  --output text \
  --no-cli-pager)

if [[ -z "${POLICY_ID}" || "${POLICY_ID}" == "None" ]]; then
  echo "Creating cache policy ${POLICY_NAME}..."
  POLICY_ID=$(aws cloudfront create-cache-policy \
    --cache-policy-config "{
      \"Name\": \"${POLICY_NAME}\",
      \"Comment\": \"Static files: cache key is the v query string\",
      \"MinTTL\": 0,
      \"DefaultTTL\": 300,
      \"MaxTTL\": 31536000,
      \"ParametersInCacheKeyAndForwardedToOrigin\": {
        \"EnableAcceptEncodingGzip\": true,
        \"EnableAcceptEncodingBrotli\": true,
        \"HeadersConfig\": { \"HeaderBehavior\": \"none\" },
        \"CookiesConfig\": { \"CookieBehavior\": \"none\" },
        \"QueryStringsConfig\": {
          \"QueryStringBehavior\": \"whitelist\",
          \"QueryStrings\": { \"Quantity\": 1, \"Items\": [\"v\"] }
        }
      }
    }" \
    --query CachePolicy.Id \
    --output text \
    --no-cli-pager)
else
  echo "Cache policy ${POLICY_NAME} already exists."
fi

echo "Cache policy: ${POLICY_ID}"

# Get current config
CONFIG=$(aws cloudfront get-distribution-config \
  --id "${DIST_ID}" \
  --no-cli-pager)

ETAG=$(echo "${CONFIG}" | jq -r '.ETag')
DIST_CONFIG=$(echo "${CONFIG}" | jq '.DistributionConfig')

# Check if the /static/* behavior already exists
if echo "${DIST_CONFIG}" | jq -e \
  '.CacheBehaviors.Items[]? | select(.PathPattern == "/static/*")' \
  >/dev/null 2>&1; then
  echo "/static/* behavior already exists, skipping."
  exit 0
fi

# Add /static/* cache behavior on the web Lambda origin. It keeps the
# default behavior's CloudFront functions, so www still redirects to the
# apex.
DIST_CONFIG=$(echo "${DIST_CONFIG}" | jq --arg policy "${POLICY_ID}" '
  .DefaultCacheBehavior.FunctionAssociations as $functions |
  .CacheBehaviors.Items = ((.CacheBehaviors.Items // []) + [{
    "PathPattern": "/static/*",
    "TargetOriginId": "lambda",
    "ViewerProtocolPolicy": "redirect-to-https",
    "AllowedMethods": {
      "Quantity": 2,
      "Items": ["GET", "HEAD"],
      "CachedMethods": { "Quantity": 2, "Items": ["GET", "HEAD"] }
    },
    "CachePolicyId": $policy,
    "Compress": true,
    "SmoothStreaming": false,
    "FieldLevelEncryptionId": "",
    "LambdaFunctionAssociations": { "Quantity": 0 },
    "FunctionAssociations": ($functions // { "Quantity": 0 })
  }]) | .CacheBehaviors.Quantity = (.CacheBehaviors.Items | length)')

echo "Updating distribution..."
aws cloudfront update-distribution \
  --id "${DIST_ID}" \
  --if-match "${ETAG}" \
  --distribution-config "${DIST_CONFIG}" \
  --no-cli-pager > /dev/null

echo "Done. /static/* is now cached at the edge."
echo "Allow 5-10 minutes for CloudFront to deploy."
