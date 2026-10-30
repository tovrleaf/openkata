---
title: Releasing the MCP server
description: How a change to the MCP server reaches production, what its version number is, and how to check and roll back a deploy.
order: 4
---

# Releasing the MCP server

The MCP server is the Go program in `cmd/openkata-mcp`. It runs as the
`openkata-mcp` Lambda function, and agents reach it at
`https://openkata.dev/mcp`. It has no release process of its own. A
change goes out when it is merged into `main`, in the same Deploy
workflow as the website, described in
[releasing a new version of the site](release-site.md).

## Its version

The server is not published like an artifact, so it has no release tags
of its own. It has the release number of the whole project, which comes
from the `vMAJOR.MINOR.PATCH` tags described in
[releasing a new version of the site](release-site.md). Its version has
two parts, which clients see together in the `serverInfo` of the
handshake, for example `0.15.0+4f2a9c1`:

- **The base**, `0.15.0`, is the release the binary was built from, put
  in by the deploy script. A build that is not from a release, such as
  `go run`, starts with `0.0.0`, which is `baseVersion` in
  `cmd/openkata-mcp/version.go`.
- **The part after the `+`** is the git commit that the binary was
  built from, which the Go toolchain records in every build. It is added
  automatically, so the version always says which commit is deployed. A
  build from a working tree with uncommitted changes adds `.dirty`, and a
  build without version control information, such as `go run`, gets
  `+dev`.

The deploy builds the server on every merge to `main`, so the commit
changes with every merge, even when the server's own code did not.

This is separate from the versions of skills, rules and profiles, which
come from their own tags. Releasing an artifact does not deploy the
server, and deploying the server publishes no artifacts. The server reads
the index of artifacts when it starts, as described in
[how MCP discovery works](../reference/mcp-discovery.md).

## Changing the server

1. Make the change in `cmd/openkata-mcp/`.
2. Build and test it:

   ```bash
   go build -o bin/openkata-mcp ./cmd/openkata-mcp/
   go test ./...
   ```

   Only the version string has tests, so `go test` mostly protects the
   rest of the repository.
3. To try it, run the program. It loads AWS credentials and reads the
   bucket and the tables, so it needs access to them. Outside Lambda it
   listens on `:8081`, or on the address in `OPENKATA_ADDR`. The bucket
   and the table names come from `OPENKATA_BUCKET` and
   `OPENKATA_TABLE`, and default to the production ones.
4. Open a pull request, and merge it. The Build workflow checks it, and
   the merge starts the deploy.

The tools are the interface that agents use, by name and by argument. A
tool that you rename or remove stops working for every client that calls
it.

## Deploying

The `deploy-mcp` job builds the server for Linux on arm64, zips it, and
replaces the code of the `openkata-mcp` function with
`scripts/deploy-mcp.sh`. It runs on every merge to `main`, whatever
changed. To deploy by hand:

```bash
make deploy mcp
```

Set `AWS_PROFILE` to a profile that may update the function.

A deploy replaces the code and nothing else. The function's memory,
timeout and environment variables are set when it is created, by
`infra/create-mcp-stack.sh`, and so are its role and its public Function
URL. The route from CloudFront is added by
`infra/add-mcp-to-cloudfront.sh`. The infrastructure README lists the
order in which to run them.

Because the new code starts fresh function instances, they read the
current `versions.json` when they start.

## Checking a deploy

Send the server a handshake and ask it for its tools. This is what a
client does when it connects:

```bash
curl -s https://openkata.dev/mcp \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"check","version":"0"}}}'
```

The answer holds the server's name and version. The commit in the version
should be the commit you merged. Replace the `method` with
`tools/list` and the `params` with `{}` to get the list of tools, which
should be the nine described in the discovery page.

## Going back

There is no rollback step. The deploy overwrites the function code, and
no versions or aliases are kept. To go back, revert the change and merge
the revert, which deploys again, or check out an earlier commit and
deploy it by hand.
