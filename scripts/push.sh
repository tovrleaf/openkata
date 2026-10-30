#!/usr/bin/env bash
# Push the current branch and every release tag that exists only in this
# clone, one push for each, so a tag is never forgotten. It shows what it
# will push and asks first. The answer defaults to no. With --yes it does
# not ask, for an agent that has been told to push.
set -euo pipefail

# shellcheck source=scripts/tag_lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/tag_lib.sh"

main() {
  local branch base commits tags answer tag assume_yes=0

  if [[ "${1:-}" == "--yes" || "${1:-}" == "-y" ]]; then
    assume_yes=1
  fi

  branch="$(git rev-parse --abbrev-ref HEAD)"
  if [[ "${branch}" == "HEAD" ]]; then
    echo "error: not on a branch" >&2
    exit 1
  fi

  # Compare with the branch on the remote, whether or not this clone has it
  # set as its upstream, and with main when the branch is new there.
  if git rev-parse --verify --quiet \
    "refs/remotes/${REMOTE}/${branch}" >/dev/null; then
    base="${REMOTE}/${branch}"
  else
    echo "Branch ${branch} is new on ${REMOTE}, compared with main."
    base="${REMOTE}/main"
  fi
  commits="$(git log --oneline "${base}..HEAD")"
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

  # Runs after the push, so a failure here is reported but the push stands
  "$(dirname "${BASH_SOURCE[0]}")/check-published.sh"
}

main "$@"
