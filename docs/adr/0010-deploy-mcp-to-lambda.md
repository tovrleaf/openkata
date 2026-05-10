---
status: PROPOSED
date: 2026-05-10
authors: [niko.kivela]
---

# 0010. Deploy MCP server to AWS Lambda

## Context

ADR 0002 established the MCP server as the distribution
mechanism for OpenKata artifacts. The server currently runs
locally via stdio or HTTP (`OPENKATA_ADDR`). To make skills
and rules discoverable and installable by anyone without
cloning the repository, the MCP server needs a public
deployment.

The web server (ADR 0007) already runs on Lambda with a
Function URL. The same pattern works for MCP — Streamable
HTTP transport over a Function URL.

Key constraint: Lambda has no git repository at runtime, so
version resolution via `git tag` (ADR 0005) won't work.
Skills and rules must be embedded or bundled.

## Decision Drivers

- Must serve skills/rules without a git clone at runtime
- Must resolve versions without git tags
- Must work with Streamable HTTP (MCP transport)
- Should reuse existing infra patterns (Lambda, Function URL)
- Should deploy via the same CI/CD pipeline
- Public access for v1 — no auth required

## Decision

Deploy the MCP server as a separate Lambda function
(`openkata-mcp`) with embedded skills and rules, a
DynamoDB table for download counts, and its own
Function URL.

### Skill storage

Released skills and rules are stored in S3. A separate
CI workflow triggers on tag creation:

1. CI checks out the repo with full history
2. Extracts the skill/rule directory at the tagged commit
   (`git archive <tag> -- <path>`)
3. Uploads the archive to S3 at
   `s3://openkata-artifacts/skills/<name>/<version>/`
4. Updates `versions.json` in S3 (list of all released
   artifacts and their latest versions)

The Lambda reads from S3 at runtime. No skills are
embedded in the binary — the Lambda stays thin and
deploys don't require full git history.

### Version resolution

`versions.json` in S3 maps artifact names to their
latest released version:

```json
{
  "skills/create-adr": "1.1.0",
  "skills/commit-conventions": "2.0.0",
  "rules/bash-style": "0.1.0"
}
```

The Lambda reads this on startup (cached in memory).
Updated only when a new tag is pushed.

### Transport

Streamable HTTP via `mcp-go`'s `StreamableHTTPServer`,
fronted by Lambda Function URL. Same adapter pattern as
the web server (`httpadapter.NewV2`).

Lambda works as an MCP server because the Streamable HTTP
transport (2025-03-26 spec) is stateless request/response
over HTTP — no persistent connections required. Each MCP
tool call is an independent HTTP POST. This fits Lambda's
execution model perfectly: request arrives, handler runs,
response returns. No WebSockets, no long-lived sessions.

The older SSE (Server-Sent Events) transport would not
work on Lambda due to its streaming/connection-hold
requirements. Streamable HTTP was designed specifically
to support serverless deployments.

The server runs in stateless mode (`WithStateLess(true)`).
No sessions, no server-initiated notifications — pure
request/response. Lambda has no memory between
invocations, so stateful sessions are impossible. Our
tools (list, install) are all single request/response
with no need for progress streaming or push
notifications.

### Infrastructure

- Separate Lambda: `openkata-mcp`
- S3 bucket: `openkata-artifacts` (skill/rule archives)
- DynamoDB table: `openkata-downloads` (atomic counters)
- Same region: `eu-north-1`
- Function URL: public, no auth
- MCP Lambda role: S3 read + DynamoDB read/write
- Web Lambda role: S3 read + DynamoDB read/write
- CI workflows: deploy (shallow clone), publish on tag
  (full clone, uploads to S3)

### Install behavior

`install_skill` and `install_rule` take the artifact
name and an optional version parameter. The server
returns:

- All files with their relative paths and content
- A pre-generated `.manifest.json` with name, version,
  source, installedAt, and per-file checksums

If no version is specified, the latest is served. The
calling agent writes files to `.agents/skills/<name>/`
or `.agents/rules/<name>/` in its project.

### Version listing

A separate `skill_versions` tool returns all available
versions for a given skill (resolved by listing S3
prefixes). Not included in `list_skills` to keep that
response lightweight.

Web downloads support versioning via query parameter:
- `/skills/create-adr/download` → latest
- `/skills/create-adr/download?v=1.0.0` → pinned

### Manifest format

Extends ADR 0005's manifest with per-file checksums for
modification detection:

