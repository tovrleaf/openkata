#!/usr/bin/env bash
# List all specs with status and title
set -euo pipefail

for dir in specs/[0-9]*/; do
  [[ -d "${dir}" ]] || continue
  name="$(basename "${dir}")"
  spec="${dir}spec.md"

  if [[ -f "${spec}" ]]; then
    status="$(grep -m1 '^status:' "${spec}" | sed 's/status: *//')"
    title="$(grep -m1 '^# ' "${spec}" | sed 's/^# //')"
  else
    status="—"
    title="(no spec.md)"
  fi

  printf "  %s  %-12s %s\n" "${name}" "${status}" "${title}"
done
