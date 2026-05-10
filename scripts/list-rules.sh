#!/usr/bin/env bash
# List all rules with type, version, and change status
set -euo pipefail

for dir in .agents/rules/*/; do
  name="$(basename "${dir}")"
  rule_dir="${dir}"
  type="local"

  if [[ -L ".agents/rules/${name}" ]]; then
    rule_dir="rules/${name}"
    type="dist"
  fi

  # Resolve version: git tag only
  tag="$(git tag -l "${rule_dir}/v*" 2>/dev/null | sort -V | tail -1)"
  version="$(echo "${tag}" | grep -o 'v[0-9].*' || true)"

  if [[ -z "${version}" ]]; then
    version="unreleased"
  fi

  changes=""
  if [[ -n "${tag}" ]]; then
    changes="$(git diff "${tag}" -- "${rule_dir}" 2>/dev/null)"
  else
    changes="$(git diff HEAD -- "${rule_dir}" 2>/dev/null)"
  fi

  status=""
  if [[ -n "${changes}" ]]; then status=" *"; fi

  # Color the type
  if [[ "${type}" == "dist" ]]; then
    type_color="\033[32m${type}\033[0m"
  else
    type_color="\033[2m${type}\033[0m"
  fi

  printf "  \033[36m%-25s\033[0m ${type_color}  %s%s\n" "${name}" "${version}" "${status}"
done
