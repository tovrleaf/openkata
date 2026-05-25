#!/usr/bin/env bash
set -euo pipefail

name="${1:?Usage: eval-finish.sh <skill-name>}"
[[ -d "skills/${name}" ]] || { echo "Error: skills/${name} not found"; exit 1; }
worktree=".worktrees/eval-${name}"
branch="eval/${name}"

if [[ ! -d "${worktree}" ]]; then
  echo "No worktree at ${worktree}"
  exit 1
fi

if [[ -z "$(git -C "${worktree}" status --porcelain)" ]]; then
  echo "No uncommitted changes in ${worktree}."
else
  echo "Uncommitted changes in ${worktree}:"
  git -C "${worktree}" status --short
  echo "Commit them first, then re-run."
  exit 1
fi

git merge --ff-only "${branch}"
git worktree remove "${worktree}"
git branch -d "${branch}"
echo "Merged and cleaned up ${branch}."
