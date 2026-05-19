#!/usr/bin/env bash
set -euo pipefail

echo "Agent Profiles"
echo "=============="
echo ""

echo "Distributable:"
for file in profiles/*.md; do
  [[ -f "${file}" ]] || continue
  name="$(basename "${file}" .md)"
  desc="$(sed -n '/^[^#]/{ /^$/d; p; q; }' "${file}")"
  printf "  \033[36m%-20s\033[0m %s\n" "${name}" "${desc}"
done

echo ""
echo "Local:"
for file in .agents/profiles/*.md; do
  [[ -f "${file}" ]] || continue
  [[ -L "${file}" ]] && continue
  name="$(basename "${file}" .md)"
  desc="$(sed -n '/^[^#]/{ /^$/d; p; q; }' "${file}")"
  printf "  \033[36m%-20s\033[0m %s\n" "${name}" "${desc}"
done
