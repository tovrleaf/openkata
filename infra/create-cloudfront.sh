#!/usr/bin/env bash
# Create CloudFront distribution with Lambda Function URL origin.
set -euo pipefail

ZONE_ID="${1:?Usage: $0 <hosted-zone-id> <cert-arn>}"
CERT_ARN="${2:?Usage: $0 <hosted-zone-id> <cert-arn>}"
DOMAIN="openkata.dev"
REGION="${AWS_REGION:-eu-north-1}"

# Get the Lambda Function URL origin
FUNCTION_URL=$(aws lambda get-function-url-config \
  --function-name openkata-web \
  --region "${REGION}" \
  --query FunctionUrl \
  --output text \
  --no-cli-pager)

# Strip https:// and trailing slash for origin domain
ORIGIN_DOMAIN=$(echo "${FUNCTION_URL}" | sed 's|https://||;s|/$||')

echo "Creating CloudFront distribution..."
echo "Origin: ${ORIGIN_DOMAIN}"

DIST_RESULT=$(aws cloudfront create-distribution \
  --distribution-config "{
    \"CallerReference\": \"openkata-$(date +%s)\",
    \"Comment\": \"Open Kata website\",
    \"Enabled\": true,
    \"Aliases\": {
      \"Quantity\": 2,
      \"Items\": [\"${DOMAIN}\", \"www.${DOMAIN}\"]
    },
    \"Origins\": {
      \"Quantity\": 1,
      \"Items\": [{
        \"Id\": \"lambda\",
        \"DomainName\": \"${ORIGIN_DOMAIN}\",
        \"CustomOriginConfig\": {
          \"HTTPPort\": 80,
          \"HTTPSPort\": 443,
          \"OriginProtocolPolicy\": \"https-only\"
        }
      }]
    },
    \"DefaultCacheBehavior\": {
      \"TargetOriginId\": \"lambda\",
      \"ViewerProtocolPolicy\": \"redirect-to-https\",
      \"AllowedMethods\": {
        \"Quantity\": 7,
        \"Items\": [\"GET\", \"HEAD\", \"OPTIONS\", \"PUT\", \"POST\", \"PATCH\", \"DELETE\"]
      },
      \"CachePolicyId\": \"4135ea2d-6df8-44a3-9df3-4b5a84be39ad\",
      \"Compress\": true
    },
    \"ViewerCertificate\": {
      \"ACMCertificateArn\": \"${CERT_ARN}\",
      \"SSLSupportMethod\": \"sni-only\",
      \"MinimumProtocolVersion\": \"TLSv1.2_2021\"
    },
    \"HttpVersion\": \"http2and3\",
    \"PriceClass\": \"PriceClass_100\"
  }" \
  --output json \
  --no-cli-pager)

DIST_ID=$(echo "${DIST_RESULT}" | grep '"Id"' | head -1 | cut -d'"' -f4)
CF_DOMAIN=$(echo "${DIST_RESULT}" | grep '"DomainName"' | head -1 | cut -d'"' -f4)

echo ""
echo "=== CloudFront distribution created ==="
echo "Distribution ID: ${DIST_ID}"
echo "CloudFront domain: ${CF_DOMAIN}"
echo ""
echo "Next run: ./infra/create-dns-records.sh ${ZONE_ID} ${CF_DOMAIN}"
