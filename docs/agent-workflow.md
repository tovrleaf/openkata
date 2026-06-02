# Agent Workflow

## Feature Development

```mermaid
flowchart LR
    Plan["1. Plan<br/>spec-planner"] --> Implement["2. Implement<br/>frontend-developer"]
    Implement --> Validate["3. Validate<br/>spec-validator"]
    Validate -->|pass| Done["Done"]
    Validate -->|fail| Implement
```

## Agent Responsibilities

```mermaid
graph TD
    User([User]) --> DM[dojo-master]
    User -.->|direct 1-1| SP[spec-planner]
    DM --> FD[frontend-developer]
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
    DM --> IsUI{UI work?}
    DM --> IsDecision{Decision?}
    DM --> IsConvention{Convention?}
    DM --> IsValidation{Validate?}
    IsUI -->|yes| FD[frontend-developer]
    IsDecision -->|yes| ADR[adr-author]
    IsConvention -->|yes| KA[kata-author]
    IsValidation -->|yes| SV[spec-validator]
```

## Agents

| Agent | Writes to |
|-------|----------|
| spec-planner | specs/NNNN/ |
| frontend-developer | cmd/, templates, CSS |
| spec-validator | specs/NNNN/validation-report.md |
| adr-author | docs/adr/ |
| kata-author | skills/, rules/, .agents/ |
| dojo-master | nothing (coordinates only) |
