---
status: Draft
depth: Quick
---

# Brief: Rename skill-* CSS classes to artifact-*

Scope: Rename all `.skill-*` CSS classes to `.artifact-*`
across `style.css` and all `.templ` files that reference
them. This makes the listing and detail components reusable
for rules and future artifact types. Purely mechanical
rename — no visual changes.

Files touched:
- `web/static/css/style.css` (23 occurrences)
- `cmd/openkata-web/templates/skills.templ` (15)
- `cmd/openkata-web/templates/rules.templ` (6)
- `cmd/openkata-web/templates/profiles.templ` (6)
- `cmd/openkata-web/templates/skill_detail.templ` (3)

Verify: `templ generate`, build, visual spot-check that
skills listing and detail pages look unchanged.

Date: 2026-06-02
