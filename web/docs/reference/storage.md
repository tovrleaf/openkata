---
title: Where artifacts are stored
description: The S3 bucket layout, how an artifact version gets published, and how versions.json is generated.
order: 2
---

# Where artifacts are stored

Skills, rules and profiles live in the git repository, and a release
is a git tag. The servers do not read git. They read a copy in an S3
bucket, and the publish workflow keeps that copy in step with the
tags. This page describes the bucket, the publish steps, and how the
index file `versions.json` is generated. The infrastructure is
defined as code in the repository's
[`infra/` folder](https://github.com/tovrleaf/openkata/tree/main/infra).

## The bucket

The bucket is `openkata-artifacts`, in the `eu-north-1` region. It is
private: all four S3 public-access blocks are on, so nothing is
served from it directly. Visitors reach artifacts through the
website or the MCP server, and both read the bucket with an IAM
role.

| Who | Access |
|-----|--------|
| MCP server (`openkata-mcp-role`) | Read |
| Website (`openkata-web-role`) | Read |
| CI (`openkata-ci`) | Read and write |

## Layout

```text
openkata-artifacts/
├── versions.json
├── skills/
│   └── <name>/
│       └── <version>/
│           ├── SKILL.md
│           └── ...
├── rules/
│   └── <name>/
│       └── <version>/
│           ├── RULE.md
│           └── ...
└── profiles/
    └── <name>/
        └── <version>/
            ├── <name>.md
            └── ...
```

- `<version>` is the semantic version without the leading `v`, for
  example `1.0.1`. Every published version keeps its own folder, so
  older versions stay available.
- A version folder holds all the files of the artifact folder as
  they were at the tag. That includes files that are not installed,
  such as `CHANGELOG.md`, `RATIONALE.md` and `evals/`. They are left
  out when an artifact is installed or downloaded, not when it is
  stored.
- `versions.json` is the index at the root. It is not part of any
  artifact.

## Publishing a version

Pushing a tag of the form `skills/<name>/v<version>`,
`rules/<name>/v<version>` or `profiles/<name>/v<version>` starts the
publish workflow. It can also be run by hand with a tag as its input.
The workflow:

1. Splits the tag into the type, the name and the version.
2. Extracts the artifact's folder as it was at the tag, with
   `git archive`.
3. Uploads it with `aws s3 sync` to
   `<type>/<name>/<version>/` in the bucket.
4. Regenerates `versions.json` and uploads it to the bucket root.

A second workflow runs when a pull request is merged into `main`. It
re-pushes the artifact tags that point at commits in that pull
request, which starts the publish workflow again.

Publishing the same tag twice uploads to the same folder, so it
replaces that version's files. Nothing in the workflow stops it.

## How `versions.json` is made

The `generate-versions` command builds the index in one of two ways.

**From the bucket**, which is what the publish workflow runs. For
each type it lists the artifact names, and for each name the version
folders. It takes the highest version, reads the description and
tags from the frontmatter of the artifact's main file in that
version, and, for skills, reads the tessl score from
`.tessl-plugin/plugin.json`. A skill counts as published unless that
file marks it private. It then uploads the result as
`s3://openkata-artifacts/versions.json`.

**From the repository**, with `make versions`. It reads the artifact
folders in the working tree. The version is the newest git tag for
the artifact, or `0.0.0` if it has never been tagged. For skills it
also reads the eval results in
`skills/<name>/evals/results/<model>/`: the newest result file of
each model, kept only when it was run against the skill's current
version. The output is written to `web/static/versions.json`.

Eval results (`models` in the file) are written only by the
repository mode, and the tessl score by both.

### Which copy is used

There are two copies, and different parts of the system read
different ones:

| Copy | Read by |
|------|---------|
| `s3://openkata-artifacts/versions.json` | MCP server (at startup), production website (each request) |
| `web/static/versions.json` | Local website, and the catalog search in the browser |

The second copy is committed to the repository and compiled into
the website, so the catalog search shows what it held at the last
deploy.

## What is not in the bucket

Download counts are kept in two DynamoDB tables, both on-demand. The
`openkata-downloads` table has one counter for each artifact, keyed
`<type>/<name>`. The `openkata-download-events` table has one record
for each download, keyed by artifact and timestamp. The MCP server
and the website have read and write access to both.
