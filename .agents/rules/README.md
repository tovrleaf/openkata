# Agent Rules (local)

This directory contains symlinks to rules from `rules/` that are
actively used in this project.

## Symlink Pattern

The `rules/` directory at the project root is the distributable
catalog. `.agents/rules/` is the discovery path for agents. By
symlinking only the rules we want active, we avoid duplication
while keeping distribution and local use decoupled.

## Adding a rule for local use

```bash
cd .agents/rules
ln -s ../../rules/<rule-name> <rule-name>
```

## Removing a rule from local use

```bash
rm .agents/rules/<rule-name>  # removes symlink only
```
