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
  "required_pull_request_reviews": {
    "required_approving_review_count": 1,
    "dismiss_stale_reviews": true
  },
  "enforce_admins": true,
  "required_status_checks": null,
  "restrictions": null
}
EOF

echo ""
echo "=== Done ==="
echo "Main branch is now protected:"
echo "  - Requires PR with 1 approval"
echo "  - Stale reviews dismissed on new commits"
echo "  - Admins cannot bypass"
