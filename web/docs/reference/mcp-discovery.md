---
title: How MCP discovery works
description: The index behind the list tools, how the tag filter and version resolution behave, and when new releases show up.
order: 1
---

# How MCP discovery works

An agent finds artifacts through a single index, not by scanning
storage. The index is a file called `versions.json`. The MCP
server reads it, and the list tools answer from it. This page
explains what is in the index, what each tool does with it, and
when a new release becomes visible.

## The index

`versions.json` sits at the root of the artifact bucket,
`s3://openkata-artifacts/versions.json`. It has one map per artifact
type, keyed by artifact name:

```json
{
  "skills": {
    "grill-me": {
      "version": "1.0.1",
      "description": "Challenge a plan, spec, ADR, or design ...",
      "tags": "category:review, category:architecture"
    }
  },
  "rules": { },
  "profiles": { }
}
```

- `version` is the newest published version of the artifact.
- `description` and `tags` come from the frontmatter of the
  artifact's main file (`SKILL.md`, `RULE.md` or `<name>.md`).
  `tags` is one string of comma-separated tags, and is left out
  when the artifact has none.
- The file can carry more than the MCP server uses. The website
  also reads the tessl score of each skill from it, for example.
  The server reads only `version`, `description` and `tags`.

The publish workflow rebuilds the file each time an artifact is
published: it copies the tagged version to the bucket, then runs
`go run ./cmd/generate-versions/`, which reads the bucket and
uploads a new `versions.json`.

## The tools

The server is named `openkata` and speaks MCP over streamable HTTP.
It keeps no session state, so every call stands alone. It has nine
tools:

| Tool | Arguments | Answers from |
|------|-----------|--------------|
| `list_skills` | `tag` (optional) | The index and the download counts |
| `list_rules` | `tag` (optional) | The index and the download counts |
| `list_profiles` | `tag` (optional) | The index and the download counts |
| `install_skill` | `skill` (required), `version` (optional) | The index and the bucket |
| `install_rule` | `rule` (required), `version` (optional) | The index and the bucket |
| `install_profile` | `profile` (required), `version` (optional) | The index and the bucket |
| `skill_versions` | `skill` (required) | The bucket |
| `rule_versions` | `rule` (required) | The bucket |
| `profile_versions` | `profile` (required) | The bucket |

### Listing

A list call goes through the artifacts of one type in the index,
applies the tag filter, and returns them sorted by name. Each entry
has `name`, `version`, `description`, `tags` and `downloads`. It does
not read the bucket. The download counts come from the
`openkata-downloads` table, and if that cannot be read the counts
are zero and the list is still returned.

### The tag filter

The `tag` argument takes one tag. The server splits each
artifact's `tags` string on commas, trims the spaces, and keeps the
artifact if one of the tags is equal to the filter.

- The comparison is exact and case-sensitive. `category:review`
  matches an artifact tagged `category:review`. `review` does not,
  and neither does `Category:Review`.
- Only one tag works. A filter such as `a,b` is compared as one
  string and matches nothing, even though the tool description
  mentions comma-separated tags: that describes the tags in the
  artifact's metadata, not the filter.
- The search on the website behaves differently. It is
  case-insensitive, matches parts of names and descriptions, and
  treats a term with a colon as a tag.

### Resolving a name and a version

`install_*` first looks the name up in the index. A name that is
not there returns an error such as `skill "x" not found`, even if
the bucket holds files for it.

Without a `version`, the install uses the `version` from the
index. With one, it reads `<type>/<name>/<version>/` from the
bucket, and if that holds no files the call returns
`version <version> not found for <name>`.

`*_versions` does not use the index. It lists the version folders
under `<type>/<name>/` in the bucket, sorted from oldest to
newest, and returns an error if there are none.

## When a new release shows up

The server reads `versions.json` once, when it starts, and keeps it
in memory. If it cannot read the file at startup, it still starts,
with empty lists. On Lambda, a start happens when a function
instance is created, so an instance can keep serving the index it
loaded for as long as it lives.

This means a release becomes visible in steps:

- **Right away:** `*_versions`, and `install_*` with an explicit
  `version`, read the bucket on every call. A version just published
  to an artifact that is already in the index can be installed by
  asking for it.
- **After the server restarts:** `list_*` and the default (latest)
  version of `install_*` follow the index, so they change only when
  an instance loads the new file. A brand-new artifact cannot be
  installed until then, because its name is not in the index the
  instance holds.

The website does not have this delay. It reads `versions.json` from
the bucket on every request.
