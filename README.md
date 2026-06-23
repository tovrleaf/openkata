# Open Kata

<!-- badges:start -->
[![Skills](https://img.shields.io/badge/skills-13-blue?style=flat-square)](skills/)
[![Rules](https://img.shields.io/badge/rules-5-blue?style=flat-square)](rules/)
[![Profiles](https://img.shields.io/badge/profiles-2-blue?style=flat-square)](profiles/)
[![MCP Server](https://img.shields.io/badge/MCP_server-ready-green?style=flat-square)](https://openkata.dev/getting-started/)
[![Build](https://github.com/tovrleaf/openkata/actions/workflows/build.yaml/badge.svg)](https://github.com/tovrleaf/openkata/actions/workflows/build.yaml)
[![Go](https://img.shields.io/badge/go-1.26.1-blue?style=flat-square)](https://go.dev/)
[![Stars](https://img.shields.io/github/stars/tovrleaf/openkata?style=flat-square)](https://github.com/tovrleaf/openkata)
[![License](https://img.shields.io/badge/license-MIT-lightgrey?style=flat-square)](LICENSE)
<!-- badges:end -->

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

## Using kata

The fastest way to use kata is via the MCP server — add
`https://openkata.dev/mcp` to your agent config. Zero install,
always up-to-date.

For offline or local-first workflows, download and extract
skill folders into your project's `.agents/skills/<name>/`
directory. Your agent discovers them automatically.

- Skills live in `skills/`, rules in `rules/`, profiles in
  `profiles/`
- See the [full catalog](https://openkata.dev/catalog/) for
  all available artifacts
- See [Getting Started](https://openkata.dev/getting-started/)
  for client-specific setup (Kiro, Claude Code, OpenCode,
  and others)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to add skills
and rules. See [AGENTS.md](AGENTS.md) for build commands, code
style, and commit conventions. See [RELEASING.md](RELEASING.md)
for how to publish skills and deploy.

## License

MIT — see [LICENSE](LICENSE).
