# OpenKata MCP Server

A local MCP server that distributes OpenKata skills to your projects.

## Build

```bash
go build -o bin/openkata-mcp ./cmd/openkata-mcp/
```

## Configure

Add to your project's `.kiro/settings/mcp.json`:

```json
{
  "mcpServers": {
    "openkata": {
      "command": "<OPENKATA_REPO>/bin/openkata-mcp",
      "args": [],
      "cwd": "<OPENKATA_REPO>",
      "env": {
        "OPENKATA_SKILLS_DIR": "<OPENKATA_REPO>/skills"
      }
    }
  }
}
```

Replace `<OPENKATA_REPO>` with the absolute path to your clone of
this repository.

## Tools

| Tool | Description |
|------|-------------|
| `list_skills` | Lists available skills with their descriptions |
| `install_skill` | Copies a skill into `<target_dir>/.agents/skills/` |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OPENKATA_SKILLS_DIR` | `skills` (relative to cwd) | Path to the directory containing skills |
| `OPENKATA_ADDR` | *(unset)* | Set to `host:port` to run as HTTP server instead of stdio |
