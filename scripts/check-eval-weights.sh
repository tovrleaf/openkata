#!/usr/bin/env bash
set -euo pipefail

# Validates that criteria.json max_score values sum to 100
# for each eval scenario.
#
# Usage:
#   ./scripts/check-eval-weights.sh skills/<name>
#   ./scripts/check-eval-weights.sh  (checks all skills)

failed=0

if [[ $# -ge 1 ]]; then
  base="${1}"
else
  base="skills"
fi

while IFS= read -r file; do
  sum=$(grep '"max_score":' "${file}" \
    | sed 's/.*"max_score": //' \
    | awk '{s+=$1} END {print s}')

  if [[ "${sum}" != "100" ]]; then
    scenario=$(dirname "${file}")
    echo "FAIL: ${scenario} — sum=${sum} (expected 100)"
    failed=1
  fi
done < <(find "${base}" -path '*/evals/*/criteria.json' | sort)

if [[ "${failed}" -eq 0 ]]; then
  echo "All eval criteria weights sum to 100."
fi

exit "${failed}"
