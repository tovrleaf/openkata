# Open Kata

Codified practices for AI agents — teach your agents the way.

See the [Manifesto](MANIFESTO.md) for why this exists.

## Vocabulary

| Concept | Kata term | Description |
|---------|-----------|-------------|
| Skills | **kata** | Codified practices agents follow |
| Rules | **dojo kun** | Always-on constraints applied to every session |
| Roles | **sensei profiles** | Agent role definitions with scoped permissions |
| Prompts | **kata forms** | Standardized templates for commits, PRs, reviews |
| ADRs | **dojo records** | Architecture decisions preserved for the school |
| MCP server | **the dojo** | Where kata are served and practiced |
| Framework | **the ryu** (school) | The complete system of practices |

## Design philosophy

Kata in this repo are **platform-agnostic**. They follow the
[Agent Skills specification](https://github.com/anthropics/agent-skills-spec)
and work with any agent that supports it — Kiro, Claude Code,
OpenCode, or others. Skills avoid coupling to specific tooling,
MCP servers, or platform-specific features so they remain
portable across environments.

## Using kata from this repository

Each kata lives in its own folder under `skills/`. To use a kata:

1. **Copy the skill folder** into your project's `.agents/skills/` directory
2. Your agent will discover it automatically on the next session

```bash
cp -r skills/create-adr /path/to/your-project/.agents/skills/
```

Or use the dojo (MCP server) to install kata into your project — see
[cmd/openkata-mcp/README.md](cmd/openkata-mcp/README.md) for setup.

## Available kata

| Kata | Description |
|------|-------------|
| [create-adr](skills/create-adr/) | Detects architectural decisions in conversation and guides creation of Architecture Decision Records |
| [commit-conventions](skills/commit-conventions/) | Enforces Conventional Commits format and branch naming conventions |
| [create-skill](skills/create-skill/) | Creates agent skills by investigating repo conventions, designing workflows, and writing SKILL.md files |
| [create-rule](skills/create-rule/) | Creates always-on agent rules by investigating repo conventions and writing RULE.md files |
| [grill-me](skills/grill-me/) | Challenges a plan, spec, or ADR by interviewing the user until all corners are covered |
| [makefile-conventions](skills/makefile-conventions/) | Structures Makefiles as a universal command interface using modular includes and self-documenting help |
| [spec-workflow](skills/spec-workflow/) | Drives feature development through a phased workflow: specify, design, tasks, implement, and validate |

## Available dojo kun

| Rule | Description |
|------|-------------|
| [markdown-style](rules/markdown-style/) | Consistent markdown formatting conventions applied to all generated files |
| [bash-style](rules/bash-style/) | Bash scripting conventions based on the Google Shell Style Guide |
| [design-system](rules/design-system/) | Token-based CSS design system conventions |

## Available sensei profiles

| Profile | Description |
|---------|-------------|
| [frontend](profiles/frontend.md) | Frontend developer scoped to web UI — templates, styles, and handlers only |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to add skills
and rules. See [AGENTS.md](AGENTS.md) for build commands, code
style, and commit conventions. See [RELEASING.md](RELEASING.md)
for how to publish skills and deploy.

## License

MIT — see [LICENSE](LICENSE).
