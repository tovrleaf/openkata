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
  elif [[ -f "${dir}brief.md" ]]; then
    status="$(grep -m1 '^status:' "${dir}brief.md" | sed 's/status: *//')"
    title="$(grep -m1 '^# ' "${dir}brief.md" | sed 's/^# //')"
  else
    status="—"
    title="(no spec.md)"
  fi

  # Color the status
  case "${status}" in
    Done) status_color="\033[32m%-12s\033[0m" ;;
    "In Progress") status_color="\033[33m%-12s\033[0m" ;;
    *) status_color="\033[33m%-12s\033[0m" ;;
  esac

  printf "  \033[36m%s\033[0m  ${status_color} %s\n" "${name}" "${status}" "${title}"
done
