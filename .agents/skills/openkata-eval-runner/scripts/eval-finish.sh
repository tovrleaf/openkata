#!/usr/bin/env bash
set -euo pipefail

name="${1:?Usage: eval-finish.sh <skill-name>}"
[[ -d "skills/${name}" ]] || { echo "Error: skills/${name} not found"; exit 1; }

echo "Running tessl skill review --optimize..."
tessl skill review --optimize "skills/${name}"
