#!/usr/bin/env bash
# Configure branch protection for main using gh CLI.
set -euo pipefail

REPO="tovrleaf/openkata"

echo "Enabling branch protection on main..."
gh api \
  --method PUT \
  "repos/${REPO}/branches/main/protection" \
  --input - <<'EOF'
{
  "required_pull_request_reviews": null,
  "enforce_admins": true,
  "required_status_checks": null,
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false
}
EOF

echo ""
echo "=== Done ==="
echo "Main branch is now protected:"
echo "  - Direct pushes blocked"
echo "  - Force pushes blocked"
echo "  - No review required (solo project)"
