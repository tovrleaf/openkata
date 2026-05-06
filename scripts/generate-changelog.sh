#!/usr/bin/env bash
# Generate root CHANGELOG.md from individual artifact changelogs
set -euo pipefail

output="CHANGELOG.md"

cat > "${output}" << 'HEADER'
# Changelog

All releases across Open Kata skills and rules, newest first.
HEADER

# Collect all releases: date|name|type|version|summary
entries=""
for f in skills/*/CHANGELOG.md rules/*/CHANGELOG.md; do
  [[ -f "${f}" ]] || continue
  name="$(echo "${f}" | cut -d/ -f2)"
  type="$(echo "${f}" | cut -d/ -f1)"

  # Parse each version heading and grab first bullet as summary
  awk -v name="${name}" -v type="${type}" '
    /^## / {
      if (version) print date "|" name "|" type "|" version "|" summary
      line = $0
      gsub(/^## /, "", line)
      gsub(/\[|\]/, "", line)
      split(line, parts, " ")
      version = parts[1]
      # Find date (YYYY-MM-DD pattern)
      match(line, /[0-9]{4}-[0-9]{2}-[0-9]{2}/)
      date = substr(line, RSTART, RLENGTH)
      summary = ""
      next
    }
    /^- / && summary == "" {
      s = $0
      gsub(/^- /, "", s)
      # Trim to ~60 chars
      if (length(s) > 60) s = substr(s, 1, 57) "..."
      summary = s
    }
    END { if (version) print date "|" name "|" type "|" version "|" summary }
  ' "${f}"
done | sort -t'|' -k1,1r -k2,2 > /tmp/changelog_entries.txt

# Group by date and write
current_date=""
while IFS='|' read -r date name type version summary; do
  if [[ "${date}" != "${current_date}" ]]; then
    current_date="${date}"
    printf "\n## %s\n\n" "${date}" >> "${output}"
  fi
  printf -- "- **[%s](%s/%s/CHANGELOG.md)** v%s — %s\n" \
    "${name}" "${type}" "${name}" "${version}" "${summary}" >> "${output}"
done < /tmp/changelog_entries.txt

rm -f /tmp/changelog_entries.txt
echo "Generated ${output}"
