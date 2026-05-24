#!/usr/bin/env bash
# List all skills with type, version, and change status
set -euo pipefail

subcmd="${1:-help}"

if [[ "${subcmd}" == "help" ]]; then
  echo "Skills"
  echo "======"
  echo ""
  printf "  \033[36m%-12s\033[0m %s\n" "list" "List all skills with type and version"
  printf "  \033[36m%-12s\033[0m %s\n" "status" "List skills with registry status (network)"
  exit 0
fi

if [[ "${subcmd}" == "status" ]]; then
  TESSL_CHECK=1
fi

dist_output=""
local_output=""

for dir in .agents/skills/*/; do
  name="$(basename "${dir}")"
  skill_dir="${dir}"
  type="local"

  if [[ -L ".agents/skills/${name}" ]]; then
    skill_dir="skills/${name}"
    type="dist"
  fi

  # Resolve version
  if [[ "${type}" == "dist" ]]; then
    # Dist: git tag only
    tag="$(git tag -l "${skill_dir}/v*" 2>/dev/null | sort -V | tail -1)"
    version="$(echo "${tag}" | grep -o 'v[0-9].*' || true)"
    if [[ -z "${version}" ]]; then
      version="unreleased"
    fi
  else
    # Local: CHANGELOG.md
    tag=""
    version="$(grep -m1 '^## ' "${skill_dir}/CHANGELOG.md" 2>/dev/null \
      | sed 's/^## //' | tr -d '[]' | sed 's/[[:space:]].*//' || true)"
    if [[ -n "${version}" ]]; then
      version="v${version}"
    fi
  fi

  changes=""
  if [[ -n "${tag}" ]]; then
    changes="$(git diff "${tag}" -- "${skill_dir}" 2>/dev/null)"
  else
    changes="$(git diff HEAD -- "${skill_dir}" 2>/dev/null)"
  fi

  if [[ -n "${changes}" ]]; then
    version="${version}*"
  fi

  # Eval status for distributable skills
  eval_indicator=""
  registry_indicator=""
  if [[ "${type}" == "dist" ]]; then
    # Eval indicator: ✓ = has evals and current, ✓* = stale, × = no evals
    if [[ -d "${skill_dir}/evals" ]]; then
      skill_ts="$(git log -1 --format=%ct -- "${skill_dir}" ':!'"${skill_dir}"'/evals' 2>/dev/null || echo 0)"
      eval_ts="$(git log -1 --format=%ct -- "${skill_dir}/evals" 2>/dev/null || echo 0)"
      if [[ -z "${eval_ts}" || "${eval_ts}" == "0" || "${skill_ts}" -gt "${eval_ts}" ]]; then
        eval_indicator="\033[33m✓*\033[0m"
      else
        eval_indicator="\033[32m✓\033[0m"
      fi
    else
      eval_indicator="\033[31m×\033[0m"
    fi

    # Registry indicator (only with TESSL_CHECK=1)
    if [[ "${TESSL_CHECK:-}" == "1" ]]; then
      tile_name="$(grep -o '"name": "[^"]*"' "${skill_dir}/tile.json" 2>/dev/null \
        | cut -d'"' -f4 || true)"
      if [[ -n "${tile_name}" ]]; then
        reg_ver="$(tessl tile info "${tile_name}" 2>/dev/null \
          | sed 's/\x1b\[[0-9;]*m//g' \
          | grep "Latest Version" | sed 's/.*Latest Version[[:space:]]*//' || true)"
        if [[ -n "${reg_ver}" ]]; then
          # Dirty if source changed since tag
          if [[ -n "${changes}" ]]; then
            registry_indicator="${reg_ver}*"
          else
            registry_indicator="${reg_ver}"
          fi
        else
          registry_indicator="\033[31m×\033[0m         "
        fi
      else
        registry_indicator="\033[31m×\033[0m         "
      fi
    fi
  fi

  # Format: name  local_version  registry_version  evals
  if [[ "${TESSL_CHECK:-}" == "1" && "${type}" == "dist" ]]; then
    line="$(printf "  \033[36m%-25s\033[0m %-14s %-10b %b\n" "${name}" "${version}" "${registry_indicator}" "${eval_indicator}")"
  elif [[ "${type}" == "dist" ]]; then
    line="$(printf "  \033[36m%-25s\033[0m %-14s %b\n" "${name}" "${version}" "${eval_indicator}")"
  else
    line="$(printf "  \033[36m%-25s\033[0m %s\n" "${name}" "${version}")"
  fi

  if [[ "${type}" == "dist" ]]; then
    dist_output+="${line}\n"
  else
    local_output+="${line}\n"
  fi
done

printf "Distributable:\n"
if [[ "${TESSL_CHECK:-}" == "1" ]]; then
  printf "  \033[2m%-25s %-14s %-10s %s\033[0m\n" "" "local" "registry" "evals"
else
  printf "  \033[2m%-25s %-14s %s\033[0m\n" "" "local" "evals"
fi
printf "${dist_output}"
printf "\nLocal:\n"
printf "${local_output}"
