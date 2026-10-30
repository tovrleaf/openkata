# Validation: A Version for the Site

Validated against `spec.md` and `history.md` only. Method: read the
diff and the current files, ran `go test` for `internal/siteversion`,
`cmd/site-version`, `cmd/openkata-mcp` and `cmd/openkata-web/...` (all
pass), `go vet` on the new packages (clean), `bash -n` on the five
touched scripts (clean), `go run ./cmd/site-version` (prints
`v0.14.0`), a build with `-ldflags -X ...siteversion.Release=v9.9.9`
(the string lands in the binary, so the ldflags path is correct), and
read-only `git`/`gh release list`.

Not verifiable here: anything that runs only on GitHub Actions (the
ref guard, the `v[0-9]*` trigger filter, `origin/main` being present
in the tag-triggered checkout, the AWS deploy). Also not run:
`scripts/push.sh` (it pushes) and the live site footer.

## Requirements Check

- [x] Every page shows `Release v<version>` in the footer, linked to
  the GitHub release. MET. `Footer()` renders
  `<a href=releases/tag/vX>Release vX</a>`
  (`cmd/openkata-web/templates/footer.templ:33-39`,
  `footer.go:5-14`). `@Footer()` is used in all 13 page templates,
  docs index and docs page included (`docs.templ:23,44`). Tested by
  `footer_test.go:24-35`. `file_viewer.templ` has no footer, but it is
  a fragment inside the detail pages. No test renders a whole page
  and asserts the footer is present; coverage by inspection only.
- [x] The version is the newest `v*` tag reachable from the deployed
  commit, `dev` when none. MET. `FromGit` runs
  `git tag --merged HEAD --list 'v[0-9]*'`
  (`internal/siteversion/siteversion.go:73`), compared by number
  (`:36-65`). Tested with real repos, including a tag on another
  branch (`siteversion_test.go:98-110`). See finding 3 on how this
  differs from the `git describe` wording in the spec.
- [x] A tag that is not `v<number>.<number>.<number>` is ignored. MET.
  Regex at `siteversion.go:30`. Tests cover artifact tags, `-rc1`,
  `+build`, no `v`, 1/2/4 parts, leading zeros
  (`siteversion_test.go:21-32`).
- [x] The MCP server's reported version starts with the same number.
  MET. `serverVersion()` strips the `v` from `siteversion.Current()`
  and appends `+<commit>` (`cmd/openkata-mcp/version.go:33-65`),
  passed to `NewMCPServer` (`cmd/openkata-mcp/main.go:65`).
  `scripts/deploy-mcp.sh:17` sets the ldflags. A build with no release
  reports `0.0.0+...` (`version.go:13`), so "same number" holds for
  deployed builds only, which is what the spec implies.
- [x] `release(site)` is a valid scope in the git-naming rule. MET.
  `rules/git-naming/RULE.md:46-50`, changelog entry
  `rules/git-naming/CHANGELOG.md:12`. See finding 8 (rule version).
- [x] The release page describes the tag, the GitHub release and what
  a deploy shows. MET. `web/docs/contributing/release-site.md:66-106`
  (tag and GitHub release), `:135-139` (what a deploy shows: fails
  without a tag, `Release dev` locally). Mirrored in `RELEASING.md`
  and the `openkata-ryu-release` skill.
- [~] The proposed history is accepted or changed before any tag is
  made. PARTLY MET / unverifiable. Fourteen tags `v0.1.0`..`v0.14.0`
  and 14 GitHub releases exist, locally and on the remote, on the
  commits in `history.md` (`git for-each-ref`, `gh release list`).
  Nothing in the spec or `history.md` records that the list was
  accepted. The spec is still `status: Draft`, and open question 1 is
  still open. See findings 1 and 2.

## Decisions Check

- One project version, artifact versions separate: met
  (`versioning.md:141-159`, `siteversion.go:1-4`).
- Tag `vX.Y.Z`, starts at `v0.1.0`: met (tags exist).
- A release is a tag and a GitHub release, no root `CHANGELOG.md`:
  met (no root changelog in the diff; GitHub releases exist).
- Deploy reads the newest `v*` tag: PARTLY. The mechanism is
  `git tag --merged`, not `git describe --tags --match` (finding 3).
  Full-history checkout: `deploy.yaml:20-22,61-63`. Number built in by
  ldflags: `deploy-web.sh:24-26`, `deploy-mcp.sh:16-18`.
- Pushing a release tag deploys: met in the workflow
  (`deploy.yaml:3-8`). Unverified on GitHub.
- Tag deploys only when on the head of `main`: met in the workflow
  (`deploy.yaml:26-35`, `:65-74`). Unverified on GitHub.
- Deploy fails if the commit has no release tag: met
  (`cmd/site-version/main.go:29-36`, `deploy-*.sh`). Tested
  (`main_test.go:54-62`).
- Past tags pushed before the first deploy: met, all 14 are on the
  remote (`git ls-remote`).
