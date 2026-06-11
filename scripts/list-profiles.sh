#!/usr/bin/env bash
# List all profiles with type, version, and change status
set -euo pipefail

dist_output=""
local_output=""

for file in .agents/profiles/*.md; do
  [[ -f "${file}" ]] || continue
  name="$(basename "${file}" .md)"
  type="local"

  if [[ -L "${file}" ]]; then
    profile_dir="profiles/${name}"
    type="dist"
  else
    profile_dir=".agents/profiles"
  fi

  # Resolve version
  if [[ "${type}" == "dist" ]]; then
    tag="$(git tag -l "profiles/${name}/v*" 2>/dev/null | sort -V | tail -1)"
    version="$(echo "${tag}" | grep -o 'v[0-9].*' || true)"
    if [[ -z "${version}" ]]; then
      version="unreleased"
    fi
  else
    version="$(grep -m1 '^## ' ".agents/profiles/CHANGELOG-${name}.md" 2>/dev/null \
      | sed 's/^## //' | tr -d '[]' | sed 's/[[:space:]].*//' || true)"
    if [[ -n "${version}" ]]; then
      version="v${version}"
    else
      version="local"
    fi
  fi

  changes=""
  if [[ "${type}" == "dist" && -n "${tag}" ]]; then
    changes="$(git diff "${tag}" -- "${profile_dir}" 2>/dev/null)"
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
