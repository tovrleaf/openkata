#!/usr/bin/env bash
# Push the current branch and every release tag that exists only in this
# clone, one push for each, so a tag is never forgotten. It shows what it
# will push and asks first. The answer defaults to no. With --yes it does
# not ask, for an agent that has been told to push.
set -euo pipefail

# shellcheck source=scripts/tag_lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/tag_lib.sh"

main() {
  local branch commits tags answer tag assume_yes=0

  if [[ "${1:-}" == "--yes" || "${1:-}" == "-y" ]]; then
    assume_yes=1
  fi

  branch="$(git rev-parse --abbrev-ref HEAD)"
  if [[ "${branch}" == "HEAD" ]]; then
    echo "error: not on a branch" >&2
    exit 1
  fi

  if git rev-parse --verify --quiet "${branch}@{upstream}" >/dev/null; then
    commits="$(git log --oneline "${branch}@{upstream}..HEAD")"
  else
    echo "Branch ${branch} is new on ${REMOTE}, compared with main."
    commits="$(git log --oneline "${REMOTE}/main..HEAD")"
  fi
  tags="$(comm -23 <(local_tags) <(remote_tags))"

  if [[ -z "${commits}" && -z "${tags}" ]]; then
    echo "Nothing to push."
    exit 0
  fi

  if [[ -n "${commits}" ]]; then
    echo "Commits to push to ${REMOTE}/${branch}:"
    echo "${commits}" | sed 's/^/  /'
  fi
  if [[ -n "${tags}" ]]; then
    echo "Tags to push:"
    echo "${tags}" | sed 's/^/  /'
  fi

  if (( assume_yes )); then
    answer="y"
  else
    read -r -p "Push? [y/N] " answer || answer=""
  fi
  if [[ "${answer}" != "y" && "${answer}" != "Y" ]]; then
    echo "Nothing was pushed."
    exit 0
  fi

  if [[ -n "${commits}" ]]; then
    git push "${REMOTE}" "${branch}"
  fi
  while read -r tag; do
    [[ -n "${tag}" ]] || continue
    git push "${REMOTE}" "refs/tags/${tag}"
  done <<< "${tags}"

  "$(dirname "${BASH_SOURCE[0]}")/check-published.sh" || true
}

main "$@"
