# Agent Skills (local)

This directory contains symlinks to skills from `skills/` that are actively used in this project.

## Why symlinks?

The `skills/` directory at the project root is the distributable catalog — it contains all skills this project offers. But not every skill in the catalog needs to be active locally.

`.agents/skills/` is the standard discovery path for agents (VS Code/Copilot, Claude Code, etc.). By symlinking only the skills we want to use, we avoid duplication while keeping distribution and local use decoupled.

## Adding a skill for local use

```bash
cd .agents/skills
ln -s ../../skills/<skill-name> <skill-name>
```

## Removing a skill from local use

```bash
rm .agents/skills/<skill-name>  # removes symlink only, not the skill itself
```
