#!/usr/bin/env bash
# Create Route 53 hosted zone for openkata.dev
# After running, update nameservers at your registrar.
set -euo pipefail

DOMAIN="openkata.dev"
REGION="${AWS_REGION:-eu-north-1}"

echo "Creating hosted zone for ${DOMAIN}..."
RESULT=$(aws route53 create-hosted-zone \
  --name "${DOMAIN}" \
  --caller-reference "openkata-$(date +%s)" \
  --no-cli-pager \
  --output json)

ZONE_ID=$(echo "${RESULT}" | grep -o '"Id": "/hostedzone/[^"]*"' | cut -d'/' -f3 | tr -d '"')
echo ""
echo "Hosted zone created: ${ZONE_ID}"
echo ""
echo "=== ACTION REQUIRED ==="
echo "Set these nameservers at your domain registrar:"
echo ""
aws route53 get-hosted-zone \
  --id "${ZONE_ID}" \
  --query "DelegationSet.NameServers" \
  --output text \
  --no-cli-pager | tr '\t' '\n'
echo ""
echo "After updating nameservers, wait for propagation"
echo "(can take up to 48h, usually minutes)."
echo ""
echo "Then run: ./infra/create-certificate.sh ${ZONE_ID}"
