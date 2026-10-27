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

| Agent | Scope | Writes to |
|-------|-------|-----------|
| spec-planner | Feature planning | `specs/` |
| frontend-developer | Templates, CSS, JS | `cmd/openkata-web/templates/`, `web/static/` |
| backend-developer | Go handlers, types, tests | `cmd/openkata-web/` (`handlers.go`, `handlers_test.go`, `main.go`), `templates/types.go`, `templates/helpers.go`, `cmd/openkata-mcp/` |
| spec-validator | Acceptance review | `specs/NNNN-slug/validation-report.md` |
| adr-author | Architecture decisions | `docs/adr/` |
| kata-author | Skills, rules, profiles | `skills/`, `rules/`, `profiles/`, `.agents/skills/`, `.agents/rules/` |
| dojo-master | Coordination | nothing (delegates only) |

## What each agent uses

Each profile names the skills and rules its agent works with. A
skill or rule that no profile names is still available to any agent.

| Agent | Skills | Rules |
|-------|--------|-------|
| spec-planner | `spec-workflow` (phases 1-3); optional `grill-me`, `grill-with-docs`, `critical-thinking` | markdown-style, git-naming |
| spec-validator | `spec-workflow` (phase 5) | markdown-style |
| adr-author | `create-adr` | markdown-style, git-naming |
| kata-author | `create-skill` with `openkata-skill-conventions`; `create-rule` with `openkata-rule-conventions`; `create-profile`; `openkata-review-skill`; `openkata-review-rule`; `openkata-ryu-release`; `openkata-eval-runner` | markdown-style, git-naming |
| frontend-developer | none named | design-system, git-naming |
| backend-developer | none named | git-naming |
| dojo-master | none; it delegates | git-naming, for the `Assisted-by` trailer on delegated commits |

The distributable skills, rules and profiles are in `skills/`,
`rules/` and `profiles/`. The local skills, which only this
repository uses, are in `.agents/skills/`.
