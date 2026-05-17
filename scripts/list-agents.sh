#!/usr/bin/env bash
set -euo pipefail

echo "Kiro Agents"
echo "==========="
echo ""

for file in .kiro/agents/*.json; do
  [[ -f "${file}" ]] || continue
  name="$(basename "${file}" .json)"
  desc="$(grep -o '"description": "[^"]*"' "${file}" | cut -d'"' -f4)"
  printf "  \033[36m%-16s\033[0m %s\n" "${name}" "${desc}"
done
