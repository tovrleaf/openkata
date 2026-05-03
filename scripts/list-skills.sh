#!/usr/bin/env bash
# List all skills with type, version, and change status
set -euo pipefail

for dir in .agents/skills/*/; do
  name="$(basename "${dir}")"
  skill_dir="${dir}"
  type="local"

  if [[ -L ".agents/skills/${name}" ]]; then
    skill_dir="skills/${name}"
    type="dist"
  fi

  # Resolve version: git tag first, then CHANGELOG.md
  tag="$(git tag -l "${skill_dir}/v*" 2>/dev/null | sort -V | tail -1)"
  version="$(echo "${tag}" | grep -o 'v[0-9].*' || true)"

  if [[ -z "${version}" ]]; then
    version="$(grep -m1 '^## ' "${skill_dir}/CHANGELOG.md" 2>/dev/null \
      | sed 's/^## //' | tr -d '[]' | sed 's/[[:space:]].*//' || echo "?")"
    version="v${version}"
  fi

  changes=""
  if [[ -n "${tag}" ]]; then
    changes="$(git diff "${tag}" -- "${skill_dir}" 2>/dev/null)"
  else
    changes="$(git diff HEAD -- "${skill_dir}" 2>/dev/null)"
  fi

  status=""
  if [[ -n "${changes}" ]]; then status=" *"; fi

  printf "  %-25s %-7s %s%s\n" "${name}" "${type}" "${version}" "${status}"
done
