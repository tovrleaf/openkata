#!/usr/bin/env bash
# Shared by check-published.sh and push.sh: lists the release tags of this
# clone and of the remote. A release tag is either an artifact tag,
# <type>/<name>/vX.Y.Z, or a site tag, vX.Y.Z. Source this file.

readonly REMOTE="${OPENKATA_REMOTE:-origin}"
readonly TAG_RE='^((skills|rules|profiles)/[^/]+/)?v[0-9]+\.[0-9]+\.[0-9]+$'

# Prints the release tags of the local clone, sorted.
local_tags() {
  git tag -l 'skills/*/v*' 'rules/*/v*' 'profiles/*/v*' 'v[0-9]*' \
    | grep -E "${TAG_RE}" | sort || true
}

# Prints the release tags on the remote, sorted.
remote_tags() {
  git ls-remote --tags "${REMOTE}" \
    'refs/tags/skills/*' 'refs/tags/rules/*' 'refs/tags/profiles/*' \
    'refs/tags/v[0-9]*' \
    | sed -e 's/^[^[:space:]]*[[:space:]]*refs\/tags\///' \
      -e 's/\^{}$//' \
    | grep -E "${TAG_RE}" | sort -u || true
}
