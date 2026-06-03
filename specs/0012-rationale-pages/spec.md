---
status: Draft
depth: Quick
---

# RATIONALE.md Support

Add RATIONALE.md as a convention for skills and rules.
The file explains why a skill is built the way it is —
design decisions, token economics, structural trade-offs.
Visible on the website, not distributed in archives.

Files touched:
- `cmd/openkata-web/handlers.go` (exclude from archive,
  load for detail page)
- `cmd/openkata-web/templates/skill_detail.templ` (add
  Rationale tab)
- `cmd/openkata-web/templates/rule_detail.templ` (add
  Rationale tab)
- `cmd/openkata-web/templates/profile_detail.templ` (add
  Rationale tab, after spec 0011)
- `cmd/openkata-web/templates/types.go` (add Rationale
  field to ArtifactDetail — or SkillDetail/RuleDetail
  until refactored)
- `cmd/openkata-web/handlers_test.go` (test exclusion
  and rendering)
- `skills/spec-workflow/RATIONALE.md` (first example)
- `.agents/skills/openkata-skill-conventions/SKILL.md`
  (document convention)

Behavior:
- RATIONALE.md excluded from archive downloads (same
  as CHANGELOG.md)
- RATIONALE.md rendered as a tab on detail pages (tab
  only visible if file exists, same pattern as
  Acknowledgments)
- Tab name: "Rationale"
- Tab position: after Changelog, before Acknowledgments

Verify: build passes, detail page shows Rationale tab
for spec-workflow, archive excludes the file.

Date: 2026-06-03
