#!/usr/bin/env bash
set -euo pipefail

name="${1:?Usage: eval-start.sh <skill-name>}"
[[ -d "skills/${name}" ]] || { echo "Error: skills/${name} not found"; exit 1; }

git worktree add ".worktrees/eval-${name}" -b "eval/${name}" 2>/dev/null \
  || echo "Worktree already exists, reusing."
cd ".worktrees/eval-${name}"
tessl scenario generate "skills/${name}"
tessl eval run "skills/${name}/" --variant with-context
