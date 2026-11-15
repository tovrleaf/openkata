# Git hooks

Two git hooks run the check on every commit, so nobody has to
remember it. `pre-commit` checks the changed lines of the staged text
files, and `commit-msg` checks the message. Either one stops the
commit when a sentence is flagged and prints only those sentences.
When nothing is flagged, the hooks print nothing.

Set them up only when the user asks for it. A hook blocks every
commit with a flagged sentence, including one the author keeps on
purpose. To keep such a sentence, commit with `git commit --no-verify`.

## Where the hooks go

- **One clone:** save the scripts in `.git/hooks/`. Git runs them
  from there, but does not track that folder, so they stay in that
  clone.
- **Everyone on the project:** commit them to a folder, such as
  `scripts/githook/`, and have each clone turn them on once:

  ```bash
  git config core.hooksPath scripts/githook
  ```

Make both scripts executable with `chmod +x`. Set `check` in each
script to where `scripts/check.py` is, relative to the project root.

## pre-commit

```bash
#!/usr/bin/env bash
# Check the changed lines of staged text files for hard-to-read
# sentences, and stop the commit when one is flagged.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"
check=".agents/skills/clear-writing/scripts/check.py"

files=()
while IFS= read -r file; do
  case "${file}" in
    *.md | *.markdown | *.mdx | *.txt | *.templ) files+=("${file}") ;;
  esac
done < <(git diff --cached --name-only --diff-filter=ACMR)

if [[ ${#files[@]} -eq 0 ]]; then
  exit 0
fi

if ! python3 -B "${check}" --changed HEAD "${files[@]}"; then
  echo "pre-commit: rewrite the flagged sentences, or keep them with" \
    "git commit --no-verify" >&2
  exit 1
fi
```

Some files hold hard-to-read text on purpose, such as test fixtures.
To skip them, add a `case` pattern that does nothing, like
`*/fixtures/*) ;;`, before the line for text files.

## commit-msg

Git comments and the diff that `git commit -v` adds below the
scissors line are removed before the check.

```bash
#!/usr/bin/env bash
# Check the commit message for hard-to-read sentences, and stop the
# commit when one is flagged.
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
check="${root}/.agents/skills/clear-writing/scripts/check.py"

if ! sed '/^# -\{24\} >8 -\{24\}$/,$d' "${1}" \
  | git stripspace --strip-comments \
  | python3 -B "${check}" -; then
  echo "commit-msg: rewrite the flagged sentences, or keep them with" \
    "git commit --no-verify" >&2
  exit 1
fi
```
