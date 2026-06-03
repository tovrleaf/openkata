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

| Agent | Scope | Writes to |
|-------|-------|----------|
| spec-planner | Feature planning | specs/NNNN/ |
| frontend-developer | Templates, CSS, JS | web/, templates/ |
| backend-developer | Go handlers, types, tests | cmd/, handlers, types |
| spec-validator | Acceptance review | specs/NNNN/validation-report.md |
| adr-author | Architecture decisions | docs/adr/ |
| kata-author | Skills, rules, profiles | skills/, rules/, .agents/ |
| dojo-master | Coordination | nothing (delegates only) |

## Skill & Rule Usage

### Distributable Skills (skills/)

| Skill | Used by | When |
|-------|---------|------|
| spec-workflow | spec-planner | Feature planning lifecycle |
| commit-conventions | all agents | Every commit |
| create-skill | kata-author | Creating new skills |
| create-rule | kata-author | Creating new rules |
| create-profile | kata-author | Creating new agent profiles |
| create-adr | adr-author | Recording architecture decisions |
| github-create-pr | any agent (on push) | Opening pull requests |
| git-worktree | any agent | Parallel branch work |
| makefile-conventions | backend-developer | Makefile changes |
| grill-me | spec-validator | Challenging implementations against spec |

### Local Skills (.agents/skills/)

| Skill | Used by | When |
|-------|---------|------|
| openkata-skill-conventions | kata-author | Alongside create-skill for project-specific conventions |
| openkata-rule-conventions | kata-author | Alongside create-rule for project-specific conventions |
| openkata-review-skill | kata-author | Reviewing/auditing skills |
| openkata-review-rule | kata-author | Reviewing/auditing rules |
| openkata-ryu-release | kata-author | Release workflow (diff, bump, changelog, commit, tag) |
| openkata-eval-runner | kata-author | Running eval scenarios |

### Rules (rules/)

| Rule | Used by | When |
|------|---------|------|
| git-naming | all agents | Branch names, commit messages, Assisted-by trailers |
| markdown-style | kata-author, adr-author, spec-planner | Any markdown file |
| bash-style | any agent writing scripts | Shell scripts |
| design-system | frontend-developer | All template/CSS work |