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

  # Color the type
  if [[ "${type}" == "dist" ]]; then
    type_color="\033[32m${type}\033[0m"
  else
    type_color="\033[2m${type}\033[0m"
  fi

  printf "  \033[36m%-25s\033[0m ${type_color}%s\n" "${name}" "${status}"
done
