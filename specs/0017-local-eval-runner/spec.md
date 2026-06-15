---
status: Draft
depth: Standard
---

# Local Eval Runner

## Story

As a skill author, I want to run eval scenarios locally
so I can iterate on skills without depending on external
services.

## Requirements

### CLI Binary

- New binary: `cmd/openkata-eval/main.go`
- Invocation: `./bin/openkata-eval skills/<name>`
- Also: `./bin/openkata-eval skills/<name>/evals/scenario-0`
  to run a single scenario
- **Full skill mode** (all scenarios): exit 0 if
  overall average >= threshold, exit 1 otherwise.
  This is the quality gate.
- **Single scenario mode**: always exit 0. Reports the
  score for debugging/iteration but does not gate on
  threshold.

### Eval File Format (unchanged)

Existing structure consumed as-is:

```
skills/<name>/evals/
  scenario-0/
    scenario.json   # { "description": "..." }
    criteria.json   # { "context", "type", "checklist": [...] }
    task.md         # user prompt
    inputs/         # optional context files
```

Criteria checklist items have `name`, `description`,
and `max_score` (integer weight).

### LLM Backend

Two execution paths depending on sandbox mode:

**Direct execution (sandbox: false)**

For text-only skills. Generates eval-agent.json in a
temp directory on the host, runs kiro-cli with
`--agent eval-agent --no-interactive` (no tools flag).
Judge call also uses kiro-cli via Completer.

```
Runner → generate eval-agent.json (temp dir)
Runner → kiro-cli --agent eval-agent --no-interactive
       → agent response (text)
Runner → Completer → judge verdict (JSON)
```

**Sandboxed execution (sandbox: true)**

For skills that use tools. Agent runs as a full
session inside Docker with `--agent eval-agent
--trust-all-tools`. Judge still uses Completer:

```
Runner → Docker container
       → kiro-cli --agent eval-agent --trust-all-tools
       → stdout + workspace artifacts
Runner → diff workspace (before/after)
Runner → Completer → judge verdict (JSON)
```

Both modes load the skill through kiro-cli's native
agent mechanism for consistent behavior. The only
difference is tool trust and isolation.

The `Completer` interface:

```go
type Completer interface {
    Complete(system, user string) (string, error)
}
```

Two implementations:

**kiro-cli (default)** — shells out to:

```bash
kiro-cli chat --model <model> --no-interactive "<prompt>"
```

**HTTP** — OpenAI-compatible chat completions API for
local models (Ollama) or direct cloud access (Bedrock).

Backend selection via `--backend` flag or config:
- `kiro` (default) — uses kiro-cli
- `http` — uses OpenAI-compatible endpoint

The Completer is used for:
- Agent calls in direct mode (sandbox: false)
- Judge calls in both modes (always single-shot)

Judge output parsing must be robust: strip ANSI codes,
skip header lines (`> `), find the first `[` or `{`,
and parse JSON from there. Retry once on parse failure.
If the second attempt also fails, fail the scenario
with a parse error.

### Sandbox Execution

All agent calls run inside a Docker container for
isolation. The agent can use tools, run scripts, and
write files without affecting the host machine.

**Per-scenario flow:**

1. Start container from eval base image
2. Mount `~/.aws/` read-only (kiro-cli authenticates
   via AWS SSO tokens stored there)
3. Copy skill files (SKILL.md, references/, scripts/,
   assets/) into container
4. Copy scenario inputs (task.md, inputs/) into
   container working directory
5. Run kiro-cli inside the container with
   `--agent eval-agent --trust-all-tools --no-interactive`
   using a generated agent definition that references
   the skill's files via `resources` (SKILL.md,
   references/, scripts/, assets/). This lets kiro-cli
   load the skill through its native discovery
   mechanism rather than concatenating into a prompt.
6. Capture stdout and copy produced artifacts out
7. Destroy container

The runner generates a temporary `.kiro/agents/eval-agent.json`
inside the container:

```json
{
  "name": "eval-agent",
  "tools": ["fs_read", "fs_write", "execute_bash",
            "grep", "glob"],
  "allowedTools": ["fs_read", "fs_write",
                   "execute_bash", "grep", "glob"],
  "resources": [
    "file://skill/SKILL.md",
    "file://skill/references/**",
    "file://skill/scripts/**",
    "file://skill/assets/**"
  ]
}
```

**Base image** (`Dockerfile.eval`):

