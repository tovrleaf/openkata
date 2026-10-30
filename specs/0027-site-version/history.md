# Proposed Release History

A proposal for the tags that give the project its past, made from the
161 commits on `main` (first parent), 2026-04-12 to 2026-10-24. Each
release is placed at a commit where a coherent piece had landed, mostly
at the merge of a pull request. The numbers are minor bumps, as the
project has had no stable interface yet (`0.x`).

Nothing is tagged. If you accept it, the tags are made on these
commits, so each release has the date of its commit, and a GitHub
release is written for each from the summary below.

| Version | Date | Commit | What it holds |
|---------|------|--------|---------------|
| v0.1.0 | 2026-04-12 | `dfbe47f` | The first MCP server with HTTP transport, the ADR tooling, the kata vocabulary |
| v0.2.0 | 2026-04-16 | `5077932` | The artifact versioning scheme and manifest, release skill, Conventional Commits skill, Keep a Changelog |
| v0.3.0 | 2026-04-25 | `61cb255` | Rules as an artifact type, the Tessl toolchain, `create-skill`, `make skills` and `make rules` |
| v0.4.0 | 2026-04-30 | `b3139c7` | `create-skill` and `create-rule` as distributable skills, `bash-style`, `AGENTS.md`, `CONTRIBUTING.md` |
| v0.5.0 | 2026-05-05 | `b966859` | The spec workflow, the manifesto, ADRs 0007 to 0009, version resolution from git tags |
| v0.6.0 | 2026-05-10 | `f19992d` | The website, the Lambda infrastructure, the custom domain, CI/CD (specs 0001 to 0004) |
| v0.7.0 | 2026-05-15 | `a3fe2d8` | The design system, the MCP server deployment, consistent artifact names, versions reset to 1.x |
| v0.8.0 | 2026-05-26 | `9069862` | The skills listing, the skill detail page, the `create-profile` skill, the skill quality audit, publish fixes |
| v0.9.0 | 2026-06-03 | `4e622ad` | Rules pages, the CSS rename, logo and navigation, rule releases |
| v0.10.0 | 2026-06-11 | `291d74a` | Skill portability, the eval runner, rename redirects, the release skill, the CI trust policy |
| v0.11.0 | 2026-06-12 | `c173db4` | Catalog search, SVG navigation icons, the mobile navigation |
| v0.12.0 | 2026-06-14 | `ee0fd67` | Download analytics, the getting-started page, the stats dashboard |
| v0.13.0 | 2026-06-23 | `dc13ec2` | New skills and the `go-testing` rule |
| v0.14.0 | 2026-10-24 | `99128fd` | The `openkata` manifest name, the docs engine (PR 33), squash merging forbidden |

The next release, `v0.15.0`, is what is on `feature/0024-docs-content`
and the work since: the docs content, the MCP server version stamp, the
benchmarks and version list in production, `make published`, and the
site version itself. It is tagged when that branch is merged.

## Judgement calls

- **The gap between 2026-06-23 and 2026-10-16** has no commits on `main`.
  Releases do not need to be evenly spaced, so v0.14.0 follows v0.13.0
  directly.
- **v0.7.0 and v0.10.0 are the least clear.** Both merge several
  things that landed over a few days, and the commit I picked is the
  last of the group.
- **Local-only work** such as the docs route is in v0.14.0 because it
  was merged, but it is not visible on the live site yet. The release
  notes should say so rather than hide it.
- **Pull requests that squash** were disabled on 2026-10-16. Earlier
  history has a few squashed commits (for example `#18`), so one of
  those carries several changes, and is described by its message.
- A release between two of these can be inserted later, with a tag
  on an earlier commit, but the version numbers would then have to
  change, so decide the list before tagging.
