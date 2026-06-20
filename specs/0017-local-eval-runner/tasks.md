# Tasks: Local Eval Runner

## Tasks

### 1. CLI skeleton and config loading
- **Status**: Done
- **Goal**: New binary that parses CLI args, loads
  config from file/env/flags, validates a skill path
  argument, discovers scenario directories, reads
  sandbox/network settings from scenario.json
- **Boundary**: `cmd/openkata-eval/main.go`,
  `cmd/openkata-eval/config.go`
- **Key files**: `cmd/generate-versions/main.go`
  (CLI pattern), existing eval dirs in
  `skills/commit-conventions/evals/`
- **Depends**: None
- **Done when**: `go build -o bin/openkata-eval
  ./cmd/openkata-eval/` compiles; running with a
  skill path prints discovered scenarios; missing
  path prints usage; single scenario mode detected
- **Verify**: `go build -o bin/openkata-eval ./cmd/openkata-eval/`

### 2. Completer interface and kiro-cli backend
- **Status**: Done
- **Goal**: Define `Completer` interface and implement
  kiro-cli backend that shells out to
  `kiro-cli chat --model X --no-interactive "prompt"`,
  strips ANSI codes and header lines, parses clean
  response text, handles errors and timeouts
- **Boundary**: `cmd/openkata-eval/completer.go`,
  `cmd/openkata-eval/kiro.go`,
  `cmd/openkata-eval/kiro_test.go`
- **Key files**: None (new code)
- **Depends**: 1
- **Done when**: Kiro backend can send a prompt and
  return clean response text; tests verify ANSI
  stripping and output parsing
- **Verify**: `go test ./cmd/openkata-eval/...`

### 3. Agent prompt construction
- **Status**: Done
- **Goal**: Build agent prompt from skill context:
  SKILL.md + references/ (minus ACKNOWLEDGMENTS.md)
  + scripts/ + assets/. Append task.md + inputs/ as
  fenced code blocks. Exclude framework internals
  (CHANGELOG.md, RATIONALE.md, .tessl-plugin/, evals/)
- **Boundary**: `cmd/openkata-eval/agent.go`,
  `cmd/openkata-eval/agent_test.go`
- **Key files**: `skills/commit-conventions/SKILL.md`,
  `skills/commit-conventions/evals/`
- **Depends**: 2
- **Done when**: Given a scenario path, builds correct
  prompt with all skill context and scenario inputs;
  tests verify inclusion/exclusion rules
- **Verify**: `go test ./cmd/openkata-eval/...`

### 4. Judge call and verdict parsing
- **Status**: Done
- **Goal**: Build judge system prompt, send agent
  response (text + diff for sandboxed) + criteria,
  parse structured JSON verdict, retry once on parse
  failure. Implement sanity check: empty/trivial
  agent response + all-pass → fail as suspicious.
- **Boundary**: `cmd/openkata-eval/judge.go`,
  `cmd/openkata-eval/judge_test.go`
- **Key files**: `skills/commit-conventions/evals/scenario-0/criteria.json`
- **Depends**: 2
- **Done when**: Judge returns pass/fail per criterion
  with reasons; handles malformed JSON with retry;
  sanity check catches rubber-stamping; tests cover
  happy path, parse failure, and suspicious results
- **Verify**: `go test ./cmd/openkata-eval/...`

### 5. Scoring, threshold, and exit logic
- **Status**: Done
- **Goal**: Compute per-scenario and overall scores.
  Full skill mode: exit 1 if below threshold. Single
  scenario mode: always exit 0 (debug tool).
- **Boundary**: `cmd/openkata-eval/scoring.go`,
  `cmd/openkata-eval/scoring_test.go`
- **Key files**: None (new code)
- **Depends**: 4
- **Done when**: Scoring computes correct percentages;
  exit codes match mode; tests cover edge cases
- **Verify**: `go test ./cmd/openkata-eval/...`

### 6. Terminal output formatting
- **Status**: Done
- **Goal**: Print colored terminal report with
  per-criterion results, scenario scores, and overall
  pass/fail with threshold
- **Boundary**: `cmd/openkata-eval/output.go`
- **Key files**: None (new code)
- **Depends**: 5
- **Done when**: Running the binary prints formatted
  results matching the spec's terminal output format
- **Verify**: `go build -o bin/openkata-eval ./cmd/openkata-eval/`

### 7. JSON output file
- **Status**: Done
- **Goal**: When `--output` flag is set, write results
  to a JSON file matching the spec schema
- **Boundary**: `cmd/openkata-eval/output.go`,
  `cmd/openkata-eval/output_test.go`
- **Key files**: None (new code)
- **Depends**: 5
- **Done when**: `--output results.json` produces valid
  JSON matching spec schema; test verifies structure
- **Verify**: `go test ./cmd/openkata-eval/...`

