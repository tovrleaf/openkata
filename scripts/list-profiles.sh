#!/usr/bin/env bash
set -euo pipefail

echo "Agent Profiles"
echo "=============="
echo ""
printf "  \033[36m%-16s\033[0m %s\n" "NAME" "DESCRIPTION"
echo ""

for file in profiles/*.md; do
  [[ -f "${file}" ]] || continue
  name="$(basename "${file}" .md)"
  # Extract first non-heading, non-empty line as description
  desc="$(sed -n '/^[^#]/{ /^$/d; p; q; }' "${file}")"
  printf "  \033[36m%-16s\033[0m %s\n" "${name}" "${desc}"
done
