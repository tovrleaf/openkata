# Open Kata

Codified practices for AI agents — teach your agents the way.

## Vocabulary

| Concept | Kata term | Description |
|---------|-----------|-------------|
| Skills | **kata** | Codified practices agents follow |
| Roles | **sensei profiles** | Agent role definitions with scoped permissions |
| Prompts | **kata forms** | Standardized templates for commits, PRs, reviews |
| ADRs | **dojo records** | Architecture decisions preserved for the school |
| MCP server | **the dojo** | Where kata are served and practiced |
| Framework | **the ryu** (school) | The complete system of practices |

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

## For contributors

Kata in this repo follow the Agent Skills specification. Each skill folder
must contain a `SKILL.md` file with YAML frontmatter and markdown instructions.

### Local development

This repo uses skills from `skills/` locally via symlinks in `.agents/skills/`.

```bash
# Symlinks are already committed — just clone and go
git clone https://github.com/tovrleaf/openkata.git
```

See [.agents/skills/README.md](.agents/skills/README.md) for details on the symlink pattern.

## License

TBD
