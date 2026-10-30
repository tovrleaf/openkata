#!/usr/bin/env bash
# Check that release tags were pushed and that the pushed artifact tags
# are on the site. Exits 1 if a tag was never pushed or a version is
# missing, and 2 if the site could not be reached.
set -euo pipefail

readonly SITE="${OPENKATA_SITE:-https://openkata.dev}"

# shellcheck source=scripts/tag_lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/tag_lib.sh"

# Prints the URL of a file of the tagged version. The site answers 200
# only for a file that is in the bucket, so it is the check.
probe_url() {
  local tag="$1" type name version
  type="${tag%%/*}"
  name="${tag#*/}"
  name="${name%%/*}"
  version="${tag##*/v}"
  case "${type}" in
    skills) echo "${SITE}/skills/${name}/${version}/raw/SKILL.md" ;;
    rules) echo "${SITE}/rules/${name}/${version}/raw/RULE.md" ;;
    profiles) echo "${SITE}/profiles/${name}/${version}/raw" ;;
  esac
}

main() {
  local unpushed missing=() tag code failed=0 renamed=0

  unpushed="$(comm -23 <(local_tags) <(remote_tags))"
  if [[ -n "${unpushed}" ]]; then
    echo "Tags that exist only in this clone (never pushed):"
    echo "${unpushed}" | sed 's/^/  /'
    failed=1
  fi

  # Site tags (vX.Y.Z) have no files in the bucket to look for.
  while read -r tag; do
    [[ "${tag}" == */* ]] || continue
    if ! code="$(curl -s -o /dev/null -w '%{http_code}' \
      "$(probe_url "${tag}")")"; then
      echo "error: cannot reach ${SITE}" >&2
      exit 2
    fi
    if [[ "${code}" == "404" ]]; then
      missing+=("${tag}")
    elif [[ "${code}" == "301" ]]; then
      # A renamed artifact: its old files were deleted on purpose.
      renamed=$((renamed + 1))
    elif [[ "${code}" != "200" ]]; then
      echo "error: ${tag} answered ${code}" >&2
      exit 2
    fi
  done < <(remote_tags)

  if (( ${#missing[@]} > 0 )); then
    echo "Pushed tags whose files are not on ${SITE}:"
    printf '  %s\n' "${missing[@]}"
    echo "Publish one with: gh workflow run publish.yaml -f tag=<tag>"
    failed=1
  fi
  if (( renamed > 0 )); then
    echo "Skipped ${renamed} tags of renamed artifacts (redirected)."
  fi
  if (( failed == 0 )); then
    echo "All release tags are pushed and published."
  fi
  exit "${failed}"
}

main "$@"
