---
title: Installing artifacts through the MCP server
description: Connect an agent to the OpenKata MCP server, then list and install skills, rules and profiles.
order: 3
---

# Installing artifacts through the MCP server

The OpenKata MCP server, the dojo, serves skills, rules and profiles
to any agent that supports MCP. An install copies the files into
your project. After that the project works without the server, and
you go back to it only to install something new or to update.

## Connect your agent

Add the server URL to your agent's MCP configuration:

```text
https://openkata.dev/mcp
```

There is nothing to install and no local runtime. In Claude Code,
for example, add this to `.mcp.json`:

```json
{
  "mcpServers": {
    "openkata": {
      "url": "https://openkata.dev/mcp"
    }
  }
}
```

Other agents use the same URL in their own config file. The
[Getting Started](/getting-started/) page lists the file for each
client.

## The tools

The server has nine tools, three for each artifact type.

| Tool | Arguments | Returns |
|------|-----------|---------|
| `list_skills`, `list_rules`, `list_profiles` | `tag` (optional) | Every artifact of that type |
| `install_skill`, `install_rule`, `install_profile` | `skill`, `rule` or `profile` (required); `version` (optional) | The files and a manifest |
| `skill_versions`, `rule_versions`, `profile_versions` | `skill`, `rule` or `profile` (required) | Every available version |

You rarely call them yourself. You ask your agent in plain
language, and it picks the tool.

## List what is available

A list call returns the artifacts sorted by name. Each entry has
its name, latest version, description, tags and download count.

Pass a `tag` to narrow the list. The tags are the ones in an
artifact's metadata, such as `category:conventions` or
`tool:git`, and the filter must match one tag exactly. Asking for
`category` alone matches nothing.

For example: "List the OpenKata rules tagged `category:conventions`."

## Install an artifact

Ask your agent to install by name: "Install the `grill-me` skill
from OpenKata." The agent calls `install_skill` and gets back a
JSON document with the name, the version, a manifest and the
files. It then writes the files into your project:

| Type | Standard location (Kiro and others) | Claude Code |
|------|-------------------------------------|-------------|
| Skills | `.agents/skills/<name>/` | `.claude/skills/<name>/` |
| Rules | `.agents/rules/<name>/` | `.claude/rules/<name>/` |
| Profiles | `.agents/profiles/<name>/` | `.claude/profiles/<name>/` |

The server does not detect your client. Its tool descriptions
name both locations and the agent writes to the one its client
uses. If your project already follows one convention, tell the
agent to use it.

The manifest is saved as `.openkata-manifest.json` next to the
files. It records the artifact name, the version, the source
(`github.com/tovrleaf/openkata`) and a SHA-256 checksum for every
installed file.

An install leaves out the files that only matter to the project
that maintains the artifact: the changelog, the rationale, the
`evals/` folder and the tessl metadata.

## Choose a version

An install gets the latest version unless you ask for another one.
To pin one, name it in the request: "Install `grill-me` at version
`1.0.1`." To see what exists, ask for the versions of an artifact;
the version tools list them from oldest to newest.

If the name or the version does not exist, the tool returns an
error that says which one was not found, for example
`skill "x" not found`.

## Download instead

If you work offline, or your agent has no MCP support, each
artifact page on the site links a `.tar.gz` archive. Extract it
into the same folders as above. It holds the same files as an MCP
install. Browse the [skills](/skills/), [rules](/rules/) and
[profiles](/profiles/).

Each MCP install and each archive download counts as a download in
the public download count shown in the catalog.
