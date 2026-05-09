#!/usr/bin/env bash
# List all ADRs with number, status, and title
set -euo pipefail

for file in docs/adr/[0-9]*.md; do
  [[ -f "${file}" ]] || continue
  number="$(basename "${file}" .md | cut -d- -f1)"
  status="$(grep -m1 '^status:' "${file}" | sed 's/status: *//')"
  title="$(grep -m1 '^# ' "${file}" | sed 's/^# [0-9]*\. //')"
  printf "  \033[36m%s\033[0m  %-10s %s\n" "${number}" "${status}" "${title}"
done
