# Open Kata

A collection of Agent Skills — reusable capabilities and expertise for AI coding agents.

## Using skills from this repository

Each skill lives in its own folder under `skills/`. To use a skill:

1. **Copy the skill folder** into your project's `.agents/skills/` directory
2. Your agent will discover it automatically on the next session

Example:
```bash
# Copy a skill into your project
cp -r skills/create-adr /path/to/your-project/.agents/skills/
```

## Available skills

| Skill | Description |
|-------|-------------|
| [create-adr](skills/create-adr/) | Detects architectural decisions in conversation and guides creation of Architecture Decision Records |

## For contributors

Skills in this repo follow the Agent Skills specification. Each skill folder must contain a `SKILL.md` file with YAML frontmatter and markdown instructions.

### Local development

This repo uses skills from `skills/` locally via symlinks in `.agents/skills/`. To set up:

```bash
# Symlinks are already committed — just clone and go
git clone https://github.com/tovrleaf/openkata.git
```

See [.agents/skills/README.md](.agents/skills/README.md) for details on the symlink pattern.

## License

TBD
