# Tasks: Rename skill-* CSS classes to artifact-*

## Tasks

### 1. Rename CSS classes in style.css
- **Status**: Pending
- **Goal**: Rename all `.skill-*` class definitions to
  `.artifact-*` in the stylesheet
- **Boundary**: `web/static/css/style.css`
- **Depends**: None
- **Done when**: No `.skill-*` classes remain in CSS,
  all are `.artifact-*`
- **Verify**: `grep -c 'skill-' web/static/css/style.css`
  returns 0

### 2. Update template references
- **Status**: Pending
- **Goal**: Update all `.templ` files to use `.artifact-*`
  class names
- **Boundary**: `cmd/openkata-web/templates/skills.templ`,
  `rules.templ`, `profiles.templ`, `skill_detail.templ`
- **Depends**: 1
- **Done when**: No `skill-` class references in `.templ`
  files, `templ generate` succeeds
- **Verify**: `templ generate ./cmd/openkata-web/templates/ && go build -o bin/openkata-web ./cmd/openkata-web/`

### 3. Visual verification
- **Status**: Pending
- **Goal**: Confirm skills listing and detail pages render
  identically after rename
- **Boundary**: None (verification only)
- **Depends**: 2
- **Done when**: Skills listing page and skill detail page
  look unchanged, no broken styles
- **Verify**: Run dev server, check `/skills/` and
  `/skills/commit-conventions`

## Progress Log
