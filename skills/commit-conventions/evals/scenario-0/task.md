# Committing a Database Module Refactor

## Background

Your team has been working on a Node.js backend service. A senior engineer recently
refactored the database connection module from a callback-based design to async/await,
collapsing three separate files into one and adding a connection-pool size cap. The
changes have been made to the working directory but nothing has been committed yet.

The files affected are in `src/db/`: the new module is `src/db/connection.js`. Two
old files, `src/db/connect.js` and `src/db/pool.js`, have been deleted and their
logic merged into `connection.js`. A new environment variable `DB_POOL_MAX` is now
read in `src/config/env.js`.

There are no lint configuration files or `CONTRIBUTING.md` in this repository.

## Your Task

Initialize a git repository in the current directory, stage the relevant files from
`src/`, and commit the refactoring work with an appropriate commit message.

Once the commit is complete, save the output of `git log --format=fuller` to a file
named `git-log.txt` in the current directory.
