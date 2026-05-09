#!/usr/bin/env bash
# List all ADRs with number, status, and title
set -euo pipefail

for file in docs/adr/[0-9]*.md; do
  [[ -f "${file}" ]] || continue
  number="$(basename "${file}" .md | cut -d- -f1)"
  status="$(grep -m1 '^status:' "${file}" | sed 's/status: *//')"
  title="$(grep -m1 '^# ' "${file}" | sed 's/^# [0-9]*\. //')"
  # Color the status
  if [[ "${status}" == "ACCEPTED" ]]; then
    status_color="\033[32m%-10s\033[0m"
  else
    status_color="\033[33m%-10s\033[0m"
  fi

  printf "  \033[36m%s\033[0m  ${status_color} %s\n" "${number}" "${status}" "${title}"
done
