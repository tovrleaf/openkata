#!/usr/bin/env bash
set -euo pipefail

name="${1:?Usage: eval-start.sh <skill-name>}"
[[ -d "skills/${name}" ]] || { echo "Error: skills/${name} not found"; exit 1; }
[[ -d "skills/${name}/evals" ]] || { echo "Error: no evals/ directory in skills/${name}"; exit 1; }

make eval-local "skills/${name}"