### 8. Docker sandbox orchestration
- **Status**: Done
- **Goal**: Implement container lifecycle: start from
  eval image, mount kiro-cli auth read-only, copy
  skill + scenario files in, configure network mode
  (host/none), run kiro-cli with --trust-all-tools
  inside container, capture stdout, copy /workspace
  out, destroy container. Respect timeout. Check
  Docker availability only when sandbox scenarios
  are present.
- **Boundary**: `cmd/openkata-eval/sandbox.go`,
  `cmd/openkata-eval/sandbox_test.go`
- **Key files**: None (new code)
- **Depends**: 1
- **Done when**: Can start a container, mount auth,
  run a command, extract output/files, destroy it.
  Lazy Docker check works. Tests use lightweight
  test image.
- **Verify**: `go test ./cmd/openkata-eval/...`

### 9. Workspace diff for sandboxed runs
- **Status**: Done
- **Goal**: Snapshot workspace before agent runs,
  compare after. Produce diff summary: added,
  modified, deleted, unchanged files. Include full
  contents of added/modified files for the judge.
- **Boundary**: `cmd/openkata-eval/diff.go`,
  `cmd/openkata-eval/diff_test.go`
- **Key files**: None (new code)
- **Depends**: 8
- **Done when**: Diff correctly identifies file
  changes; outputs structured summary; tests cover
  add, modify, delete, and no-change cases
- **Verify**: `go test ./cmd/openkata-eval/...`

### 10. Parallel execution
- **Status**: Done
- **Goal**: Run scenarios in parallel using goroutines
  (text-only) or concurrent containers (sandboxed).
  Configurable concurrency (default: 2). Collect
  results and print in scenario order. Implement
  retry with exponential backoff for rate limits.
- **Boundary**: `cmd/openkata-eval/runner.go`,
  `cmd/openkata-eval/runner_test.go`
- **Key files**: None (new code)
- **Depends**: 3, 4, 8
- **Done when**: Multiple scenarios run concurrently;
  output order is stable; rate limit retries work;
  tests verify concurrency and ordering
- **Verify**: `go test -race ./cmd/openkata-eval/...`

### 11. Dockerfile for eval environment
- **Status**: Done
- **Goal**: Create `Dockerfile.eval` with kiro-cli,
  git, bash, jq, and minimal tooling. Add Makefile
  target to build the image.
- **Boundary**: `Dockerfile.eval`, `mk/dev.mk`
- **Key files**: None (new file)
- **Depends**: None
- **Done when**: `docker build -f Dockerfile.eval -t
  openkata-eval .` produces a working image with
  kiro-cli available
- **Verify**: `docker run --rm openkata-eval kiro-cli --version`

### 12. HTTP backend (Ollama/Bedrock)
- **Status**: Done
- **Goal**: Implement HTTP-based Completer that calls
  OpenAI-compatible chat completions endpoint. Handles
  base_url, api_key, model config, and timeouts.
- **Boundary**: `cmd/openkata-eval/http.go`,
  `cmd/openkata-eval/http_test.go`
- **Key files**: None (new code)
- **Depends**: 2
- **Done when**: HTTP backend can send completions and
  parse responses; tests use httptest mock server
- **Verify**: `go test ./cmd/openkata-eval/...`

### 13. Makefile target
- **Status**: Done
- **Goal**: Add `eval-local` target to Makefile
- **Boundary**: `mk/dev.mk` or `Makefile`
- **Key files**: `mk/dev.mk` (existing targets)
- **Depends**: 1
- **Done when**: `make eval-local skills/commit-conventions`
  builds and runs the eval binary
- **Verify**: `make -n eval-local`

### 14. End-to-end integration test
- **Status**: Done
- **Goal**: Test full pipeline: config → agent prompt
  → judge → scoring → output using a mock Completer
  and testdata scenario. Cover both sandbox:true
  (mocked container) and sandbox:false paths.
- **Boundary**: `cmd/openkata-eval/eval_test.go`,
  `cmd/openkata-eval/testdata/`
- **Key files**: `skills/commit-conventions/evals/scenario-0/`
  (reference format)
- **Depends**: 3, 4, 5, 6, 7, 9, 10
- **Done when**: `go test ./cmd/openkata-eval/...` passes
  with full pipeline test using mock completer
- **Verify**: `go test -race ./cmd/openkata-eval/...`

## Progress Log

- 2026-06-19: Tasks 1-14 implemented. CLI skeleton, config
  loading, completer interface (kiro + HTTP), agent prompt
  builder, judge with retry/sanity check, scoring, terminal
  + JSON output, parallel runner, Docker sandbox with
  workspace diff, Dockerfile, Makefile targets, and
  integration tests. Validation found 7 gaps (sandbox not
  wired, HTTP URL bug, missing availability checks, direct
  mode agent config, timeout not configurable, dead code).
  All fixed and verified.
