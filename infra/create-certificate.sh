#!/usr/bin/env bash
# Request ACM certificate and add DNS validation to Route 53.
# ACM certs for CloudFront must be in us-east-1.
set -euo pipefail

ZONE_ID="${1:?Usage: $0 <hosted-zone-id>}"
DOMAIN="openkata.dev"

echo "Requesting certificate for ${DOMAIN} and *.${DOMAIN}..."
CERT_ARN=$(aws acm request-certificate \
  --domain-name "${DOMAIN}" \
  --subject-alternative-names "*.${DOMAIN}" \
  --validation-method DNS \
  --region us-east-1 \
  --query CertificateArn \
  --output text \
  --no-cli-pager)

echo "Certificate ARN: ${CERT_ARN}"
echo "Waiting for validation details..."
sleep 5

# Get DNS validation records
VALIDATION=$(aws acm describe-certificate \
  --certificate-arn "${CERT_ARN}" \
  --region us-east-1 \
  --query "Certificate.DomainValidationOptions[0].ResourceRecord" \
  --output json \
  --no-cli-pager)

CNAME_NAME=$(echo "${VALIDATION}" | grep -o '"Name": "[^"]*"' | cut -d'"' -f4)
CNAME_VALUE=$(echo "${VALIDATION}" | grep -o '"Value": "[^"]*"' | cut -d'"' -f4)

echo "Adding DNS validation record to Route 53..."
aws route53 change-resource-record-sets \
  --hosted-zone-id "${ZONE_ID}" \
  --change-batch "{
    \"Changes\": [{
      \"Action\": \"UPSERT\",
      \"ResourceRecordSet\": {
        \"Name\": \"${CNAME_NAME}\",
        \"Type\": \"CNAME\",
        \"TTL\": 300,
        \"ResourceRecords\": [{\"Value\": \"${CNAME_VALUE}\"}]
      }
    }]
  }" \
  --no-cli-pager

echo ""
echo "Validation record added. Waiting for certificate validation..."
echo "(This can take a few minutes)"
aws acm wait certificate-validated \
  --certificate-arn "${CERT_ARN}" \
  --region us-east-1

echo ""
echo "=== Certificate validated ==="
echo "ARN: ${CERT_ARN}"
echo ""
echo "Next run: ./infra/create-cloudfront.sh ${ZONE_ID} ${CERT_ARN}"
