---
title: Downloads
description: What counts as a download, what is recorded for each one, and where the counts appear.
order: 3
---

# Downloads

Every skill, rule and profile has a download count. It appears next
to the artifact on the site and in the MCP list tools. This page
explains what increases it, what is recorded each time, and what the
number does and does not tell you.

## What counts as a download

A download is counted when the server has assembled an artifact's
files and is about to hand them over. There are two ways to get
them, and both count:

- **An archive download.** A request for the `.tar.gz` of an
  artifact, such as the link on its page, that finds files to put in
  it.
- **An MCP install.** A call to `install_skill`, `install_rule` or
  `install_profile` that finds the artifact's files.

These do not count:

- viewing a listing, an artifact page or one of its tabs
- the MCP list and version tools
- a request for an unknown artifact or version, which returns an
  error and not files
- archive downloads from a locally running website, which is a
  development setup

Every successful request counts. Nothing removes repeats from the
same person or address, and nothing filters out scripts or bots.

## What is recorded

Each counted download makes two writes in DynamoDB.

### The counter

The `openkata-downloads` table has one row for each artifact. Its key
is `<type>/<name>`, for example `skills/grill-me`, and its
`downloads` attribute goes up by one with an atomic add. All
versions of an artifact share the same counter, so the count says how
often the artifact was fetched, not which version.

### The event

The `openkata-download-events` table gets one row for each download,
keyed by artifact and timestamp. It keeps these fields:

| Field | Value |
|-------|-------|
| `artifact` | `<type>/<name>` |
| `timestamp` | The time of the download, in UTC |
| `version` | The version that was served |
| `source` | `web` for an archive download, `mcp` for an MCP install |
| `client` | The kind of client, as described below |
| `country` | The viewer's country, for archive downloads |

For an archive download, the client is guessed from the
`User-Agent` header, checking in this order: a name containing
`claude` gives `Claude-Desktop`, `cursor` gives `Cursor`, `kiro`
gives `Kiro`, `curl` gives `curl`, and `mozilla` gives `browser`.
Anything else is `unknown`. The country comes from the
`CloudFront-Viewer-Country` header that CloudFront adds.

An MCP install always records the client as `unknown` and leaves the
country empty, because the server does not look at either.

## Failures

A download does not depend on its bookkeeping. If the counter update
fails, the error is ignored and that download is not counted. If the
event write fails, the error is logged and the download still
completes. The counts are therefore best-effort: a failed write
means a missed count. The count is also taken when the files are
ready, so a download that the client then abandons is still counted.

If the counter table cannot be read, pages and the MCP list tools
show `0` for every artifact.

## Where the counts appear

- The skills, rules and profiles listings, as `42 ↓` next to each
  entry.
- The artifact's own page, next to the `.tar.gz` link.
- The `downloads` field of every entry the `list_skills`,
  `list_rules` and `list_profiles` tools return.

The counts are read from the table each time a listing, an artifact
page or an MCP list result is built, so a new download shows up on
the next load.