- Commits after the tag not shown: met, footer shows the tag only.
- Footer on every page, linked, `Release dev` unlinked without a tag:
  met (`footer.go:9-13`, `footer_test.go:61-69`).
- Who picks the number: process only. Documented in
  `release-site.md:74-86`, `RELEASING.md:92-95`, the ryu-release
  skill. No code. Cannot be verified beyond the docs.
- `make push` pushes branch then local-only tags one at a time, only
  the two tag forms, ends with `make published`: PARTLY. Branch then
  tags one by one (`push.sh:60-66`), tag forms via `TAG_RE`
  (`tag_lib.sh:7`). It ends by calling `check-published.sh` with
  `|| true` (`push.sh:68`), the same script `make published` runs, but
  not `make published` itself, and a failed check does not change the
  exit status. Also see finding 5.
- Agents do not push on their own: `push.sh` asks, default no, but has
  `--yes`/`YES=1` (`push.sh:14-16,50-52`, `mk/dev.mk:20`) that the spec
  does not mention (finding 6).
- MCP `<release>+<commit>`: met.
- Local footer `Release dev`: met.

## Out of Scope Verified

- Confirmed: no release notes page on the site and no root
  `CHANGELOG.md` in the diff.
- Confirmed: no automatic version bump from commit messages. The
  number is chosen by hand.
- Confirmed: infrastructure is not versioned separately.

## Issues Found

1. Medium. Open question 2 is not settled and the first deploy will
   show the wrong number. `v0.14.0` is the newest tag, and the docs use
   `v0.15.0` as the example (`release-site.md:10`). A merge to `main`
   of this branch deploys with `Release v0.14.0`. Tagging `v0.15.0`
   later redeploys (tag on the head of `main`). Tagging before the
   merge only works with a merge commit, not a rebase (finding 4).
   Spec open question 2 has no answer recorded.
2. Low. `history.md:9` still says "Nothing is tagged", and the spec is
   `Draft` with open question 1 unanswered, yet the tags and GitHub
   releases exist. The artifacts disagree with the spec on whether the
   history was accepted. `scripts/history-commands.sh` (not in the
   spec's requirements) is a helper for a step that is already done.
3. Medium. The spec says the deploy uses
   `git describe --tags --match 'v[0-9]*'`. The code uses
   `git tag --merged HEAD --list` plus a numeric compare. This is
   arguably better (describe orders by ancestry distance, not number),
   but it is a deviation. Consequence: a mistaken high tag that is an
   ancestor, such as `v0.150.0` for `v0.15.0`, wins over every later
   correct tag, and the deploy raises no alarm. Nothing checks that the
   new tag is the successor of the previous one, or that it is on
   `main`.
4. Medium. Rebased merge and tags on a feature branch. A tag made on a
   feature branch is not contained in `main` after a rebase merge, so
   the site keeps the old number with no error. The docs warn about it
   (`release-site.md:104-106`) but no check enforces it. A tag pushed
   from a branch (not at the head of `main`) deploys nothing and gives
   only a workflow notice (`deploy.yaml:33`).
5. Medium. Two deploy runs for one release. `make push` pushes the
   branch and then the tag, and a push to `main` plus a tag push both
   trigger Deploy, which has no `concurrency` group
   (`deploy.yaml:1-15`). Two runs update the same Lambda functions at
   the same time. The code is identical, so the end state is the same,
   but a tag-triggered run that checks out before the branch push lands
   may skip with "not head of main". Unverified on GitHub.
6. Low. `make push YES=1` skips confirmation (`push.sh:14-16`,
   `mk/dev.mk:20`). The spec says it is the user's command and agents
   do not push on their own. The flag is documented in the script
   header and the ryu-release skill, not in the spec.
7. Low. Shallow clones. `deploy-web.sh` fails on a shallow clone via
   `docs-dates -require`, but `deploy-mcp.sh` has no such guard. A
   shallow clone that has some tags would silently report an older
   release, since `git tag --merged` sees only the commits it has.
   `-require` only catches the case with no tag at all. CI uses
   `fetch-depth: 0`, so this affects manual deploys only. The comment
   in `deploy-mcp.sh:8` says "the site", which is a copy of the web
   script's comment.
8. Low. `rules/git-naming` has the change under `[Unreleased]`
   (`CHANGELOG.md:8-13`), so the published rule version does not yet
   include `release(site)`. The spec requirement is met in the repo
   but not in a released artifact. Whether to release it is the
   implementer's choice.
9. Low. Test gaps. No test for the deploy workflow guard, for
   `push.sh`/`tag_lib.sh` (shell, no tests), or for a full-page render
   containing the footer. `version_test.go` still uses `0.3.0` as the
   sample base, which no longer matches the real base.
10. Info. `check-published.sh` skips site tags (`:37-39`) and so does
    not verify that a `vX.Y.Z` tag was deployed or has a GitHub
    release. The docs for `make published` ("artifact tags") are
    consistent with this.
