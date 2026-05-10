#!/usr/bin/env bash
# Create Route 53 alias records pointing to CloudFront.
set -euo pipefail

ZONE_ID="${1:?Usage: $0 <hosted-zone-id> <cloudfront-domain>}"
CF_DOMAIN="${2:?Usage: $0 <hosted-zone-id> <cloudfront-domain>}"
DOMAIN="openkata.dev"

# CloudFront hosted zone ID is always Z2FDTNDATAQYW2
CF_ZONE_ID="Z2FDTNDATAQYW2"

echo "Creating DNS records for ${DOMAIN} and www.${DOMAIN}..."

aws route53 change-resource-record-sets \
  --hosted-zone-id "${ZONE_ID}" \
  --change-batch "{
    \"Changes\": [
      {
        \"Action\": \"UPSERT\",
        \"ResourceRecordSet\": {
          \"Name\": \"${DOMAIN}\",
          \"Type\": \"A\",
          \"AliasTarget\": {
            \"HostedZoneId\": \"${CF_ZONE_ID}\",
            \"DNSName\": \"${CF_DOMAIN}\",
            \"EvaluateTargetHealth\": false
          }
        }
      },
      {
        \"Action\": \"UPSERT\",
        \"ResourceRecordSet\": {
          \"Name\": \"${DOMAIN}\",
          \"Type\": \"AAAA\",
          \"AliasTarget\": {
            \"HostedZoneId\": \"${CF_ZONE_ID}\",
            \"DNSName\": \"${CF_DOMAIN}\",
            \"EvaluateTargetHealth\": false
          }
        }
      },
      {
        \"Action\": \"UPSERT\",
        \"ResourceRecordSet\": {
          \"Name\": \"www.${DOMAIN}\",
          \"Type\": \"A\",
          \"AliasTarget\": {
            \"HostedZoneId\": \"${CF_ZONE_ID}\",
            \"DNSName\": \"${CF_DOMAIN}\",
            \"EvaluateTargetHealth\": false
          }
        }
      },
      {
        \"Action\": \"UPSERT\",
        \"ResourceRecordSet\": {
          \"Name\": \"www.${DOMAIN}\",
          \"Type\": \"AAAA\",
          \"AliasTarget\": {
            \"HostedZoneId\": \"${CF_ZONE_ID}\",
            \"DNSName\": \"${CF_DOMAIN}\",
            \"EvaluateTargetHealth\": false
          }
        }
      }
    ]
  }" \
  --no-cli-pager

echo ""
echo "=== Done ==="
echo "DNS records created. Both ${DOMAIN} and www.${DOMAIN}"
echo "now point to CloudFront."
echo ""
echo "Allow a few minutes for propagation, then test:"
echo "  curl https://${DOMAIN}"
echo "  curl https://www.${DOMAIN}"
