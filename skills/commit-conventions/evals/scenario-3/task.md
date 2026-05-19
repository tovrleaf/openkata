# Releasing a Breaking API Change

## Background

Your team maintains a REST API used by several downstream clients. The current
authentication scheme passes API keys as a query parameter (`?api_key=...`), which
is being retired for security reasons. The new scheme requires all clients to send
the key in the `Authorization: Bearer <token>` header instead.

This is a deliberate, planned migration. The old query-parameter method has been
removed entirely from the codebase in this change — any downstream clients still
using the old scheme will receive a 401 response. Clients must update their
integration before this release is deployed.

The change touches `src/api/router.js` and `src/middleware/auth.js`. No other files
are affected.

## Your Task

Initialize a git repository in the current directory, create minimal placeholder
files for `src/api/router.js` and `src/middleware/auth.js` (a few lines each is
fine), then write and commit the change that represents this migration.

The commit must convey the nature and impact of the change accurately to any
engineer reading the history later.

After committing, save the full output of `git log --format=fuller` to `git-log.txt`.