```json
{
  "name": "create-adr",
  "version": "1.0.0",
  "source": "github.com/tovrleaf/openkata",
  "installedAt": "2026-05-10T20:00:00Z",
  "checksums": {
    "SKILL.md": "sha256:abc123...",
    "references/checklist.md": "sha256:def456..."
  }
}
```

Checksums enable the update flow to detect local
modifications before overwriting.

### Update behavior

No separate update tools on the server. `install_skill`
and `install_rule` always return the latest released
content with checksums. The calling agent is responsible
for:

1. Checking if the skill/rule already exists locally
2. Reading the existing `.manifest.json`
3. Comparing checksums to detect local modifications
4. Warning the user before overwriting modified files

Install and update are the same server operation. The
update intelligence lives in the agent, which has
filesystem access and user context.

### Tags and categories

Skills and rules carry tags in frontmatter `metadata`
for filtering and discovery:

```yaml
---
name: create-adr
description: Creates Architecture Decision Records
metadata:
  tags: "documentation, architecture, decisions"
---
```

Tags are used by:
- `list_skills`/`list_rules` — filter by tag
- Website — browsing and search
- Future: recommendation and related skills

### Download counting

Every install (via MCP) and download (via website tar.gz)
increments an atomic counter in DynamoDB. Single table,
partition key = artifact name, atomic `ADD` operation —
no read-modify-write races.

Both Lambdas have DynamoDB read/write access:
- MCP Lambda increments on `install_skill`/`install_rule`
- Web Lambda increments on tar.gz download
- Web Lambda reads counts for display on skill pages
- MCP Lambda reads counts for popularity queries

DynamoDB free tier (25 GB, 25 WCU/RCU) covers this
indefinitely at our scale.

### Rate limiting

Lambda reserved concurrency on both functions (e.g. 10)
caps total throughput at zero cost. Excess requests
receive 429 responses. This prevents runaway costs from
bot abuse on the public endpoints.

WAF rate-limiting rules via CloudFront are available as
a future upgrade (~$1.60/month) if per-IP throttling is
needed.

## Alternatives Considered

### Same Lambda as web server

- **Pros:** Less infrastructure, one deployment, one
  Function URL, simpler CI
- **Cons:** Larger binary, slower cold starts, coupled
  deployments, mixed concerns
- **Rejected because:** Separate binaries have faster
  cold starts and cleaner separation. Infra complexity
  is acceptable.

### Separate Lambda

- **Pros:** Clean separation of concerns, independent
  scaling, faster cold starts, smaller binaries
- **Cons:** Two deploys, two Function URLs, more CI steps
- **Chosen because:** Better cold-start performance for
  MCP tool calls, clean separation, independent
  deployment lifecycle.

### S3 for skill storage

- **Pros:** Update skills without redeploying, no binary
  bloat, CI only needs full history on tag (not deploy),
  can serve individual files
- **Cons:** Extra infra (bucket, IAM), S3 read latency
- **Chosen because:** Decouples skill releases from
  Lambda deploys. Deploy CI stays fast with shallow
  clone. Skills can grow without affecting binary size.

### Embedding skills in binary

- **Pros:** No external deps at runtime, fast reads
- **Cons:** Every skill change requires redeploy, CI
  needs full git history every deploy, binary grows
  unbounded, large responses for big skills
- **Rejected because:** Couples skill releases to Lambda
  deploys and requires full git history on every CI run.

### Keep MCP local-only, use HTTP API for remote

- **Pros:** Simpler, no MCP transport concerns
- **Cons:** Agents already speak MCP natively. A custom
  HTTP API requires each agent to implement a client.
- **Rejected because:** MCP is the standard protocol
  agents already support.

## Consequences

### Positive

- Public MCP endpoint for skill discovery and install
- Skills update without redeploying Lambda
- Small Lambda binary — fast cold starts
- Deploy workflow uses shallow clone (fast CI)
- Agents can connect directly via MCP config
- Checksums protect users from silent overwrites
- Tags enable discovery via MCP and website
- Download counts provide visibility into adoption

### Negative

- S3 bucket adds infrastructure
- DynamoDB table adds infrastructure (minimal)
- S3 read latency on cold start (~50-100ms)
- Tags require frontmatter discipline from authors

### Neutral

- Auth can be layered on later via CloudFront or
  Lambda authorizer
- ADR 0005 manifest format is extended, not replaced