```dockerfile
FROM ubuntu:24.04
RUN apt-get update && apt-get install -y \
    curl git bash jq
# Install kiro-cli
COPY --from=kiro /usr/local/bin/kiro-cli /usr/local/bin/
WORKDIR /workspace
```

The exact image contents depend on what tools skills
need. Start minimal, extend as needed.

**Configuration:**

```yaml
# .openkata-eval.yaml
sandbox:
  image: openkata-eval:latest
  timeout: 300  # seconds per scenario
  network: host # none | host (default: host)
```

Network mode defaults to `host` because kiro-cli
inside the container needs to reach the LLM backend
(Bedrock). Set to `none` only for scenarios that
pre-load a model response (future: offline replay
mode). Per-scenario override available:

```json
{
  "description": "...",
  "sandbox": true,
  "network": "none"
}
```

**Artifact extraction:**

After the agent finishes, the runner copies the
container's `/workspace` directory out and makes it
available to the judge. The runner compares before
and after states to produce a diff summary:

- **Added**: files that didn't exist before
- **Modified**: files whose contents changed
- **Deleted**: files that were removed
- **Unchanged**: files that remain identical

The judge receives:
- Agent stdout (text output from kiro-cli)
- Diff summary (which files were added/modified/
  deleted/unchanged)
- Full contents of added and modified files

This lets MUST_NOT criteria ("does not modify X") be
evaluated unambiguously — the judge knows what the
agent *changed*, not just what exists.

**No sandbox fallback:**

For text-only skills that don't use tools (e.g.,
critical-thinking, grill-me), the runner can skip
Docker and use a direct kiro-cli or HTTP call. This
is controlled by the scenario:

```json
{
  "description": "...",
  "sandbox": false
}
```

Default: `"sandbox": true`. Skills that only produce
conversational output can opt out for faster execution.

**Safety constraint:**

`--trust-all-tools` must NEVER be used outside a
container. When `sandbox: false`, the agent call uses
kiro-cli without tool trust (or the HTTP backend for
a raw completion). This is a hard rule, not
configurable.

**Note:** kiro-cli `--no-interactive` auto-approves
tool use regardless of `--trust-all-tools`. The two
flags are equivalent from a safety perspective. This
means `sandbox: false` is only safe for skills whose
tasks will not trigger tool use (purely conversational
skills). If in doubt, use `sandbox: true`.

### Configuration

Resolved in order (later wins):

1. Defaults (backend: kiro, model: claude-sonnet-4.6,
   threshold: 95)
2. Config file: `.openkata-eval.yaml` in project root
3. Environment variables
4. CLI flags

```yaml
# .openkata-eval.yaml
backend: kiro
model: claude-sonnet-4.6
threshold: 95

# Optional: separate judge model
judge_model: claude-sonnet-4.6

# Only needed for http backend:
http:
  base_url: http://localhost:11434/v1
  api_key: ""
```

Environment variables:

- `OPENKATA_EVAL_BACKEND` (kiro | http)
- `OPENKATA_EVAL_MODEL`
- `OPENKATA_EVAL_JUDGE_MODEL`
- `OPENKATA_EVAL_THRESHOLD`
- `OPENKATA_EVAL_HTTP_BASE_URL`
- `OPENKATA_EVAL_HTTP_API_KEY`

CLI flags: `--backend`, `--model`, `--judge-model`,
`--threshold`, `--output`

If only one model is configured, it is used for both
agent and judge.

### Execution Flow

For each scenario:

1. **Agent call** — single chat completion:
   - Skill context (concatenated into the prompt):
     - Contents of `SKILL.md`
     - All files in `references/` (excluding
       `ACKNOWLEDGMENTS.md`)
     - All files in `scripts/` (as readable context,
       not executed)
     - All files in `assets/` (if directory exists)
   - Excluded from context (framework internals):
     - `CHANGELOG.md`, `RATIONALE.md`
     - `.tessl-plugin/`, `tile.json`, `.tesslignore`
     - `evals/`
   - Task (appended after skill context):
     - Contents of `task.md`
     - If `inputs/` exists in the scenario, append file
       contents as fenced code blocks with paths as labels
   - This mirrors what agents receive when a skill is
     installed per the Agent Skills specification.
2. **Judge call** — single chat completion:
   - System message: judge prompt (score the response
     against criteria, output structured JSON)
   - User message: the agent's response + full criteria
     checklist
   - Expected output: JSON array of
     `{ "name": "...", "pass": bool, "reason": "..." }`

