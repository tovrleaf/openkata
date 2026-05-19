# Codebase Restructuring: Rename and Consolidate Utilities

## Background

A code review has flagged an inconsistency in this JavaScript project. The codebase
has utility files spread across two directories (`src/utils/` and `src/helpers/`)
with PascalCase filenames, but the project's established style guide calls for all
utility modules to live under `src/utils/` with kebab-case filenames.

The files that need to be renamed and relocated are:

- `src/utils/StringHelpers.js` → `src/utils/string-helpers.js`
- `src/utils/DateHelpers.js` → `src/utils/date-helpers.js`
- `src/helpers/APIHelpers.js` → `src/utils/api-helpers.js`

## Your Task

Initialize a git repository in the current directory with the existing source files
as the baseline (make an initial commit first), then carry out the rename operations
and commit the result.

Write a shell script named `rename.sh` that contains the git commands you would use
to perform this rename-and-commit (script should be executable but does not need to
be run — just write it). Also save the output of `git log --oneline` to `git-log.txt`
after completing your commits.
