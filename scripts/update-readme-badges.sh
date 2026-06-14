#!/usr/bin/env bash
# Update README.md badges with current artifact counts.
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readme="${repo_root}/README.md"

skill_count="$(find "${repo_root}/skills" -maxdepth 2 -name "SKILL.md" | wc -l | tr -d ' ')"
rule_count="$(find "${repo_root}/rules" -maxdepth 2 -name "RULE.md" | wc -l | tr -d ' ')"
profile_count="$(find "${repo_root}/profiles" -mindepth 1 -maxdepth 1 -type d | wc -l | tr -d ' ')"
go_version="$(grep '^go ' "${repo_root}/go.mod" | awk '{print $2}')"

badges_file="$(mktemp)"
trap 'rm -f "${badges_file}"' EXIT

cat > "${badges_file}" <<EOF
<!-- badges:start -->
[![Skills](https://img.shields.io/badge/skills-${skill_count}-blue?style=flat-square)](skills/)
[![Rules](https://img.shields.io/badge/rules-${rule_count}-blue?style=flat-square)](rules/)
[![Profiles](https://img.shields.io/badge/profiles-${profile_count}-blue?style=flat-square)](profiles/)
[![MCP Server](https://img.shields.io/badge/MCP_server-ready-green?style=flat-square)](https://openkata.dev/getting-started/)
[![Build](https://github.com/tovrleaf/openkata/actions/workflows/build.yaml/badge.svg)](https://github.com/tovrleaf/openkata/actions/workflows/build.yaml)
[![Go](https://img.shields.io/badge/go-${go_version}-blue?style=flat-square)](https://go.dev/)
[![Stars](https://img.shields.io/github/stars/tovrleaf/openkata?style=flat-square)](https://github.com/tovrleaf/openkata)
[![License](https://img.shields.io/badge/license-MIT-lightgrey?style=flat-square)](LICENSE)
<!-- badges:end -->
EOF

# Replace content between markers (portable macOS + Linux)
tmp="${readme}.tmp"
awk -v bfile="${badges_file}" '
  /<!-- badges:start -->/ { skip=1; while ((getline line < bfile) > 0) print line; next }
  /<!-- badges:end -->/ { skip=0; next }
  !skip { print }
' "${readme}" > "${tmp}"
mv "${tmp}" "${readme}"
