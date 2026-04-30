#!/usr/bin/env bash
# List all rules with type and change status
set -euo pipefail

for dir in .agents/rules/*/; do
  name="$(basename "${dir}")"
  rule_dir="${dir}"
  type="local"

  if [[ -L ".agents/rules/${name}" ]]; then
    rule_dir="rules/${name}"
    type="dist"
  fi

  changes="$(git diff HEAD -- "${rule_dir}" 2>/dev/null)"
  status=""
  if [[ -n "${changes}" ]]; then status=" *"; fi

  printf "  %-25s %-7s%s\n" "${name}" "${type}" "${status}"
done
