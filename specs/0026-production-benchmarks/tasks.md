# Tasks: Benchmarks and Version List in Production

## Tasks

### 1. Pass the model to kiro-cli in the eval sandbox
- **Status**: Pending
- **Goal**: A sandboxed result is produced by the model it is filed under
- **Boundary**: `cmd/openkata-eval/sandbox.go`,
  `cmd/openkata-eval/sandbox_test.go`, `cmd/openkata-eval/runner.go`,
  `web/docs/reference/evals.md`
- **Key files**: `sandboxRun` and `NewSandboxRunner` in
  `cmd/openkata-eval/sandbox.go`, `KiroCompleter` in `kiro.go`
- **Depends**: None
- **Done when**: the `docker exec` command includes `--model <model>`
  when a model is set, and omits it when it is empty; a test asserts
  the arguments with the fake executor for both cases; the evals
  reference page no longer says the sandbox ignores `--model`
- **Verify**: `go test ./cmd/openkata-eval/...`

### 2. Share the index types
- **Status**: Pending
- **Goal**: One definition of `versions.json`
- **Boundary**: `internal/versionindex/` (new),
  `cmd/generate-versions/main.go`, `cmd/openkata-web/handlers.go`
- **Key files**: `artifactInfo`, `versionsFile`, `modelInfo` in
  `cmd/generate-versions/main.go`; the anonymous structs in
  `loadArtifactList` and `loadArtifactDetailLocal`
- **Depends**: None
- **Done when**: the types live in `internal/versionindex`; the
  generator and the web server use them; the JSON output is byte for
  byte the same as before for the same input; existing tests pass
- **Verify**: `go test ./...`, and `go run ./cmd/generate-versions/
  --local --out <tmp>` gives the same file as before

### 3. Share the eval result parsing
- **Status**: Pending
- **Goal**: One place decides which results count
- **Boundary**: `internal/evalresults/` (new),
  `cmd/generate-versions/main.go`
- **Key files**: `loadEvalResults` in `cmd/generate-versions/main.go`
- **Depends**: 2
- **Done when**: a pure function takes report files (path and contents)
  and the current version and returns effectiveness by model; it keeps
  the newest report of each model, ignores reports for another version,
  and rounds as before; the local mode uses it; table tests cover
  several models, an old version, a malformed report and no reports
- **Verify**: `go test ./internal/evalresults/ ./cmd/generate-versions/`

### 4. Read eval results from the bucket in the index generator
- **Status**: Pending
- **Goal**: The index the publish workflow uploads has `models`
- **Boundary**: `cmd/generate-versions/main.go`,
  `cmd/generate-versions/s3_test.go`
- **Key files**: `runS3` in `cmd/generate-versions/main.go`
- **Depends**: 3
- **Done when**: for each skill, the bucket mode lists
  `evals/results/` in the latest version folder, reads the files and
  gives them to the shared function; the entry gets `models` in the same
  shape as the local mode; a bucket interface lets tests use a fake;
  tests cover a skill with results, one without, and one with results
  for another version
- **Verify**: `go test ./cmd/generate-versions/`

### 5. Production artifact page: benchmarks
- **Status**: Pending
- **Goal**: The live Benchmarks tab, shown only for the indexed version
- **Boundary**: `cmd/openkata-web/handlers.go`,
  `cmd/openkata-web/main.go`, `cmd/openkata-web/handlers_test.go`
- **Key files**: `loadArtifactDetailS3`, `loadArtifactDetailLocal`
- **Depends**: 2
- **Done when**: the web server uses a bucket interface with `GetObject`
  and `ListObjectsV2`; the production loader sets the quality score, the
  published flag and the models from the index through the same function
  as the local loader; on a page for a version other than the indexed
  one the tab and the score are not set; tests use a fake bucket with a
  skill that has benchmarks, one without, and an older version page
- **Verify**: `go test ./cmd/openkata-web/...`

### 6. Production artifact page: version list
- **Status**: Pending
- **Goal**: The version control lists every version in the bucket
- **Boundary**: `cmd/openkata-web/handlers.go`,
  `cmd/openkata-web/handlers_test.go`
- **Key files**: `loadArtifactDetailS3`, `compareVersions`
- **Depends**: 5
- **Done when**: the loader lists the version folders under
  `<type>/<name>/`, sorts them newest first, and marks the newest as
  the latest; rules and profiles get the same; tests cover several
  versions, one version, and sorting of `1.10.0` after `1.9.0`
- **Verify**: `go test ./cmd/openkata-web/...`

### 7. Document the release order
- **Status**: Pending
- **Goal**: A release carries results for its own version
- **Boundary**: `web/docs/contributing/release-artifact.md`,
  `web/docs/reference/storage.md`,
  `web/docs/reference/scoring.md`,
  `.agents/skills/openkata-ryu-release/SKILL.md`, its changelog
- **Key files**: the release steps in
  `.agents/skills/openkata-ryu-release/SKILL.md`
- **Depends**: 4
- **Done when**: the release steps say to run the evals after the
  version bump and commit the report with the release commit; the pages
  say that both modes of the generator write `models`; no page says the
  bucket mode lacks them
- **Verify**: `go test ./cmd/openkata-web/...` for the docs tree

### 8. Check that tags and the bucket agree
- **Status**: Pending
- **Goal**: An unpushed tag or a missing version is noticed
- **Boundary**: `scripts/check-published.sh` (new), `Makefile`,
  `web/docs/contributing/release-artifact.md`
- **Key files**: `scripts/list-skills.sh`, `mk/dev.mk`
- **Depends**: None
- **Done when**: the script lists local tags that are not on the remote
  and pushed tags that have no files in the bucket, and exits non-zero
  if either exists; `make published` runs it; it follows the bash-style
  rule; the release page mentions it
- **Verify**: run it, and check it against the lists in the spec

## Progress Log
