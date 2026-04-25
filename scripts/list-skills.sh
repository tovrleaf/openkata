#!/usr/bin/env bash
# List all skills with type, version, and change status

for dir in .agents/skills/*/; do
  name=$(basename "$dir")
  skill_dir="$dir"
  type="local"

  if [ -L ".agents/skills/$name" ]; then
    skill_dir="skills/$name"
    type="dist"
  fi

  version=$(awk '/^version:/{print $2; exit}' "$skill_dir/SKILL.md" 2>/dev/null || echo "?")
  tag="$skill_dir/v$version"

  changes=""
  if [ "$type" = "dist" ] && git tag -l "$tag" | grep -q .; then
    changes=$(git diff "$tag" -- "$skill_dir" 2>/dev/null)
  else
    changes=$(git diff HEAD -- "$skill_dir" 2>/dev/null)
  fi

  status=""
  if [ -n "$changes" ]; then status=" *"; fi

  printf "  %-25s %-7s v%s%s\n" "$name" "$type" "$version" "$status"
done
