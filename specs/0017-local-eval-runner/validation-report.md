# Validation Report: 0017-local-eval-runner

Date: 2026-06-19
Validator: spec-validator agent
Build: PASS (go build + go test both exit 0)

## Summary

The implementation covers the core eval runner logic —
config resolution, scenario discovery, scoring, judge
parsing, terminal output, JSON output, parallel
execution, and retry logic. However, several spec
requirements have gaps or deviations.

Overall: 19 requirements met, 7 gaps/deviations found.

## Requirements Status

### CLI Binary

| Requirement | Status | Notes |
|-------------|--------|-------|
| Binary at `cmd/openkata-eval/main.go` | ✓ PASS | |
| Invocation: `./bin/openkata-eval skills/<name>` | ✓ PASS | |
| Single scenario mode: `skills/<name>/evals/scenario-0` | ✓ PASS | |
| Full skill mode: exit 0 if >= threshold, exit 1 otherwise | ✓ PASS | |
| Single scenario mode: always exit 0 | ✓ PASS | |

### Eval File Format

| Requirement | Status | Notes |
|-------------|--------|-------|
| Reads scenario.json, criteria.json, task.md, inputs/ | ✓ PASS | |
| Criteria have name, description, max_score | ✓ PASS | |

### LLM Backend

| Requirement | Status | Notes |
|-------------|--------|-------|
| Completer interface with Complete(system, user) | ✓ PASS | |
| kiro-cli backend (default) | ✓ PASS | |
| HTTP backend (OpenAI-compatible) | ✗ FAIL | URL construction bug (see below) |
| Direct mode uses eval-agent.json in temp dir | ✗ FAIL | Not implemented (see below) |
| Backend selection via --backend flag | ✓ PASS | |
| Judge output parsing: strip ANSI, skip headers, find JSON | ✓ PASS | |
| Retry once on judge parse failure | ✓ PASS | |

### Sandbox Execution

| Requirement | Status | Notes |
|-------------|--------|-------|
| Docker container lifecycle | ✓ PASS | Code exists in sandbox.go |
| Mount ~/.aws/ read-only | ✓ PASS | |
| Copy skill files into container | ✓ PASS | |
| Copy scenario inputs | ✓ PASS | |
| Run kiro-cli --agent eval-agent --trust-all-tools | ✓ PASS | |
| Generate eval-agent.json inside container | ✓ PASS | |
| Workspace diff (before/after) | ✓ PASS | |
| Sandbox wired into main.go | ✗ FAIL | Not connected (see below) |
| Docker availability check | ✗ FAIL | CheckDocker exists but never called |
| kiro-cli availability check | ✗ FAIL | Not implemented |
| Container timeout → kill + fail scenario | ✓ PASS | Uses `timeout` command wrapper |

### Configuration

| Requirement | Status | Notes |
|-------------|--------|-------|
| Defaults (backend: kiro, model: claude-sonnet-4.6, threshold: 95) | ✓ PASS | |
| YAML config file (.openkata-eval.yaml) | ✓ PASS | |
| Environment variables | ✓ PASS | |
| CLI flags override | ✓ PASS | |
| Resolution order: defaults → yaml → env → flags | ✓ PASS | |
| Judge model defaults to agent model | ✓ PASS | |
| Concurrency configurable (default: 2) | ✓ PASS | |
| Sandbox config (image, timeout, network) | ✓ PASS | |

### Execution Flow

| Requirement | Status | Notes |
|-------------|--------|-------|
| Agent prompt: SKILL.md + references + scripts + assets + task.md + inputs | ✓ PASS | |
| Excludes CHANGELOG.md, RATIONALE.md, .tessl-plugin/, tile.json, .tesslignore, evals/ | ✓ PASS | |
| Excludes references/ACKNOWLEDGMENTS.md | ✓ PASS | |
| Inputs as fenced code blocks with paths as labels | ✓ PASS | |

### Judge Strategy

| Requirement | Status | Notes |
|-------------|--------|-------|
| Batch all criteria in single judge call | ✓ PASS | |
| Sanity check: trivial response + all pass → fail as suspicious | ✓ PASS | |
| Threshold: <20 chars | ✓ PASS | |

### Scoring

| Requirement | Status | Notes |
|-------------|--------|-------|
| Per criterion: pass → max_score, fail → 0 | ✓ PASS | |
| Per scenario: percentage of earned vs max | ✓ PASS | |
| Overall: average percentage across scenarios | ✓ PASS | |
| Threshold: 95% default, configurable | ✓ PASS | |

### Parallel Execution

| Requirement | Status | Notes |
|-------------|--------|-------|
| Scenarios run in parallel | ✓ PASS | |
| Bounded concurrency | ✓ PASS | |
| Results printed in order | ✓ PASS | |

### Terminal Output

| Requirement | Status | Notes |
|-------------|--------|-------|
| Colored ✓/✗ per criterion with scores | ✓ PASS | |
| Scenario score line with PASS/FAIL | ✓ PASS | |
| Overall line with threshold | ✓ PASS | |

### JSON Output

| Requirement | Status | Notes |
|-------------|--------|-------|
| Timestamped filename | ✓ PASS | |
| Directory support | ✓ PASS | |
| JSON structure matches spec | ✓ PASS | |

### Makefile Integration

| Requirement | Status | Notes |
|-------------|--------|-------|
| `eval-local` target in mk/dev.mk | ✓ PASS | |
| Build binary target | ✓ PASS | |

### Error Handling

