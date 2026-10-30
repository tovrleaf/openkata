#!/usr/bin/env bash
# Print, without running them, the commands that create the releases in
# specs/0027-site-version/history.md: a tag on each commit, then a GitHub
# release for each tag. Read them, then run them.
set -euo pipefail

readonly HISTORY="${1:-specs/0027-site-version/history.md}"
readonly MAIN="${OPENKATA_MAIN:-main}"

# Prints one "version<TAB>date<TAB>commit<TAB>summary" line for each
# table row of the history.
rows() {
  grep -E '^\| v[0-9]+\.[0-9]+\.[0-9]+ \|' "${HISTORY}" \
    | awk -F'|' '{
        for (i = 2; i <= 5; i++) gsub(/^ +| +$/, "", $i)
        gsub(/`/, "", $4)
        printf "%s\t%s\t%s\t%s\n", $2, $3, $4, $5
      }'
}

# Prints its argument in single quotes, safe to paste into a shell.
quote() {
  printf "'%s'" "${1//\'/\'\\\'\'}"
}

main() {
  local version date commit summary

  while IFS=$'\t' read -r version date commit summary; do
    if ! git merge-base --is-ancestor "${commit}" "${MAIN}"; then
      echo "error: ${commit} of ${version} is not on ${MAIN}" >&2
      exit 1
    fi
    if git rev-parse --verify --quiet "refs/tags/${version}" >/dev/null; then
      echo "error: tag ${version} already exists" >&2
      exit 1
    fi
  done < <(rows)

  echo "# Tags, on the commits of ${HISTORY}"
  while IFS=$'\t' read -r version date commit summary; do
    echo "git tag ${version} ${commit}  # ${date}"
  done < <(rows)

  echo
  echo "# Push them"
  echo "make push"

  echo
  echo "# GitHub releases, oldest first. The newest is marked latest."
  local last
  last="$(rows | tail -1 | cut -f1)"
  while IFS=$'\t' read -r version date commit summary; do
    local latest="--latest=false"
    [[ "${version}" == "${last}" ]] && latest="--latest"
    printf "gh release create %s --verify-tag %s --title %s --notes %s\n" \
      "${version}" "${latest}" "${version}" \
      "$(quote "${summary} (${date})")"
  done < <(rows)
}

main "$@"
