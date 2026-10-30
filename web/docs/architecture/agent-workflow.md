---
title: Agent workflow
description: How agents move a feature from spec to release.
order: 1
---

# Agent Workflow

## Feature Development

```mermaid
flowchart LR
    Plan["1. Plan<br/>spec-planner"] --> Implement["2. Implement<br/>frontend + backend"]
    Implement --> Validate["3. Validate<br/>spec-validator"]
    Validate -->|pass| Done["Done"]
    Validate -->|fail| Implement
```

## Parallel Implementation

```mermaid
flowchart TD
    DM[dojo-master] --> FD[frontend-developer]
    DM --> BD[backend-developer]
    FD --> |templates, CSS, JS| Merge[Converge]
    BD --> |handlers, types, tests| Merge
    Merge --> Integration[Integration tasks<br/>both agents sequential]
```

## Task Parallelization Example

```text
Time ──────────────────────────────────────────────►

frontend-developer │ Nav ─── JS ─── Components ──┐
                   │                              ├─► Rules listing ─► Detail ─► Tests
backend-developer  │ Handlers ─── Types ─────────┘
```

## Agent Responsibilities

```mermaid
graph TD
    User([User]) --> DM[dojo-master]
    User -.->|direct 1-1| SP[spec-planner]
    DM --> FD[frontend-developer]
    DM --> BD[backend-developer]
    DM --> SV[spec-validator]
    DM --> AA[adr-author]
    DM --> KA[kata-author]
```

## Routing

```mermaid
flowchart TD
    Request([Request]) --> IsPlanning{Planning?}
    IsPlanning -->|yes| Direct["Switch to spec-planner<br/>kiro-cli --agent spec-planner"]
    IsPlanning -->|no| DM[dojo-master]
    DM --> IsUI{Templates,<br/>CSS, JS?}
    DM --> IsGo{Handlers,<br/>types, tests?}
    DM --> IsDecision{Decision?}
    DM --> IsConvention{Convention?}
    DM --> IsValidation{Validate?}
    IsUI -->|yes| FD[frontend-developer]
    IsGo -->|yes| BD[backend-developer]
    IsDecision -->|yes| ADR[adr-author]
    IsConvention -->|yes| KA[kata-author]
    IsValidation -->|yes| SV[spec-validator]
```

## Agents

The scopes come from each agent's profile in `.agents/profiles/`.

- **spec-planner** plans features. It writes to `specs/`.
- **frontend-developer** works on templates, CSS and JS. It writes to
  `cmd/openkata-web/templates/` and `web/static/`.
- **backend-developer** works on Go handlers, types and tests. It writes
  to `handlers.go`, `handlers_test.go` and `main.go` in
  `cmd/openkata-web/`, to `templates/types.go` and
  `templates/helpers.go`, and to `cmd/openkata-mcp/`.
- **spec-validator** reviews an implementation against its spec. It
  writes to `specs/NNNN-slug/validation-report.md`.
- **adr-author** records architecture decisions. It writes to
  `docs/adr/`.
- **kata-author** writes skills, rules and profiles. It writes to
  `skills/`, `rules/`, `profiles/`, `.agents/skills/` and
  `.agents/rules/`.
- **dojo-master** coordinates the others. It writes nothing, and
  delegates all work.

## What each agent uses

Each profile names the skills and rules its agent works with. A skill
or rule that no profile names is still available to any agent.

- **spec-planner**
  - Skills: `spec-workflow` (phases 1-3), and optionally `grill-me`,
    `grill-with-docs` and `critical-thinking`
  - Rules: markdown-style, git-naming
- **spec-validator**
  - Skills: `spec-workflow` (phase 5)
  - Rules: markdown-style
- **adr-author**
  - Skills: `create-adr`
  - Rules: markdown-style, git-naming
- **kata-author**
  - Skills: `create-skill` with `openkata-skill-conventions`,
    `create-rule` with `openkata-rule-conventions`, `create-profile`,
    `openkata-review-skill`, `openkata-review-rule`,
    `openkata-ryu-release` and `openkata-eval-runner`
  - Rules: markdown-style, git-naming
- **frontend-developer**
  - Skills: none named
  - Rules: design-system, git-naming
- **backend-developer**
  - Skills: none named
  - Rules: git-naming
- **dojo-master**
  - Skills: none, because it delegates
  - Rules: git-naming, for the `Assisted-by` trailer on delegated
    commits

The distributable skills, rules and profiles are in `skills/`,
`rules/` and `profiles/`. The local skills, which only this
repository uses, are in `.agents/skills/`.
