#!/usr/bin/env bash
# List all rules with type, version, and change status
set -euo pipefail

dist_output=""
local_output=""

for dir in .agents/rules/*/; do
  name="$(basename "${dir}")"
  rule_dir="${dir}"
  type="local"

  if [[ -L ".agents/rules/${name}" ]]; then
    rule_dir="rules/${name}"
    type="dist"
  fi

  # Resolve version
  if [[ "${type}" == "dist" ]]; then
    # Dist: git tag only
    tag="$(git tag -l "${rule_dir}/v*" 2>/dev/null | sort -V | tail -1)"
    version="$(echo "${tag}" | grep -o 'v[0-9].*' || true)"
    if [[ -z "${version}" ]]; then
      version="unreleased"
    fi
  else
    # Local: CHANGELOG.md
    tag=""
    version="$(grep -m1 '^## ' "${rule_dir}/CHANGELOG.md" 2>/dev/null \
      | sed 's/^## //' | tr -d '[]' | sed 's/[[:space:]].*//' || true)"
    if [[ -n "${version}" ]]; then
      version="v${version}"
    fi
  fi

  changes=""
  if [[ -n "${tag}" ]]; then
    changes="$(git diff "${tag}" -- "${rule_dir}" 2>/dev/null)"
  else
    changes="$(git diff HEAD -- "${rule_dir}" 2>/dev/null)"
  fi

  status=""
  if [[ -n "${changes}" ]]; then status=" *"; fi

  line="$(printf "  \033[36m%-25s\033[0m %s%s\n" "${name}" "${version}" "${status}")"

  if [[ "${type}" == "dist" ]]; then
    dist_output+="${line}\n"
  else
    local_output+="${line}\n"
  fi
done

printf "Distributable:\n"
printf "${dist_output}"
if [[ -n "${local_output}" ]]; then
  printf "\nLocal:\n"
  printf "${local_output}"
fi