### Judge Strategy

Batch all criteria into a single judge call per
scenario. The judge prompt instructs the model to
evaluate each criterion independently and return a
JSON verdict per criterion.

**Sanity check:** If the agent response is empty or
under 20 characters and the judge passes all criteria,
flag the scenario as suspicious and fail it. This
guards against a rubber-stamping judge.

### Scoring

- Per criterion: pass → `max_score` points, fail → 0
- Per scenario: `(sum of passed scores / sum of all
  max_scores) * 100`
- Overall: average percentage across all scenarios
- Threshold: 95% (configurable via `--threshold`)

### Parallel Execution

Scenarios run in parallel by default. Each scenario
gets its own container (sandboxed) or goroutine
(text-only). Concurrency is configurable:

```yaml
# .openkata-eval.yaml
concurrency: 4  # max parallel scenarios
```

Default: 2. Keeps LLM rate limits manageable while
still reducing wall time. Increase if your backend
supports higher throughput. Results print as each
scenario completes (not buffered). Overall summary
prints last.

### Terminal Output

```
skills/commit-conventions
  scenario-0: Atomic commits with conventional format
    ✓ Atomic commits (12/12)
    ✓ No bulk staging (10/10)
    ✗ Commit type: fix (0/8) — used 'bugfix' instead of 'fix'
    ...
    Score: 88/100 (88%) FAIL

  scenario-1: ...
    Score: 100/100 (100%) PASS

  Overall: 94% FAIL (threshold: 95%)
```

### JSON Output

When `--output results.json` is specified, the runner
appends a timestamp to the filename to avoid
overwriting previous runs:
`results-2026-06-18T110000.json`. If `--output` is a
directory, files are written there with auto-generated
names.

```json
{
  "skill": "commit-conventions",
  "timestamp": "2026-06-16T01:00:00Z",
  "config": {
    "backend": "kiro",
    "agent_model": "claude-sonnet-4.6",
    "judge_model": "claude-sonnet-4.6"
  },
  "threshold": 95,
  "scenarios": [
    {
      "name": "scenario-0",
      "description": "Atomic commits with conventional format",
      "criteria": [
        { "name": "...", "pass": true, "score": 12,
          "max_score": 12, "reason": "" }
      ],
      "score": 88,
      "max_score": 100,
      "percentage": 88,
      "pass": false
    }
  ],
  "overall_percentage": 94,
  "pass": false
}
```

### Makefile Integration

```makefile
.PHONY: eval-local
eval-local: bin/openkata-eval
	@./bin/openkata-eval $(filter-out $@,$(MAKECMDGOALS))

bin/openkata-eval:
	@go build -o bin/openkata-eval ./cmd/openkata-eval/
```

Usage: `make eval-local skills/commit-conventions`

### Error Handling

- Missing SKILL.md → skip skill with error message
- No eval scenarios found → fail with error
- Missing scenario files → skip scenario with warning
- kiro-cli not found → fail fast, suggest installation
- Docker not running → fail fast only if the current
  run includes sandboxed scenarios. Text-only scenarios
  (`"sandbox": false`) never require Docker.
- Container timeout → kill container, fail scenario
- Rate limiting (HTTP 429 / throttling) → retry with
  exponential backoff, up to 3 attempts
- LLM connection failure (HTTP backend) → fail fast
  with clear error
- Malformed judge JSON → retry once, then fail the
  scenario with a parse error
- Timeout: 300s per sandboxed scenario, 120s per
  direct call (configurable)

### Dependencies

- `os/exec` for kiro-cli backend and Docker commands
- `net/http` for HTTP backend (OpenAI-compatible API)
- `gopkg.in/yaml.v3` for config file parsing (already
  indirect in module, or add)
- Docker (runtime dependency for sandboxed execution)
- No other external dependencies

## Out of Scope

- Running evals in CI (future spec)
- Conversation/multi-turn evals
- Caching agent responses across runs
- Comparing results between models
- Streaming responses

## Open Questions

1. **kiro-cli output parsing** — Stripping ANSI codes
   and `> ` header lines is the plan. Need to confirm
   this is stable across kiro-cli versions.
2. **Judge prompt tuning** — The judge system prompt
   will be iterated during implementation. Not
   specified here.
3. **kiro-cli auth in container** — Exact path to
   mount for kiro-cli's session/credentials TBD at
   implementation time.

Date: 2026-06-16