| Requirement | Status | Notes |
|-------------|--------|-------|
| Missing SKILL.md → skip with error | ✓ PASS | |
| No scenarios → fail with error | ✓ PASS | |
| Rate limiting → retry with exponential backoff (3 attempts) | ✓ PASS | |
| Malformed judge JSON → retry once | ✓ PASS | |
| kiro-cli not found → fail fast | ✗ FAIL | No check |
| Docker not running → fail fast (if sandboxed needed) | ✗ FAIL | CheckDocker not called |

### Dependencies

| Requirement | Status | Notes |
|-------------|--------|-------|
| gopkg.in/yaml.v3 in go.mod | ✓ PASS | |
| os/exec for kiro-cli and docker | ✓ PASS | |
| net/http for HTTP backend | ✓ PASS | |
| No other external dependencies | ✓ PASS | |

### Dockerfile.eval

| Requirement | Status | Notes |
|-------------|--------|-------|
| Ubuntu 24.04 base | ✓ PASS | |
| Packages: curl, git, bash, jq | ✓ PASS | |
| WORKDIR /workspace | ✓ PASS | |
| kiro-cli installation | ~ PARTIAL | Comment says "install manually" instead of COPY --from |

## Detailed Findings

### FAIL 1: Sandbox not wired into main.go

The `Runner` struct accepts a `Sandbox SandboxRunner` field
and the `runScenario` method checks `s.Sandbox && r.Sandbox
!= nil` to route to the sandbox path. However, `main.go`
never assigns this field. `NewSandboxRunner` and
`CheckDocker` are defined in `sandbox.go` but never called
from `main.go`.

Impact: Sandboxed scenarios silently fall through to direct
mode (raw completion), ignoring the `"sandbox": true` flag.
The safety constraint ("--trust-all-tools must NEVER be used
outside a container") cannot be enforced.

### FAIL 2: Direct mode does not use eval-agent.json

The spec describes direct execution (sandbox: false) as:

> Generates eval-agent.json in a temp directory on the host,
> runs kiro-cli with `--agent eval-agent --no-interactive`
> (no tools flag).

The implementation's `KiroCompleter.Complete()` runs
`kiro-cli chat --no-interactive "<prompt>"` — a raw
completion without an agent definition. This means the skill
is not loaded through kiro-cli's native agent mechanism as
specified.

Impact: Agent behavior in direct mode may differ from what
users see when the skill is installed normally.

### FAIL 3: HTTP backend URL double-path bug

`HTTPCompleter.Complete()` constructs:

```go
url := h.BaseURL + "/v1/chat/completions"
```

The spec's example config sets `base_url:
http://localhost:11434/v1`. This produces
`http://localhost:11434/v1/v1/chat/completions` — an
invalid URL with duplicated `/v1`.

Either the code should strip a trailing `/v1` from
BaseURL, or the config example should use the base
without `/v1` (e.g., `http://localhost:11434`), or the
code should use `/chat/completions` without the `/v1`
prefix since it's already in the base URL.

### FAIL 4: No kiro-cli availability check

The spec requires: "kiro-cli not found → fail fast,
suggest installation." No `exec.LookPath("kiro-cli")`
check exists. The runner will attempt execution and
produce a cryptic error from os/exec instead of a
user-friendly message.

### FAIL 5: CheckDocker never called

`CheckDocker()` exists in `sandbox.go` and checks both
PATH availability and `docker info`. However, it is never
invoked from `main.go`. The spec requires: "Docker not
running → fail fast only if the current run includes
sandboxed scenarios."

### FAIL 6: Direct-call timeout not configurable

The spec says "Timeout: 300s per sandboxed scenario,
120s per direct call (configurable)." The 120s timeout
is hardcoded in `main.go`:

```go
NewKiroCompleter(model, 120*time.Second)
NewHTTPCompleter(cfg.HTTP.BaseURL, ..., 120*time.Second)
```

No config field or flag allows overriding the
direct-call timeout.

### DEVIATION: copySkillFiles method is a no-op

`DockerSandbox.Run()` calls `copySkillFiles()` which
returns `("", nil)` immediately with a comment saying
"For now, we just ensure /workspace/skill/ exists."
However, `sandboxRun()` (the properly wired version)
implements `copySkillToContainer()` correctly. This is
a dead code issue — `DockerSandbox.Run()` is not the
active path; `sandboxRun()` is the intended entrypoint.
Not blocking since `NewSandboxRunner` uses `sandboxRun`.

## Out-of-Scope Verification

| Item | Built? | Status |
|------|--------|--------|
| Running evals in CI | No | ✓ Correct |
| Conversation/multi-turn evals | No | ✓ Correct |
| Caching agent responses | No | ✓ Correct |
| Comparing results between models | No | ✓ Correct |
| Streaming responses | No | ✓ Correct |

## Test Coverage

Tests cover: config loading, scenario discovery, agent
prompt construction (inclusion/exclusion), judge parsing
(ANSI, headers, retry, suspicious), scoring (per-scenario,
overall, exit codes), HTTP completer, kiro output cleaning,
sandbox commands, workspace diffing, JSON output, parallel
execution ordering, and throttle retry. All tests pass.

## Verdict

The core evaluation logic is solid and well-tested. The
primary gap is the sandbox integration not being wired
into `main.go`, which means the sandbox code path is
unreachable in production. The direct-mode agent call
deviates from the spec's design (should use kiro-cli's
agent mechanism, not raw chat). The HTTP URL bug would
break anyone using the documented config example.

Recommended priority for fixes:
1. Wire sandbox into main.go (critical — safety constraint)
2. Fix HTTP URL construction (breaks documented usage)
3. Add kiro-cli / Docker availability checks (UX)
4. Implement eval-agent.json for direct mode (spec fidelity)
5. Make direct-call timeout configurable (minor)
