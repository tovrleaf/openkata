#!/usr/bin/env bash
# Check development prerequisites
set -euo pipefail

ok=0
fail=0

check() {
  local name="${1}"
  local cmd="${2}"
  local version_flag="${3:-}"
  local install_hint="${4:-}"

  if command -v "${cmd}" &>/dev/null; then
    local version=""
    if [[ -n "${version_flag}" ]]; then
      version=" ($(${cmd} ${version_flag} 2>&1 | head -1))"
    fi
    printf "  ✔ %-12s found%s\n" "${name}" "${version}"
    ok=$(( ok + 1 ))
  else
    printf "  ✘ %-12s missing"  "${name}"
    if [[ -n "${install_hint}" ]]; then
      printf " — %s" "${install_hint}"
    fi
    printf "\n"
    fail=$(( fail + 1 ))
  fi
}

echo "Checking prerequisites..."
echo ""

check "git"   "git"   "--version"
check "go"    "go"    "version"    "https://go.dev/dl/"
check "make"  "make"  "--version"
check "tessl" "tessl" "--version"  "https://docs.tessl.io/introduction-to-tessl/installation"

echo ""
if (( fail > 0 )); then
  echo "${ok} found, ${fail} missing"
  exit 1
else
  echo "All ${ok} prerequisites found"
fi
